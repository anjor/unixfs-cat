package unixfs_cat

import (
	"errors"
	"fmt"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	unixfspb "github.com/ipfs/go-unixfs/pb"
)

func ConcatNodes(nodes ...*merkledag.ProtoNode) (*merkledag.ProtoNode, error) {
	nd := unixfs.NewFSNode(unixfspb.Data_File)
	var links []format.Link

	for _, node := range nodes {
		un, err := unixfs.ExtractFSNode(node)
		if err != nil {
			return nil, err
		}

		switch t := un.Type(); t {
		case unixfs.TRaw, unixfs.TFile:
			break
		default:
			return nil, errors.New(fmt.Sprintf("can only concat raw or file types, instead found %s", t))
		}

		s := un.FileSize()
		links = append(links, format.Link{
			Name: "",
			Cid:  node.Cid(),
		})

		nd.AddBlockSize(s)
	}
	ndb, err := nd.GetBytes()
	if err != nil {
		return nil, err
	}
	pbn := merkledag.NodeWithData(ndb)

	for _, l := range links {
		err := pbn.AddRawLink("", &l)
		if err != nil {
			return nil, err
		}
	}
	return pbn, nil
}
