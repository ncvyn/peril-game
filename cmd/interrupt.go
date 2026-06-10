package cmd

import (
	"fmt"
	"os"
	"os/signal"
)

func WaitForInterrupt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println( /*^C*/ "opy that, closing the connection...")
}
