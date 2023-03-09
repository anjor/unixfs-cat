package unixfs_cat

import (
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	unixfspb "github.com/ipfs/go-unixfs/pb"
	mh "github.com/multiformats/go-multihash"
)

func ConcatNodes(nodes ...*merkledag.ProtoNode) (*merkledag.ProtoNode, error) {
	nd := unixfs.NewFSNode(unixfspb.Data_File)
	var links []format.Link

	for _, node := range nodes {
		s := uint64(len(node.Data()))
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
	pbn.SetCidBuilder(cid.V1Builder{MhType: uint64(mh.SHA2_256)})

	for _, l := range links {
		err := pbn.AddRawLink("", &l)
		if err != nil {
			return nil, err
		}
	}
	return pbn, nil
}
