package unixfs_cat

import (
	"errors"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	unixfspb "github.com/ipfs/go-unixfs/pb"
	"github.com/ipld/go-ipld-prime"
)

func ConcatNodes(nodes ...ipld.Node) (*merkledag.ProtoNode, error) {
	nd := unixfs.NewFSNode(unixfspb.Data_File)
	var links []format.Link

	for _, node := range nodes {
		switch node := node.(type) {

		case *merkledag.RawNode:
			s := len(node.RawData())

			links = addLink(links, node.Cid())
			nd.AddBlockSize(uint64(s))

		case *merkledag.ProtoNode:
			un, err := unixfs.ExtractFSNode(node)
			if err != nil {
				return nil, err
			}

			switch t := un.Type(); t {
			case unixfs.TRaw, unixfs.TFile:
			default:
				return nil, errors.New(fmt.Sprintf("can only concat raw or file types, instead found %s", t))
			}

			s := un.FileSize()

			links = addLink(links, node.Cid())
			nd.AddBlockSize(s)

		default:
			return nil, errors.New("unknown node")
		}
	}

	return constructPbNode(nd, links)
}

func constructPbNode(nd *unixfs.FSNode, links []format.Link) (pbn *merkledag.ProtoNode, err error) {
	ndb, err := nd.GetBytes()
	if err != nil {
		return
	}

	pbn = merkledag.NodeWithData(ndb)

	for _, l := range links {
		err = pbn.AddRawLink("", &l)
		if err != nil {
			return
		}
	}

	return
}

func addLink(links []format.Link, cid cid.Cid) []format.Link {
	return append(links, format.Link{
		Name: "",
		Cid:  cid,
	})
}
