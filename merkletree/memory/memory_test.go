package memory

import (
	"../../merkletree"
	"../../merkletree/merkletreetest"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestNewMerkleTree(t *testing.T) {
	et := merkletreetest.WrapTesting(t)
	tree := NewMerkleTree()
	et.Assert(tree.RootNode == nil, "The tree root was not nil after init")
	et.Assert(len(tree.Nodes) == 1, "The tree did not have 1 initial level")
	_, isMerkleTree := interface{}(tree).(merkletree.MerkleTree)
	et.Assert(isMerkleTree, "The tree did not implement the MerkleTree interface")
	_, isMarshalledMerkleTree := interface{}(tree).(merkletree.MarshalledMerkleTree)
	et.Assert(isMarshalledMerkleTree, "The tree did not implement the MarshalledMerkleTree interface")
}

func TestAdd(t *testing.T) {
	et := merkletreetest.WrapTesting(t)
	tree := NewMerkleTree()

	data1 := []byte("First Leaf")
	dh1 := crypto.Keccak256Hash(data1)
	expectedHash := dh1.Hex()
	i, h := tree.Add(data1)

	et.Assert(len(dh1[:]) == 32, "The hash should be 32 bytes")
	et.Assert(i == 0, "The index of first addition was not 0")
	et.Assert(h == expectedHash, "The hash of the added node was not the keccak256 hash of the data")
	et.Assert(tree.Root() == expectedHash, "The hash of the root was not equal to the only element")

	data2 := []byte("Second Leaf")
	dh2 := crypto.Keccak256Hash(data2)
	expectedHash = dh2.Hex()
	i, h = tree.Add(data2)
	expectedRoot := crypto.Keccak256Hash(dh1[:], dh2[:]).Hex()

	et.Assert(i == 1, "The index of second addition was not 1")
	et.Assert(h == expectedHash, "The hash of the added node was not the keccak256 hash of the data")

	et.Assert(tree.Root() == expectedRoot, "The hash of the root was not the hash of the two elements")
	et.Assert(len(tree.Nodes) == 2, "The tree did not grow to 2 levels")

	data3 := []byte("Third Leaf")
	dh3 := crypto.Keccak256Hash(data3)
	expectedHash = dh3.Hex()
	i, h = tree.Add(data3)
	leftBranch := crypto.Keccak256Hash(dh1[:], dh2[:])
	rightBranch := crypto.Keccak256Hash(dh3[:], dh3[:])

	expectedRoot = crypto.Keccak256Hash(leftBranch[:], rightBranch[:]).Hex()

	et.Assert(i == 2, "The index of third addition was not 2")
	et.Assert(h == expectedHash, "The hash of the added node was not the keccak256 hash of the data")

	et.Assert(tree.Root() == expectedRoot, "The hash of the root was not correctly calculated")
	et.Assert(len(tree.Nodes) == 3, "The tree did not grow to 2 levels")
}
