package io

import (
	"context"

	dag "gx/ipfs/QmYxX4VfVcxmfsj8U6T5kVtFvHsSidy9tmPyPTW5fy7H3q/go-merkledag"
	ft "gx/ipfs/QmagwbbPqiN1oa3SDMZvpTFE5tNuegF1ULtuJvA9EVzsJv/go-unixfs"
	hamt "gx/ipfs/QmagwbbPqiN1oa3SDMZvpTFE5tNuegF1ULtuJvA9EVzsJv/go-unixfs/hamt"

	ipld "gx/ipfs/QmUSyMZ8Vt4vTZr5HdDEgEfpwAXfQRuDdfCFTt7XBzhxpQ/go-ipld-format"
)

// ResolveUnixfsOnce resolves a single hop of a path through a graph in a
// unixfs context. This includes handling traversing sharded directories.
func ResolveUnixfsOnce(ctx context.Context, ds ipld.NodeGetter, nd ipld.Node, names []string) (*ipld.Link, []string, error) {
	switch nd := nd.(type) {
	case *dag.ProtoNode:
		upb, err := ft.FromBytes(nd.Data())
		if err != nil {
			// Not a unixfs node, use standard object traversal code
			lnk, err := nd.GetNodeLink(names[0])
			if err != nil {
				return nil, nil, err
			}

			return lnk, names[1:], nil
		}

		switch upb.GetType() {
		case ft.THAMTShard:
			rods := dag.NewReadOnlyDagService(ds)
			s, err := hamt.NewHamtFromDag(rods, nd)
			if err != nil {
				return nil, nil, err
			}

			out, err := s.Find(ctx, names[0])
			if err != nil {
				return nil, nil, err
			}

			return out, names[1:], nil
		default:
			lnk, err := nd.GetNodeLink(names[0])
			if err != nil {
				return nil, nil, err
			}

			return lnk, names[1:], nil
		}
	default:
		lnk, rest, err := nd.ResolveLink(names)
		if err != nil {
			return nil, nil, err
		}
		return lnk, rest, nil
	}
}
