package messages

import (
	"bitcoin-go/utility"
	"math"
)

type MerkleTree struct {
	Nodes        [][]*[32]byte
	total        uint32
	maxDepth     uint32
	currentDepth uint32
	currentIndex uint32
}

func NewMerkleTree(total uint32) MerkleTree {

	tree := MerkleTree{}

	// Create an empty hash tree
	tree.maxDepth = uint32(math.Ceil(math.Log2(float64(total))))
	tree.total = total
	tree.Nodes = make([][]*[32]byte, tree.maxDepth+1)
	for i := 0; i < len(tree.Nodes); i++ {

		numItemsAtLevel := int(math.Ceil(float64(total) / (math.Exp2(float64(int(tree.maxDepth) - i)))))
		tree.Nodes[i] = make([]*[32]byte, numItemsAtLevel)
	}

	return tree
}
func (m *MerkleTree) PopulateTree(bitField []bool, hashes [][32]byte) {

	// Populate the tree
	fieldIdx := 0
	hashIdx := 0

	for m.Root() == nil {
		if m.isLeaf() {
			// if we are a leaf, we know this position's hash
			// get the next bit from bitFields
			fieldIdx++

			// set the current node in the merkle tree to the next has
			m.setCurrentNode(hashes[hashIdx])
			hashIdx++

			// go up a level
			m.up()
		} else {
			// get the left hash
			leftHash := m.leftNode()
			if leftHash == nil {
				// if we don't have the left hash
				// if the next flag bit is 0, the next hash is our current node
				flag := bitField[fieldIdx]
				fieldIdx++
				if !flag {
					// set the current node to be the next hash
					m.setCurrentNode(hashes[hashIdx])
					hashIdx++

					// sub-tree doesn't need calculation, go up
					m.up()
				} else {
					// go to the left node
					m.left()
				}
			} else if m.hasRight() {
				// get the right hash
				rightHash := m.rightNode()

				// if we don't have the right hash
				if rightHash == nil {
					// go to the right node
					m.right()

				} else {
					// combine the left and right hashes
					result := m.combine(*leftHash, *rightHash)
					m.setCurrentNode(result)

					// we've completed this sub-tree, go up
					m.up()
				}
			} else {
				// combine the left hash twice
				result := m.combine(*leftHash, *leftHash)
				m.setCurrentNode(result)

				// we've completed this sub-tree, go up
				m.up()
			}
		}
	}

	if hashIdx < len(hashes) {
		panic("we didn't get through all the hashes: error in parsing")
	}

	for fieldIdx < len(bitField) {
		if bitField[fieldIdx] {
			panic("not all flag bits were consumed.")
		}
		fieldIdx++
	}
}

func (m *MerkleTree) combine(left [32]byte, right [32]byte) [32]byte {
	var combined [32]byte
	result := utility.MerkelParent(left[:], right[:])
	copy(combined[:], result)
	return combined
}

func (m *MerkleTree) up() {
	m.currentDepth--
	m.currentIndex /= 2
	//fmt.Printf("Going Up to [%v][%v]\n", m.currentDepth, m.currentIndex)
}

func (m *MerkleTree) left() {
	m.currentDepth++
	m.currentIndex *= 2
	//fmt.Printf("Going Left to [%v][%v]\n", m.currentDepth, m.currentIndex)
}

func (m *MerkleTree) right() {
	m.currentDepth++
	m.currentIndex *= 2
	m.currentIndex++
	//fmt.Printf("Going Right to [%v][%v]\n", m.currentDepth, m.currentIndex)
}

func (m *MerkleTree) Root() *[32]byte {
	return m.Nodes[0][0]
}

func (m *MerkleTree) setCurrentNode(hash [32]byte) {
	m.Nodes[m.currentDepth][m.currentIndex] = &hash
	//fmt.Printf("Setting [%v][%v] => %v\n", m.currentDepth, m.currentIndex, hash)
}

// func (m *MerkleTree) currentNode() *[32]byte {
// 	return m.Nodes[m.currentDepth][m.currentIndex]
// }

func (m *MerkleTree) leftNode() *[32]byte {
	return m.Nodes[m.currentDepth+1][m.currentIndex*2]
}

func (m *MerkleTree) rightNode() *[32]byte {
	return m.Nodes[m.currentDepth+1][(m.currentIndex*2)+1]
}

func (m *MerkleTree) isLeaf() bool {
	return m.currentDepth == m.maxDepth
}

func (m *MerkleTree) hasRight() bool {
	return len(m.Nodes[m.currentDepth+1]) > (int(m.currentIndex)*2)+1
}
