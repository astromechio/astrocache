package main

import (
	"os"

	"github.com/astromechio/astrocache/modes/master"
	"github.com/astromechio/astrocache/modes/verifier"
	"github.com/astromechio/astrocache/modes/worker"
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
