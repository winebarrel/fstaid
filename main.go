package main

import (
	"fstaid"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func openLog(path string) (out io.Writer, err error) {
	if path != "" {
		var file *os.File
		file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)

		if err != nil {
			return
		}

		out = file
	} else {
		out = os.Stdout
	}

	return
}

func main() {
	flags := fstaid.ParseFlag()

	config, err := fstaid.LoadConfig(flags)

	if err != nil {
		log.Fatalf("Load config failed: %s", err)
	}

	gin.SetMode(config.Global.Mode)

	out, err := openLog(config.Global.Log)

	if err != nil {
		log.Fatalf("Open log failed: %s", config.Global.Log, err)
	}

	if file, ok := out.(*os.File); ok {
		defer file.Close()
	}

	log.SetOutput(out)

	cmds, err := fstaid.NewCommands(config)

	if err != nil {
		log.Fatalf("Create commands failed: %s", err)
	}

	result := cmds.Check()

	if !result.SelfCheckIsSuccess() || !result.Primary.IsSuccess() || !result.Secondary.IsSuccess() {
		log.Fatalf("Initial check failed")
	}

	handler, err := fstaid.NewCommand("handler", &config.Handler)

	if err != nil {
		log.Fatalf("Create handler failed: %s", err)
	}

	checker, err := fstaid.NewChecker(config, cmds, handler)

	if err != nil {
		log.Fatalf("Create checker failed: %s", err)
	}

	server := fstaid.NewServer(config, checker, out)

	checker.Run()
	defer checker.Stop()

	server.Run()
}
