package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for telnet client")
	flag.Parse()

	params := os.Args[len(os.Args)-2:]

	host := params[0]
	port := params[1]

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatalf("error connecting to server: %s\n", err.Error())
	}
	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	go process(ctx, stop, client.Send, "error sending a message")

	go process(ctx, stop, client.Receive, "error receiving a message")

	for range ctx.Done() {
		if err := client.Close(); err != nil {
			fmt.Printf("error closing net: %s\n", err.Error())
			return
		}
	}
}

func process(ctx context.Context, cancelFunc func(), processFunc func() error, errorMsg string) {
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			return
		default:
			if err := processFunc(); err != nil {
				fmt.Printf("%s: %s\n", errorMsg, err.Error())
				cancelFunc()
			}
		}
	}
}
