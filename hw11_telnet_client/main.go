package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
	//go-telnet --timeout=5s localhost 4242
	timeout := flag.Duration("timeout", 10*time.Second, "")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("arguments are missed: hosr or port")
	}

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "sending error: %v\n", err)
		} else {
			fmt.Fprint(os.Stderr, "...EOF\n")
		}
		cancel()
	}()
	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "receiving error: %v\n", err)
		} else {
			fmt.Fprint(os.Stderr, "...Connection was closed by peer")
		}
		cancel()
	}()

	<-ctx.Done()
}
