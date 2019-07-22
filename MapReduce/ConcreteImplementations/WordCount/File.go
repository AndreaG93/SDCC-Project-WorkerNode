package WordCount

import (
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenHashTable"
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenList"
	"SDCC-Project/MapReduce/ConcreteImplementations/WordCount/DataStructures/WordTokenListGroupSet"
	"SDCC-Project/MapReduce/Input"
	"SDCC-Project/utility"
	"os"
	"strings"
)

type File struct {
	MapCardinality uint
	FileDigest     string
}

func (obj File) Split() ([]Input.MiddleInput, error) {

	// TODO
	// File get from AWS

	output := make([]Input.MiddleInput, obj.MapCardinality)

	splits, err := obj.splitFile()
	utility.CheckError(err)

	for index, split := range splits {

		inputForMapTask := new(MapInput)
		(*inputForMapTask).MapCardinality = obj.MapCardinality
		(*inputForMapTask).Input = split

		output[index] = *inputForMapTask
	}

	return output, nil
}

func (obj File) Shuffle(rawDataFromMapTask [][]byte) []Input.MiddleInput {

	var err error

	output := make([]Input.MiddleInput, obj.MapCardinality)
	hashTables := make([]*WordTokenHashTable.WordTokenHashTable, len(output))

	for index := 0; index < len(rawDataFromMapTask); index++ {

		currentRawData := rawDataFromMapTask[index]

		hashTables[index], err = WordTokenHashTable.Deserialize(currentRawData)
		utility.CheckError(err)
	}

	set := WordTokenListGroupSet.New(hashTables)

	for index := 0; index < len(rawDataFromMapTask); index++ {

		currentReduceInput := new(ReduceInput)

		group := set.GetGroup(uint(index))
		currentReduceInput.Data, err = group.Serialize()
		utility.CheckError(err)

		output[index] = currentReduceInput
	}

	return output
}

func (obj File) CollectResults(rawDataFromReduceTask [][]byte) string {

	finalOutput := WordTokenList.New()

	for _, rawData := range rawDataFromReduceTask {

		data, err := WordTokenList.Deserialize(rawData)
		utility.CheckError(err)

		finalOutput.Merge(data)
	}

	finalOutput.Print()

	return ""
}

func (obj File) splitFile() ([]string, error) {

	var inputFile *os.File
	var fileInfo os.FileInfo
	var err error

	output := make([]string, (obj).MapCardinality)

	if inputFile, err = os.Open((obj).FileDigest); err != nil {
		return nil, err
	}
	defer func() {
		utility.CheckError(inputFile.Close())
	}()
	if fileInfo, err = inputFile.Stat(); err != nil {
		return nil, err
	}

	averageChunkSize := fileInfo.Size() / int64((obj).MapCardinality)
	readByte := make([]byte, 1)

	for index, currentStartByte, currentEndByte := int64(0), int64(0), averageChunkSize; ; {

		if _, err = inputFile.Seek(currentEndByte, 0); err != nil {
			return nil, err
		}
		if _, err = inputFile.Read(readByte); err != nil {
			return nil, err
		}

		currentChar := string(readByte[0])

		if strings.Compare(currentChar, " ") == 0 {

			readData := make([]byte, currentEndByte-currentStartByte)

			if _, err = inputFile.Seek(currentStartByte, 0); err != nil {
				return nil, err
			}
			if _, err = inputFile.Read(readData); err != nil {
				return nil, err
			}

			output[index] = string(readData)
			index++

			currentStartByte = currentEndByte

			if (currentEndByte + averageChunkSize) < fileInfo.Size() {

				currentEndByte += averageChunkSize

			} else {

				readData := make([]byte, fileInfo.Size()-currentStartByte)

				if _, err = inputFile.Seek(currentStartByte, 0); err != nil {
					return nil, err
				}
				if _, err = inputFile.Read(readData); err != nil {
					return nil, err
				}
				output[index] = string(readData)
				break
			}

		} else {
			currentEndByte++
		}
	}

	return output, nil
}