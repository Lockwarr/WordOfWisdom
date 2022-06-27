package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/Lockwarr/WordOfWisdom/client"
)

// TODO: start using config file or env variables
const address = "localhost:8080"

func main() {
	mode := os.Getenv("CLIENT_MODE")
	if mode == "local" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			switch scanner.Text() {
			case "exit":
				os.Exit(0)
			case "help":
				log.Println("Available commands: request-quote, start, run, exit, help")
			case "request-quote", "start", "run":
				err := client.Run(context.Background(), address)
				if err != nil {
					panic(err)
				}
			default:
				log.Println("Unknown command. Type help for more info")
			}
		}
	} else if mode == "docker" {
		time.Sleep(time.Second)
		// requests one quote and that's it
		err := client.Run(context.Background(), "server:8080")
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Unknown client mode. Use CLIENT_MODE=local or CLIENT_MODE=docker")
	}

}
