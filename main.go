package main

import (
	"liuyi-test/core"
	"liuyi-test/httpserver"
	"liuyi-test/logger"
	"liuyi-test/utils"
	"os"
	"sync"
)

func main() {

	go httpserver.Server()

	logger.CurDir = getLoggerDir()

	var waitGroutp = sync.WaitGroup{}


	config,e:=core.LoadConfig("conf/liuyi-test.yaml")
	utils.PanicError(e)

	for _,v:=range config.Traces{
		p := core.NewPipeline(&v,&config)
		go func() {
			p.Start()
			waitGroutp.Done()
		}()
		waitGroutp.Add(1)
	}

	waitGroutp.Wait()

}

func getLoggerDir()string{
	wd,err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
