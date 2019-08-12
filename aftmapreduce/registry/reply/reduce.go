package reply

type sameDigestReduceReply struct {
	sameDigestReplyNodeIds []int
	dataSize               int
}

type ReduceReplyRegistry struct {
	registry                     map[string]*sameDigestReduceReply
	requiredNumberOfMatches      int
	mostMatchedWorkerReplyDigest string
}

func NewReduceReplyRegistry(requiredNumberOfMatches int) *ReduceReplyRegistry {

	output := new(ReduceReplyRegistry)

	(*output).registry = make(map[string]*sameDigestReduceReply)
	(*output).requiredNumberOfMatches = requiredNumberOfMatches

	return output
}

func (obj *ReduceReplyRegistry) Add(replyDigest string, replyNodeId int, dataSize int) bool {

	reply := (*obj).registry[replyDigest]

	if reply == nil {

		reply = new(sameDigestReduceReply)

		(*reply).sameDigestReplyNodeIds = make([]int, 0)
		(*reply).dataSize = dataSize

		(*obj).registry[replyDigest] = reply
	}

	(*reply).sameDigestReplyNodeIds = append((*reply).sameDigestReplyNodeIds, replyNodeId)

	if len((*reply).sameDigestReplyNodeIds) == (*obj).requiredNumberOfMatches {
		(*obj).mostMatchedWorkerReplyDigest = replyDigest
		return true
	}

	return false
}

func (obj *ReduceReplyRegistry) GetMostMatchedReply() (string, []int, int) {

	mostMatchedReply := (*obj).registry[(*obj).mostMatchedWorkerReplyDigest]

	return (*obj).mostMatchedWorkerReplyDigest, mostMatchedReply.sameDigestReplyNodeIds, mostMatchedReply.dataSize
}