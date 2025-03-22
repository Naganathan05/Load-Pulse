package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "Load-Pulse/Tester"
    redisDB "Load-Pulse/Service"
)

func main() {
    arg := os.Args[1];
    testObj, err := Tester.New(arg);
    if err != nil {
        log.Fatal("[ERR]: Invalid File Arguement:", err);
    }

    stop := make(chan os.Signal, 1);
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM);

    redisDB.InitRedisClient();

    go func() {
        testObj.Run();
    }();

    <- stop;
    log.Println("\n[LOG]: Gracefully Shutting Down Test Server...");
}