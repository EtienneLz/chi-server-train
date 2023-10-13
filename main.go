package main

import (
	"go_train/data"
	"go_train/server"
	"os"
)

func main() {
	lenArgs := len(os.Args[1:])
	if lenArgs != 2 {
		println("Wrong number of arguments")
		println("Use \"./main.exe server_name server\" or \"./main.exe server_name db-update\" instead")
		os.Exit(0)
	}
	data.G_DB = server.DbConnect(os.Args[1])

	if os.Args[2] == "server" {
		server.Init()
	} else if os.Args[2] == "db-update" {
		i := 1
		i++
	} else {
		println("Incorrect argument")
		println("Use \"./main.exe server_name server\" or \"./main.exe server_name db-update\"")
	}
}
