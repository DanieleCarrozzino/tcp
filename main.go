package main

import (

	"fmt"
	"time"
	"os"

	//server client communication
	structure 	"go/tcp/carrozzino/structure"
	server    	"go/tcp/carrozzino/server"
	client    	"go/tcp/carrozzino/client"

	//cert and key creation
	creation 	"go/tcp/carrozzino/creation"
)

func main() {

	if true {
		var client_generator structure.ClientGenerator =
			&creation.Gen{} 
		client_generator.CreateKey();
		client_generator.CreateCSR();
		client_generator.CreateCRT();
	}

	go func() {
		time.Sleep(2 * time.Second)
		err := client.SendFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		time.Sleep(2 * time.Second)
		err = client.SendFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}()
	var server_iterface structure.ServerInterface = 
			&server.Server{}
	server_iterface.Start()
}