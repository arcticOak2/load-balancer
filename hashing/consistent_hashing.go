package hashing

import (
	"consistent-hashing/constant"
	"consistent-hashing/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang/glog"
	"math/big"
	"sync"
)

type ConsistentHashing struct {
	ring           []model.RingNode
	connectedNodes map[string]string
	mutex          sync.Mutex
}

func NewConsistentHashing() *ConsistentHashing {

	consistentHashing := ConsistentHashing{}

	consistentHashing.ring = []model.RingNode{}
	for i := 0; i < constant.DEFAULT_RING_SIZE; i++ {
		consistentHashing.ring = append(consistentHashing.ring, model.RingNode{
			NodeId:        "dummy" + "-" + fmt.Sprint(i),
			NodeType:      model.DUMMY,
			CurrentStatus: model.NOT_APPLICABLE,
		})
	}
	consistentHashing.connectedNodes = make(map[string]string)

	return &consistentHashing
}

func getHashValue(value string) int64 {

	hashMethod := md5.New()
	hashMethod.Write([]byte(value))
	hexStr := hex.EncodeToString(hashMethod.Sum(nil))
	bigInt := big.NewInt(0)
	bigInt.SetString(hexStr, 16)

	if bigInt.Int64() < 0 {
		return bigInt.Int64() * -1
	}

	return bigInt.Int64()
}

func remove(ring []model.RingNode, index int) []model.RingNode {

	copy(ring[index:], ring[index+1:])
	return ring[:len(ring)-1]
}

func isEligible(hashing *ConsistentHashing, nodeId string) bool {

	if _, ok := hashing.connectedNodes[nodeId]; ok {

		glog.Warning("Node is already added to the ring")

		return false
	}

	return true
}

func (consistentHashing *ConsistentHashing) AddNode(nodeId string) bool {

	if !isEligible(consistentHashing, nodeId) {

		return false
	}

	isActualNodePlace := false

	counter := ""

	for i := 0; i < constant.DUMMY_NODE_COUNT+1; {
		hash := getHashValue(nodeId + counter)

		pointer := hash % int64(len(consistentHashing.ring))

		if consistentHashing.ring[pointer].NodeId == nodeId {
			counter += "-"
			continue
		}

		tempNode := model.RingNode{
			NodeId:        nodeId,
			NodeType:      model.MAIN,
			CurrentStatus: model.ACTIVE,
		}

		if isActualNodePlace {
			tempNode.NodeType = model.DUMMY
		}

		consistentHashing.mutex.Lock()

		if int(pointer) == len(consistentHashing.ring)-1 {

			consistentHashing.ring = append(consistentHashing.ring, tempNode)

			consistentHashing.connectedNodes[nodeId] = nodeId
		} else if pointer == 0 {

			consistentHashing.ring = append([]model.RingNode{tempNode}, consistentHashing.ring...)

			consistentHashing.connectedNodes[nodeId] = nodeId
		} else {

			consistentHashing.ring = append(consistentHashing.ring[:pointer],
				append([]model.RingNode{tempNode}, consistentHashing.ring[pointer:]...)...)

			consistentHashing.connectedNodes[nodeId] = nodeId
		}

		consistentHashing.mutex.Unlock()

		if tempNode.NodeType == model.MAIN {
			isActualNodePlace = true
		}

		counter += "-"
		i++
	}

	return true
}

func (consistentHashing *ConsistentHashing) RemoveNode(nodeId string) bool {

	if _, ok := consistentHashing.connectedNodes[nodeId]; !ok {

		glog.Warning("No such node exist")

		return false
	}

	var indexToDelete []int

	consistentHashing.mutex.Lock()

	for i, node := range consistentHashing.ring {

		if node.NodeId == nodeId {

			indexToDelete = append(indexToDelete, i)
		}
	}

	// should be deleted in reverse order as deleting shorter index will ruin indexes of high order
	for count := len(indexToDelete) - 1; count >= 0; count-- {

		consistentHashing.ring = remove(consistentHashing.ring, indexToDelete[count])
	}

	delete(consistentHashing.connectedNodes, nodeId)

	consistentHashing.mutex.Unlock()

	return true
}

func (consistentHashing *ConsistentHashing) GetTargetNode(reqId string) string {

	hash := getHashValue(reqId)

	pointer := hash % int64(len(consistentHashing.ring))

	for {

		tempNode := consistentHashing.ring[pointer]

		if tempNode.NodeType == model.MAIN {

			return tempNode.NodeId
		}

		pointer++

		if len(consistentHashing.ring) == int(pointer) {

			pointer = 0
		}
	}

	return ""
}

func (consistentHashing *ConsistentHashing) LogRingStatus() {

	glog.Info("Logging ring status")

	for i, node := range consistentHashing.ring {

		glog.Info("index: ", i, " node: ", node)
	}
}
