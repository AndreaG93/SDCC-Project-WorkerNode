package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"strconv"
	"strings"
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

const (
	digestAssociationArrayLabel = "DIGEST-ASSOCIATIONS"
)

func (x *Receive) Execute(input ReceiveInput, output *ReceiveOutput) error {

	process.GetLogger().PrintInfoLevelLabeledMessage(ReceiveTaskName, fmt.Sprintf("Received Data Digest: %s Associated to Data Digest: %s", input.ReceivedDataDigest, input.AssociatedDataDigest))

	if err := process.GetDataRegistry().Set(input.ReceivedDataDigest, input.Data); err != nil {
		return err
	}

	if input.AssociatedDataDigest != "" {
		if err := SaveDigestAssociation(input.ReceivedDataDigest, input.AssociatedDataDigest); err != nil {
			return err
		}
	}

	return nil
}

func GetDigestAssociationArray(localDigest string, reduceIndex int) ([]string, error) {

	var support []string
	output := make([]string, 0)

	key := fmt.Sprintf("%s-%s", localDigest, digestAssociationArrayLabel)
	rawData := process.GetDataRegistry().Get(key)

	if err := utility.Decoding(rawData, &support); err != nil {
		return nil, err
	}

	for _, digest := range support {

		keyIndexAsString := strings.Split(digest, "-")[1]
		if keyIndex, err := strconv.Atoi(keyIndexAsString); err != nil {
			return nil, err
		} else {
			if keyIndex == reduceIndex {
				output = append(output, digest)
			}
		}
	}

	return output, nil
}

func SaveDigestAssociation(digest string, localDigest string) error {

	var digestAssociationArray []string

	key := fmt.Sprintf("%s-%s", localDigest, digestAssociationArrayLabel)
	rawData := process.GetDataRegistry().Get(key)

	if rawData == nil {

		digestAssociationArray = make([]string, 1)
		digestAssociationArray[0] = digest

	} else {

		if err := utility.Decoding(rawData, &digestAssociationArray); err != nil {
			return err
		} else {

			for _, elem := range digestAssociationArray {
				if elem == digest {
					return nil
				}
			}

			digestAssociationArray = append(digestAssociationArray, digest)
		}
	}

	if rawData, err := utility.Encoding(digestAssociationArray); err != nil {
		return err
	} else {
		return process.GetDataRegistry().Set(key, rawData)
	}
}
