package system

import (
	"SDCC-Project-WorkerNode/utility"
	"net"
	"net/rpc"
	"strconv"
)

const (
	DefaultNetwork     = "tcp"
	AmazonAWSRegion    = "us-east-1"
	AmazonS3BucketName = "graziani-filestorage"

	DefaultArbitraryFaultToleranceLevel = 3

	RPCPort = 30000
)

func StartAcceptingRPCRequest(serviceTypeRequest interface{}, nodeId int) {

	var listener net.Listener

	listener, _ = net.Listen(DefaultNetwork, "localhost"+":"+strconv.Itoa(RPCPort+int(nodeId)))
	rpc.Register(serviceTypeRequest)

	defer func() {
		utility.CheckError(listener.Close())
	}()
	for {
		rpc.Accept(listener)
	}
}
