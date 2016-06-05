package fstaid

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-shellwords"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Command struct {
	Name    string
	CmdArgs []string
	Timeout time.Duration
}

func NewCommand(name string, config *CommandConfig) (cmd *Command, err error) {
	cmdArgs, err := shellwords.Parse(config.Command)

	if err != nil {
		return
	}

	cmd = &Command{
		Name:    name,
		CmdArgs: cmdArgs,
		Timeout: time.Second * time.Duration(config.Timeout),
	}

	return
}

func makeCmd(cmdArgs []string) (cmd *exec.Cmd, outReader io.ReadCloser, errReader io.ReadCloser, err error) {
	if len(cmdArgs) > 1 {
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	} else {
		cmd = exec.Command(cmdArgs[0])
	}

	outReader, err = cmd.StdoutPipe()

	if err != nil {
		return
	}

	errReader, err = cmd.StderrPipe()

	if err != nil {
		return
	}

	return
}

func getExitCode(cmdErr error) (exitCode int, err error) {
	exitErr, ok := cmdErr.(*exec.ExitError)

	if !ok {
		err = cmdErr
		return
	}

	status, ok := exitErr.Sys().(syscall.WaitStatus)

	if !ok {
		err = cmdErr
		return
	}

	exitCode = status.ExitStatus()
	return

}

func tailf(prefix string, reader io.Reader, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		log.Printf("%s: %s\n", prefix, scanner.Text())
	}

	wg.Done()
}

func (command *Command) Run(args ...string) (exitCode int, timeout bool) {
	cmdArgs := append(command.CmdArgs, args...)
	cmd, outReader, errReader, err := makeCmd(cmdArgs)

	if err != nil {
		log.Fatalf("%s: %s", command.Name, err)
	}

	wg := &sync.WaitGroup{}

	if gin.Mode() == "debug" {
		wg.Add(1)
		go tailf(command.Name+": stdout", outReader, wg)
	}

	wg.Add(1)
	go tailf(command.Name+": stderr", errReader, wg)

	err = cmd.Start()

	if err != nil {
		log.Fatalf("%s: %s", command.Name, err)
	}

	var timer *time.Timer

	timer = time.AfterFunc(command.Timeout, func() {
		timer.Stop()
		cmd.Process.Kill()
		timeout = true
	})

	err = cmd.Wait()
	timer.Stop()
	wg.Wait()

	if err != nil {
		exitCode, err = getExitCode(err)

		if err != nil {
			log.Fatalf("%s: %s", command.Name, err)
		}
	}

	if timeout {
		log.Printf("%s: failed: timed out\n", command.Name)
	} else if exitCode != 0 {
		log.Printf("%s: failed: exitCode=%d\n", command.Name, exitCode)
	}

	return
}
