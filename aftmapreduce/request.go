package aftmapreduce

import (
	"SDCC-Project/aftmapreduce/ConcreteImplementations/wordcount"
	"SDCC-Project/aftmapreduce/registries/zookeeperclient"
	"SDCC-Project/cloud/zookeeper"
	"SDCC-Project/utility"
	"fmt"
)

const (
	PendingRequestsZNodePath  = "/PendingClientRequests"
	CompleteRequestsZNodePath = "/CompleteClientRequests"
	InitialPhase              = "0"
	AfterMapPhase             = "1"
	AfterReducePhase          = "2"
)

type Request struct {
	digest                  string
	pendingRequestZNodePath string
	transientDataZNodePath  string
	requestStatusZNodePath  string
	clientDataTypeZNodePath string
	finalOutputZNodePath    string
	zookeeperClient         *zookeeper.Client
}

func NewRequest(clientData *ClientData) *Request {

	output := new(Request)

	(*output).digest = (*clientData).GetDigest()
	(*output).zookeeperClient = zookeeperclient.GetInstance()

	(*output).pendingRequestZNodePath = fmt.Sprintf("%s/%s", PendingRequestsZNodePath, (*output).digest)
	(*output).finalOutputZNodePath = fmt.Sprintf("%s/%s", CompleteRequestsZNodePath, (*output).digest)

	(*output).transientDataZNodePath = fmt.Sprintf("%s/%s", (*output).pendingRequestZNodePath, "data")
	(*output).requestStatusZNodePath = fmt.Sprintf("%s/%s", (*output).pendingRequestZNodePath, "status")

	if !(*output).zookeeperClient.CheckZNodeExistence((*output).pendingRequestZNodePath) {

		(*output).zookeeperClient.CreateZNode((*output).pendingRequestZNodePath, nil, int32(0))
		(*output).zookeeperClient.CreateZNode((*output).requestStatusZNodePath, []byte(InitialPhase), int32(0))
		(*output).zookeeperClient.CreateZNode((*output).transientDataZNodePath, nil, int32(0))
		(*output).zookeeperClient.CreateZNode((*output).finalOutputZNodePath, nil, int32(0))
	}

	(*output).zookeeperClient.SetZNodeData((*output).pendingRequestZNodePath, (*clientData).ToByte())
	(*output).zookeeperClient.SetZNodeData((*output).clientDataTypeZNodePath, []byte((*clientData).GetTypeName()))

	return output
}

func (obj *Request) Checkpoint(data []byte) {

	currentPhase, _ := (*obj).zookeeperClient.GetZNodeData((*obj).requestStatusZNodePath)

	switch string(currentPhase) {
	case InitialPhase:
		(*obj).zookeeperClient.SetZNodeData((*obj).transientDataZNodePath, data)
		(*obj).zookeeperClient.SetZNodeData((*obj).requestStatusZNodePath, []byte(AfterMapPhase))
	case AfterMapPhase:
		(*obj).zookeeperClient.SetZNodeData((*obj).transientDataZNodePath, data)
		(*obj).zookeeperClient.SetZNodeData((*obj).requestStatusZNodePath, []byte(AfterReducePhase))
	case AfterReducePhase:
		(*obj).zookeeperClient.SetZNodeData((*obj).finalOutputZNodePath, data)

		(*obj).zookeeperClient.RemoveZNode((*obj).requestStatusZNodePath)
		(*obj).zookeeperClient.RemoveZNode((*obj).transientDataZNodePath)
		(*obj).zookeeperClient.RemoveZNode((*obj).pendingRequestZNodePath)
	}
}

func (obj *Request) GetDataFromCheckpoint() []byte {

	output, _ := (*obj).zookeeperClient.GetZNodeData((*obj).transientDataZNodePath)
	return output
}

func (obj *Request) getClientData() *ClientData {
	return getClientDataFromName((*obj).zookeeperClient, (*obj).pendingRequestZNodePath, (*obj).clientDataTypeZNodePath)
}

func (obj *Request) getStatus() string {

	output, _ := (*obj).zookeeperClient.GetZNodeData((*obj).requestStatusZNodePath)
	return string(output)
}

func InitNeededZNodePathsToManageClientsRequests(zookeeperClient *zookeeper.Client) {

	if !(*zookeeperClient).CheckZNodeExistence(PendingRequestsZNodePath) {
		(*zookeeperClient).CreateZNode(PendingRequestsZNodePath, nil, 0)
	}

	if !(*zookeeperClient).CheckZNodeExistence(CompleteRequestsZNodePath) {
		(*zookeeperClient).CreateZNode(CompleteRequestsZNodePath, nil, 0)
	}
}

func GetPendingClientsRequests(zookeeperClient *zookeeper.Client) []*Request {

	pendingClientRequests := zookeeperClient.GetChildrenList(PendingRequestsZNodePath)
	output := make([]*Request, len(pendingClientRequests))

	for index, clientRequestName := range pendingClientRequests {

		zNodePathRawClientData := fmt.Sprintf("%s/%s", PendingRequestsZNodePath, clientRequestName)
		zNodePathClientDataType := fmt.Sprintf("%s/%s", zNodePathRawClientData, "ClientDataType")

		clientData := getClientDataFromName(zookeeperClient, zNodePathRawClientData, zNodePathClientDataType)

		output[index] = NewRequest(clientData)
	}

	return output

}

func getClientDataFromName(zookeeperClient *zookeeper.Client, zNodePathRawClientData string, zNodePathClientDataType string) *ClientData {

	var output ClientData

	rawClientData, _ := zookeeperClient.GetZNodeData(zNodePathRawClientData)
	clientDataType, _ := zookeeperClient.GetZNodeData(zNodePathClientDataType)

	switch string(clientDataType) {
	case AfterReducePhase:
		output = wordcount.Input{}
		utility.CheckError(utility.Decode(rawClientData, output))
	}

	return &output

	return nil
}