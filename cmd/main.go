package main

import (
	"fmt"
	"github.com/araframework/aradg/internal/network"
	"log"
	"os"
	"os/signal"
)

func main() {
	// catch Ctrl-c and kill signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill)

	fmt.Println("Starting Ara DG...")

	cluster := network.NewCluster()
	cluster.Start()

	s := <-sc
	cluster.Stop()
	log.Printf("Got signal %s, I will cleanup and exit now\n", s)
	os.Exit(0)
}
