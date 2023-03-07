package unixfs_cat

import (
	"bytes"
	"github.com/ipfs/go-merkledag"
	"testing"
)

func TestConcatNodes(t *testing.T) {
	node1 := merkledag.NodeWithData([]byte("hello"))
	node2 := merkledag.NodeWithData([]byte("world!"))

	parent, err := ConcatNodes(node1, node2)
	if err != nil {
		t.Fatal("concat failed", err)
	}
	m, err := parent.MarshalJSON()
	if err != nil {
		t.Fatal("failed to marshal", err)
	}

	n1, err := node1.MarshalJSON()
	if err != nil {
		t.Fatal("failed to marshal", err)
	}

	n2, err := node2.MarshalJSON()
	if err != nil {
		t.Fatal("failed to marshal", err)
	}

	t.Logf("node 1 := %v", bytes.NewBuffer(n1).String())
	t.Logf("node 2 := %v", bytes.NewBuffer(n2).String())
	t.Logf("parent := %v", bytes.NewBuffer(m).String())

	t.Logf("node 1 := %v", node1.String())
	t.Logf("node 2 := %v", node2.String())
	t.Logf("parent := %v", parent.String())
}
