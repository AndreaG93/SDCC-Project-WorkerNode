package Task

import (
	"SDCC-Project/MapReduce/Input"
	"SDCC-Project/MapReduce/Registry/WorkerMutex"
	"SDCC-Project/MapReduce/Registry/WorkerResultsRegister"
)

type MapReduce struct {
}

type MapReduceInput struct {
	InputData Input.MiddleInput
}

type MapReduceOutput struct {
	Digest string
}

func (x *MapReduce) Execute(input MapReduceInput, output *MapReduceOutput) error {

	digest, rawData, err := input.InputData.PerformTask()
	if err != nil {
		return err
	}

	WorkerMutex.GetInstance().Lock()
	WorkerResultsRegister.GetInstance().Set(digest, rawData)
	WorkerMutex.GetInstance().Unlock()

	output.Digest = digest

	return nil
}