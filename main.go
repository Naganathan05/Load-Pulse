package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/andresoro/bench/bench"
)

func main() {
    arg := os.Args[1]
    b, err := bench.New(arg)
    if err != nil {
        log.Fatal("Could not open file:", err)
    }

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    go func() {
        b.Run()
    }();

    <-stop
    log.Println("\nGracefully shutting down...")
}