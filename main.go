package main

import (

	"fmt"
	"time"
	"os"
	structure "go/tcp/carrozzino/structure"
	server    "go/tcp/carrozzino/server"
	client    "go/tcp/carrozzino/client"

)

func main() {

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Send file")
		err := client.SendFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}()
	var server_iterface structure.ServerInterface = 
			&server.Server{}
	server_iterface.Start()
}