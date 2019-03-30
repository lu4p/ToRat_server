package main

import "github.com/lu4p/ToRat_server/server"
import _ "github.com/dimiro1/banner/autoload"

func main() {
	go server.Start()
	server.Shell()
}
