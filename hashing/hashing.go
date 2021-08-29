package hashing

type Hashing interface {
	AddNode(nodeId string) bool
	RemoveNode(nodeId string) bool
	GetTargetNode(reqId string) string
	// LogRingStatus For debugging purpose
	LogRingStatus()
}
