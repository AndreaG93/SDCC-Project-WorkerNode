package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/wordcount/DataStructures/WordTokenList"
	"fmt"
)

type Receive struct {
}

type ReceiveInput struct {
	Data                 []byte
	ReceivedDataDigest   string
	AssociatedDataDigest string
}

type ReceiveOutput struct {
}

func (x *Receive) Execute(input ReceiveInput, output *ReceiveOutput) error {

	node.GetLogger().PrintInfoTaskMessage(ReceiveTaskName, fmt.Sprintf("Received Data Digest: %s Associated to Data Digest: %s", input.ReceivedDataDigest, input.AssociatedDataDigest))

	receivedWordTokenList := WordTokenList.Deserialize(input.Data)

	node.GetDataRegistry().Set(input.ReceivedDataDigest, receivedWordTokenList.Serialize())
	node.GetDigestRegistry().Add(input.AssociatedDataDigest, input.ReceivedDataDigest)

	return nil
}
