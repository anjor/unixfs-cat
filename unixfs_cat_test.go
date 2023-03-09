package unixfs_cat

import (
	"context"
	"github.com/ipfs/go-cid"
	client "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"testing"
)

func TestConcatNodes(t *testing.T) {
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

// Needs a local ipfs daemon running
func TestConcatNodes2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	api, err := client.NewLocalApi()
	if err != nil {
		t.Fatal("couldn't connect to local daemon")
	}

	node1 := merkledag.NodeWithData(unixfs.FilePBData([]byte("hello"), 5))
	node2 := merkledag.NodeWithData(unixfs.FilePBData([]byte("world!"), 6))

	api.Dag().Add(ctx, node1)
	api.Dag().Add(ctx, node2)

	nd, err := ConcatNodes(node1, node2)

	err = api.Dag().Add(ctx, nd)
	if err != nil {
		t.Fatalf("failed to put: %v", err)
	}

	n1, _ := api.Object().Data(ctx, path.IpfsPath(node1.Cid()))
	n2, _ := api.Object().Data(ctx, path.IpfsPath(node2.Cid()))
	n, _ := api.Object().Data(ctx, path.IpfsPath(nd.Cid()))

	t.Logf("n1: %s", n1)
	t.Logf("n1 cid: %s", node1.Cid())
	t.Logf("n2: %s", n2)
	t.Logf("n2 cid: %s", node2.Cid())
	t.Logf("n: %s", n)
	t.Logf("n cid: %s", nd.Cid())

}
