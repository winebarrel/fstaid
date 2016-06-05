package fstaid

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"strconv"
)

type Server struct {
	Config  *Config
	Engine  *gin.Engine
	Router  gin.IRouter
	Checker *Checker
}

func NewServer(config *Config, checker *Checker, out io.Writer) (server *Server) {
	engine := gin.New()
	engine.Use(gin.Recovery())

	logger := gin.LoggerWithWriter(out)
	engine.Use(logger)

	server = &Server{
		Config:  config,
		Engine:  engine,
		Router:  engine,
		Checker: checker,
	}

	if len(config.User) > 0 {
		accounts := gin.Accounts{}

		for _, u := range config.User {
			accounts[u.Userid] = u.Password
		}

		server.Router = engine.Group("", gin.BasicAuth(accounts))
	}

	server.Engine.GET("/ping", server.Ping)
	server.Router.GET("/fail", server.Fail)

	return
}

func (server *Server) Run() {
	port := strconv.Itoa(server.Config.Global.Port)
	endless.ListenAndServe(":"+port, server.Engine)
}

func ServerShutdown() {
	pid := os.Getpid()
	self, err := os.FindProcess(pid)

	if err != nil {
		log.Fatalf("Server shutdown failed: %s", err)
	}

	self.Signal(os.Interrupt)
}

func (server *Server) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (server *Server) Fail(c *gin.Context) {
	exitCode := 1
	exitCodeStr := c.Query("exit")

	if exitCodeStr != "" {
		var err error
		exitCode, err = strconv.Atoi(exitCodeStr)

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}
	}

	result := &CheckResult{
		Primary:   &CommandResult{ExitCode: exitCode},
		Secondary: &CommandResult{ExitCode: exitCode},
	}

	server.Checker.HandleFailWithoutShutdown(result)

	c.JSON(200, gin.H{
		"accepted": true,
	})

	ServerShutdown()
}
