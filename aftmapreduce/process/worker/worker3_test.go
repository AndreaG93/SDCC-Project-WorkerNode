package worker

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"testing"
)

func Test_worker3(t *testing.T) {
	utility.CheckError(process.Initialize(3, 1, process.WorkerProcessType, "localhost", []string{"127.0.0.1:2181"}))
	StartWork()
}
