package unixfs_cat

import (
	"errors"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"github.com/ipld/go-ipld-prime"
	"testing"
)

func getCid(nd ipld.Node) (cid.Cid, error) {
	switch nd := nd.(type) {
	case *merkledag.ProtoNode:
		return nd.Cid(), nil
	case *merkledag.RawNode:
		return nd.Cid(), nil
	default:
		return cid.Undef, errors.New("unknown node")
	}
}

func TestLinks(t *testing.T) {
	nodes := []ipld.Node{
		merkledag.NodeWithData(unixfs.FilePBData([]byte("hello"), 5)),
		merkledag.NodeWithData(unixfs.FilePBData([]byte("world!"), 6)),
		merkledag.NewRawNode([]byte("foo")),
		merkledag.NewRawNode([]byte("bar")),
	}

	parents, err := ConcatNodes(nodes...)
	parent := parents[0]
	if err != nil {
		t.Fatal("concat failed", err)
	}

	links := make(map[cid.Cid]bool)
	for _, l := range parent.Links() {
		links[l.Cid] = true
	}

	for _, node := range nodes {
		cid, err := getCid(node)
		if err != nil {
			t.Fatal("getting cid failed", err)
		}

		_, ok := links[cid]
		if !ok {
			t.Fatalf("link %s not found", cid)
		}
		delete(links, cid)
	}
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

	nds, err := ConcatNodes(node1, node2)
	nd := nds[0]
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
