package wordcount

import (
	"SDCC-Project/aftmapreduce/node"
	"SDCC-Project/aftmapreduce/utility"
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
	utility.CheckError(node.GetDataRegistry().Set(input.ReceivedDataDigest, input.Data))

	if input.AssociatedDataDigest != "" {
		utility.CheckError(SaveGuidAssociation(input.AssociatedDataDigest, input.ReceivedDataDigest))
	}

	return nil
}

func GetGuidAssociation(guid string) []string {

	output := []string{}

	key := fmt.Sprintf("%s-GUID-ASSOCIATIONS", guid)
	rawData := node.GetDataRegistry().Get(key)

	utility.Decode(rawData, &output)

	return output
}

func SaveGuidAssociation(guid string, associatedGuid string) error {

	key := fmt.Sprintf("%s-GUID-ASSOCIATIONS", guid)
	rawData := node.GetDataRegistry().Get(key)

	if rawData == nil {

		newAssociationGuidVector := make([]string, 1)
		newAssociationGuidVector[0] = associatedGuid

		return node.GetDataRegistry().Set(key, utility.Encode(newAssociationGuidVector))
	} else {

		existentAssociatedGuidVector := []string{}
		utility.Decode(rawData, &existentAssociatedGuidVector)

		for _, elem := range existentAssociatedGuidVector {
			if elem == associatedGuid {
				return nil
			}
		}

		existentAssociatedGuidVector = append(existentAssociatedGuidVector, associatedGuid)
		return node.GetDataRegistry().Set(key, utility.Encode(existentAssociatedGuidVector))
	}
}
