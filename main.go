package main

import (
	"os"

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
	}
}
