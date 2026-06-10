package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/cmd"
)

func main() {
	_, err := cmd.Connect()
	if err != nil {
		panic(err)
	}

	cmd.WaitForInterrupt()
}
