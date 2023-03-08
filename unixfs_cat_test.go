package unixfs_cat

import (
	"context"
	"github.com/ipfs/go-cid"
	client "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-libipfs/files"
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

// Needs a local ipfs daemon running
func TestConcatNodes2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	api, err := client.NewLocalApi()
	if err != nil {
		t.Fatal("couldn't connect to local daemon")
	}

	p1, err := api.Unixfs().Add(ctx, strFile("hello")())
	if err != nil {
		t.Fatal("failed to add data")
	}
	t.Logf("p1: %s", p1.String())
	n1, err := api.Object().Get(ctx, p1)
	if err != nil {
		t.Fatal("failed to get node")
	}

	p2, err := api.Unixfs().Add(ctx, strFile("world")())
	if err != nil {
		t.Fatal("failed to add data")
	}
	t.Logf("p2: %s", p2.String())

	n2, err := api.Object().Get(ctx, p2)
	if err != nil {
		t.Fatal("failed to get node")
	}

	nd1, _ := n1.(*merkledag.ProtoNode)
	nd2, _ := n2.(*merkledag.ProtoNode)

	nd, err := ConcatNodes(nd1, nd2)

	err = api.Dag().Add(ctx, nd)
	if err != nil {
		t.Fatalf("failed to put: %v", err)
	}
	t.Logf("bs: %s", nd.Cid())
}

func strFile(data string) func() files.Node {
	return func() files.Node {
		return files.NewBytesFile([]byte(data))
	}
}
