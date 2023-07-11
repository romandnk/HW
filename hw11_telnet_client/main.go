package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var (
		timeout time.Duration
		address strings.Builder
	)

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for telnet client")
	flag.Parse()

	params := os.Args[len(os.Args)-2:]

	host := params[0]
	port := params[1]

	address.WriteString(host + ":" + port)

	client := NewTelnetClient(address.String(), timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatalf("error connecting to server: %s\n", err.Error())
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Send(); err != nil {
					fmt.Printf("error sendind a message: %s\n", err.Error())
					cancel()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Receive(); err != nil {
					fmt.Printf("error receiving a message: %s\n", err.Error())
					cancel()
				}
			}
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-done:
			cancel()
			return
		case <-ctx.Done():
			if err := client.Close(); err != nil {
				fmt.Printf("error closing net: %s\n", err.Error())
				return
			}
		}
	}
}
