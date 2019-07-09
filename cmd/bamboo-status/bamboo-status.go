package main

import (
	"bamboo/internal/pkg/bamboo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("stopping ...")
		done <- true
	}()

	log.Println("starting ...")

	err := bamboo.Run(done)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	log.Println("... done")
}
