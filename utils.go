package unixfs_cat

import (
	"errors"
	"fmt"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	mh "github.com/multiformats/go-multihash"
)

func (ndwl *nodeWithLinks) constructPbNode() (pbn *merkledag.ProtoNode, err error) {
	ndb, err := ndwl.node.GetBytes()
	if err != nil {
		return
	}

	pbn = merkledag.NodeWithData(ndb)
	err = pbn.SetCidBuilder(cid.V1Builder{
		Codec:  cid.DagProtobuf,
		MhType: mh.SHA2_256,
	})
	if err != nil {
		return
	}

	for _, l := range ndwl.links {
		err = pbn.AddRawLink("", &l)
		if err != nil {
			return
		}
	}

	return
}

func (ndwl *nodeWithLinks) concatFileNode(node ipld.Node) error {

	switch node := node.(type) {

	case *merkledag.RawNode:
		s := len(node.RawData())

		ndwl.links = append(ndwl.links, ipld.Link{Cid: node.Cid()})
		ndwl.node.AddBlockSize(uint64(s))

	case *merkledag.ProtoNode:
		un, err := unixfs.ExtractFSNode(node)
		if err != nil {
			return err
		}

		switch t := un.Type(); t {
		case unixfs.TRaw, unixfs.TFile:
		default:
			return errors.New(fmt.Sprintf("can only concat raw or file types, instead found %s", t))
		}

		s := un.FileSize()

		ndwl.links = append(ndwl.links, ipld.Link{Cid: node.Cid()})
		ndwl.node.AddBlockSize(s)

	default:
		return errors.New("unknown node")
	}

	return nil
}
