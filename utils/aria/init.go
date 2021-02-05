package aria

import (
	"context"
	"github.com/zyxar/argo/rpc"
	"log"
	"os/exec"
	"time"
)

var client rpc.Client
var ariaContext = context.Background()

func initAriaRpc() {
	options := []string{
		"--rpc-listen-port=6942",
		"--seed-time=0",
		"--enable-rpc",
	}

	cmd := exec.Command("aria2c", options...)
	_ = cmd.Start()
}

func init() {
	initAriaRpc()
	var err error
	client, err = rpc.New(ariaContext, "http://localhost:6942/jsonrpc", "", time.Minute, nil)
	if err != nil {
		log.Fatal(err)
	}
}