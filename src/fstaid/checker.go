package fstaid

import (
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Checker struct {
	Config    *Config
	Commands  *Commands
	Handler   *Command
	Out       io.Writer
	Running   bool
	WaitGroup *sync.WaitGroup
}

func NewChecker(config *Config, cmds *Commands, handler *Command) (checker *Checker, err error) {
	checker = &Checker{
		Config:    config,
		Commands:  cmds,
		Handler:   handler,
		Running:   true,
		WaitGroup: &sync.WaitGroup{},
	}

	return
}

func (checker *Checker) TouchLockFile() {
	lockFile := checker.Config.Global.LockFile()
	file, err := os.OpenFile(lockFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Fatalf("Create lock file failed: %s", err)
	}

	file.Close()
}

func (checker *Checker) CheckLockFile() {
	lockFile := checker.Config.Global.LockFile()
	_, err := os.Stat(lockFile)

	if err == nil {
		log.Fatalf("fstaid is is already locked: %s", lockFile)
	}
}

func (checker *Checker) HandleFailureWithoutShutdown(result *CheckResult) {
	checker.Running = false
	log.Println("Call handler")

	checker.Handler.Run(
		strconv.Itoa(result.Primary.ExitCode),
		strconv.FormatBool(result.Primary.Timeout),
		strconv.Itoa(result.Secondary.ExitCode),
		strconv.FormatBool(result.Secondary.Timeout))

	checker.TouchLockFile()
}

func (checker *Checker) HandleFailure(result *CheckResult) {
	checker.HandleFailureWithoutShutdown(result)
	ServerShutdown()
}

func (checker *Checker) Check() {
	result := checker.Commands.Check()

	if !result.SelfCheckIsSuccess() {
		if checker.Config.Global.ContinueIfSelfCheckFailed {
			log.Println("** Self check failed ** (Health check will continue)")
			return
		} else {
			log.Fatalf("** Self check failed **")
		}
	}

	if !result.Primary.IsSuccess() {
		if !result.Secondary.IsSuccess() {
			log.Println("** Health check failed **")
			checker.HandleFailure(result)
		} else {
			log.Println("Primary check failed, but Secondary check succeeded")
		}
	}
}

func (checker *Checker) Mainloop() {
	for checker.Running {
		checker.Check()
		time.Sleep(time.Second * time.Duration(checker.Config.Global.Interval))
	}

	checker.WaitGroup.Done()
}

func (checker *Checker) Run() {
	checker.CheckLockFile()

	log.Println("Health check started")
	checker.WaitGroup.Add(1)
	go checker.Mainloop()
}

func (checker *Checker) Stop() {
	checker.Running = false
	checker.WaitGroup.Wait()
	log.Println("Health check stopped")
}
