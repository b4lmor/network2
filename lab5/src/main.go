package main

import (
	"lab5/cli"
	"lab5/proxy"
	"runtime"
)

func main() {
	port := cli.Parse()
	runtime.GOMAXPROCS(1)
	proxy.Start(port)
}
