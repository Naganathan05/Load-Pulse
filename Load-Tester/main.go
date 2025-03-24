package main

import (
	"os"
	"fmt"
	"log"
	"syscall"
	"os/signal"

	"Load-Pulse/Service"
	"Load-Pulse/Load-Tester/Raft"
)

func main() {
	arg := os.Args[1];
	fmt.Println(arg);
	testObj, err := Raft.NewLoadTester(arg);
	if err != nil {
		log.Fatal("[ERR]: Invalid File Arguement:", err);
	}

	stop := make(chan os.Signal, 1);
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM);

	Service.InitRedisClient();
	Service.ResetRequestCount();

	Service.ConnectRabbitMQ();
	defer Service.CloseRabbitMQ();

	go func() {
		testObj.Run();
	}();

	<- stop;
	log.Println("\n[LOG]: Gracefully Shutting Down Test Server...");
}
