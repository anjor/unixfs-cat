package unixfs_cat

import (
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	unixfspb "github.com/ipfs/go-unixfs/pb"
)

type nodeWithLinks struct {
	node  *unixfs.FSNode
	links []ipld.Link
}

type ParentDagBuilder struct {
	maxLinks int
}

func (pdb ParentDagBuilder) ConcatFileNodes(nodes ...ipld.Node) ([]*merkledag.ProtoNode, error) {
	var pbns []*merkledag.ProtoNode
	ndwl := nodeWithLinks{node: unixfs.NewFSNode(unixfspb.Data_File)}
	for _, node := range nodes {
		if len(ndwl.node.BlockSizes()) < pdb.maxLinks {
			if err := ndwl.concatFileNode(node); err != nil {
				return nil, err
			}
		} else {
			pbn, err := ndwl.constructPbNode()
			if err != nil {
				return nil, err
			}

			pbns = append(pbns, pbn)

			ndwl = nodeWithLinks{node: unixfs.NewFSNode(unixfspb.Data_File)}
			if err := ndwl.concatFileNode(node); err != nil {
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
