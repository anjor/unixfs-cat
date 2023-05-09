package unixfs_cat

import (
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	unixfspb "github.com/ipfs/go-unixfs/pb"
	"github.com/multiformats/go-multihash"
)

type nodeWithLinks struct {
	node  *unixfs.FSNode
	links []ipld.Link
}

type NodeWithName struct {
	node ipld.Node
	name string
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

func (pdb ParentDagBuilder) ConstructParentDirectory(nodes ...NodeWithName) (*merkledag.ProtoNode, error) {
	ndbs, err := unixfs.NewFSNode(unixfspb.Data_Directory).GetBytes()
	if err != nil {
		return nil, err
	}
	nd := merkledag.NodeWithData(ndbs)
	err = nd.SetCidBuilder(cid.V1Builder{Codec: cid.DagProtobuf, MhType: multihash.SHA2_256})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		s, _ := node.node.Size()
		err = nd.AddRawLink(node.name, &ipld.Link{Cid: node.node.Cid(), Size: s})
		if err != nil {
			return nil, err
		}
	}

	return nd, nil
}
