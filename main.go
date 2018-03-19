package main

import (
	"fmt"
	"os"

	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/server/master"
	"github.com/astromechio/astrocache/server/verifier"
	"github.com/astromechio/astrocache/server/worker"
)

func main() {
	mode := os.Args[1]

	switch mode {
	case "master":
		master.StartMaster()
	case "verifier":
		verifier.StartVerifier()
	case "worker":
		worker.StartWorker()
	default:
		logger.LogError(fmt.Errorf("%s is not a valid astrocache node type, please use 'master', 'verifier', or 'worker'", mode))
	}
}
