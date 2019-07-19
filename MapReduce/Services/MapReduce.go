package wordcount

import (
	"SDCC-Project/MapReduce/Data"
	"SDCC-Project/MapReduce/InputData"
	"SDCC-Project/MapReduce/Registry/WorkerResultsRegister"
)

type MapReduce struct {
}

type MapReduceInput struct {
	InputData Data.Split
}

type MapReduceOutput struct {
	Digest string
}

func (x *MapReduce) Execute(input MapReduceInput, output *MapReduceOutput) error {

	digest, rawData, err := input.InputData.PerformTask()
	if err != nil {
		return err
	}

	WorkerResultsRegister.GetInstance().Set(digest, rawData)

	output.Digest = digest

	return nil
}