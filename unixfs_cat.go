package unixfs_cat

import (
	"errors"
	"fmt"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"github.com/ipfs/go-unixfs/importer/helpers"
	unixfspb "github.com/ipfs/go-unixfs/pb"
)

type nodeWithLinks struct {
	node  *unixfs.FSNode
	links []ipld.Link
}

func ConcatNodes(nodes ...ipld.Node) ([]*merkledag.ProtoNode, error) {
	var pbns []*merkledag.ProtoNode
	ndwl := nodeWithLinks{node: unixfs.NewFSNode(unixfspb.Data_File)}
	for _, node := range nodes {
		if len(ndwl.node.BlockSizes()) < helpers.DefaultLinksPerBlock {
			if err := ndwl.addLink(node); err != nil {
				return nil, err
			}
		} else {
			pbn, err := ndwl.constructPbNode()
			if err != nil {
				return nil, err
			}

			pbns = append(pbns, pbn)

			ndwl = nodeWithLinks{node: unixfs.NewFSNode(unixfspb.Data_File)}
			if err := ndwl.addLink(node); err != nil {
				return nil, err
			}

		}
	}

	pbn, err := ndwl.constructPbNode()
	if err != nil {
		return nil, err
	}

	pbns = append(pbns, pbn)

	return pbns, nil
}

func (ndwl *nodeWithLinks) constructPbNode() (pbn *merkledag.ProtoNode, err error) {
	ndb, err := ndwl.node.GetBytes()
	if err != nil {
		return
	}

	pbn = merkledag.NodeWithData(ndb)

	for _, l := range ndwl.links {
		err = pbn.AddRawLink("", &l)
		if err != nil {
			return
		}
	}

	return
}

func (ndwl *nodeWithLinks) addLink(node ipld.Node) error {

	switch node := node.(type) {

	case *merkledag.RawNode:
		s := len(node.RawData())

		ndwl.links = append(ndwl.links, ipld.Link{Name: "", Cid: node.Cid()})
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

		ndwl.links = append(ndwl.links, ipld.Link{Name: "", Cid: node.Cid()})
		ndwl.node.AddBlockSize(s)

	default:
		return errors.New("unknown node")
	}

	return nil
}
