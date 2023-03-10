package unixfs_cat

import (
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"testing"
)

func TestLinks(t *testing.T) {
	node1 := merkledag.NodeWithData(unixfs.FilePBData([]byte("hello"), 5))
	node2 := merkledag.NodeWithData(unixfs.FilePBData([]byte("world!"), 6))

	parent, err := ConcatNodes(node1, node2)
	if err != nil {
		t.Fatal("concat failed", err)
	}

	links := make(map[cid.Cid]bool)
	for _, l := range parent.Links() {
		links[l.Cid] = true
	}

	_, ok := links[node1.Cid()]
	if !ok {
		t.Fatalf("link %s not found", node1.Cid())
	}
	delete(links, node1.Cid())

	_, ok = links[node2.Cid()]
	if !ok {
		t.Fatalf("link %s not found", node2.Cid())
	}
	delete(links, node2.Cid())

	if len(links) != 0 {
		t.Fatalf("unexpected link")
	}

}

func TestSizes(t *testing.T) {
	str1 := "foo"
	str2 := "bar"

	expected := uint64(len(str1)) + uint64(len(str2))

	node1 := merkledag.NodeWithData(unixfs.FilePBData([]byte(str1), uint64(len(str1))))
	node2 := merkledag.NodeWithData(unixfs.FilePBData([]byte(str2), uint64(len(str2))))

	nd, err := ConcatNodes(node1, node2)
	if err != nil {
		t.Fatal("concat failed", err)
	}

	n, err := unixfs.ExtractFSNode(nd)
	if err != nil {
		t.Fatal("failed to extract node", err)
	}

	s := n.FileSize()
	if s != expected {
		t.Fatalf("expected size %d but found %d", expected, s)
	}
}

func TestDirectory(t *testing.T) {
	d1 := unixfs.NewFSNode(unixfs.TDirectory)
	d1b, _ := d1.GetBytes()
	p1 := merkledag.NodeWithData(d1b)

	expected := fmt.Sprintf("can only concat raw or file types, instead found %s", unixfs.TDirectory)
	_, err := ConcatNodes(p1, p1)

	if err.Error() != expected {
		t.Fatalf("expected error %s, but instead got %s", expected, err)
	}
}
