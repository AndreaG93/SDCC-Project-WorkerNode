package main

import (
	"SDCC-Project-WorkerNode/mapreduce/wordcount/system/workernode"
	"SDCC-Project-WorkerNode/utility"
	"os"
	"strconv"
)

func main() {

	myId, err := strconv.Atoi(os.Getenv("NODE_ID"))
	utility.CheckError(err)

	node := workernode.New(uint(myId), []string{"localhost"})
	node.StartWork()
}