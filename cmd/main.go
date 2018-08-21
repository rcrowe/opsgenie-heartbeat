package main

import (
	"context"
	"fmt"
	"os"

	heartbeat "github.com/rcrowe/opsgenie-heartbeat"
)

func main() {
	if len(os.Args) != 3 {
		exit("incorrect usage of `opsgenie-heartbeat`. Must include api key & heartbeat name as the only args.\n")
	}

	hb := heartbeat.New(os.Args[1])
	if err := hb.Ping(context.Background(), os.Args[2]); err != nil {
		exit("failed to send heartbeat: %s", err)
	}
}

func exit(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
