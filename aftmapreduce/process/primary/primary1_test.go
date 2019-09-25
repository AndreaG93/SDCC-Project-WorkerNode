package primary

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"testing"
)

func Test_primary2(t *testing.T) {
	utility.CheckError(process.Initialize(1, 0, process.PrimaryProcessType, "localhost", []string{"127.0.0.1:2181"}))
	StartWork()
}
