package main

import (
	"fmt"
	"go_train/data"
	"go_train/server"
	"os"
	"time"
)

func main() {
	lenArgs := len(os.Args[1:])
	if lenArgs != 2 {
		println("Wrong number of arguments")
		println("Use \"./main.exe server_name server\" or \"./main.exe server_name db-update\" instead")
		os.Exit(0)
	}
	var ok bool
	data.G_DB, ok = server.DbConnect(os.Args[1])
	if ok == false {
		println("Please verify that your database is running")
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for t := range ticker.C {
			println("Pinging Database...")
			fmt.Println("Tick at", t)
			checkDatabase()
		}

	}()
	if os.Args[2] == "server" {
		server.Init()
	} else if os.Args[2] == "db-update" {
		i := 1
		i++
	} else {
		println("Incorrect argument")
		println("Use \"./main.exe server_name server\" or \"./main.exe server_name db-update\"")
	}
	select {}
}

func checkDatabase() {
	err := data.G_DB.Ping()
	var ok bool
	if err != nil {
		println("Connection to Database lost, trying to reconnect...")
		for ok != true {
			data.G_DB, ok = server.DbConnect(os.Args[1])
			if ok != true {
				println("Connection to Database lost, trying to reconnect...")
			}
			time.Sleep(5 * time.Second)
		}
	}
	println("Database online")
}
