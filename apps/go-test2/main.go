package main

import (
	"flag"
	"strings"

	"github.com/golang/glog"
	"github.com/ypapax/go-test2"
)

func main() {
	var connStr, servePort, endpoints string
	flag.StringVar(&connStr, "conn", "", "mongo db connection string")
	flag.StringVar(&servePort, "port", "8181", "server binding port")
	flag.StringVar(&endpoints, "endpoints", "current_speed,temperature", "comma separated endpoints")
	flag.Parse()
	if err := go_test2.Launch(connStr, servePort, strings.Split(endpoints, ",")); err != nil {
		glog.Error(err)
	}
}
