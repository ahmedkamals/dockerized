package main

import (
	"os"
	"syscall"
	"os/signal"
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	println("Blocking indefinitely...")
	<-sigc
	println("Bye Bye!")
}