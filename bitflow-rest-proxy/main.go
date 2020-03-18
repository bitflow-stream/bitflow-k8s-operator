package main

import (
	"flag"
	"os"

	"github.com/antongulenko/golib"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	os.Exit(executeMain())
}

func executeMain() int {
	endpoint := flag.String("l", ":8080", "Listen endpoint used for the API endpoint")

	var server ProxyServer
	server.registerFlags()
	golib.RegisterLogFlags()
	golib.ParseFlags()
	golib.ConfigureLogging()

	if err := server.init(); err != nil {
		log.Errorf("Failed to initialize proxy server: %v", err)
		return 1
	}

	g := golib.NewGinTask(*endpoint)
	g.Use(cors.Default())
	server.registerEndpoints(g.Engine)

	tasks := golib.TaskGroup{
		g,
		&golib.NoopTask{
			Chan:        golib.ExternalInterrupt(),
			Description: "external interrupt",
		},
	}
	log.Debugln("Press Ctrl-C to interrupt")
	return tasks.PrintWaitAndStop()
}
