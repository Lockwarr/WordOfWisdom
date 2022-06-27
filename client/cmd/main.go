package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/Lockwarr/WordOfWisdom/client"
)

// TODO: start using config file
const address = ":8080"

func main() {
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
}
