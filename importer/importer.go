// Package importer implements utilities used to create IPFS DAGs from files
// and readers.
package importer

import (
	bal "gx/ipfs/QmagwbbPqiN1oa3SDMZvpTFE5tNuegF1ULtuJvA9EVzsJv/go-unixfs/importer/balanced"
	h "gx/ipfs/QmagwbbPqiN1oa3SDMZvpTFE5tNuegF1ULtuJvA9EVzsJv/go-unixfs/importer/helpers"
	trickle "gx/ipfs/QmagwbbPqiN1oa3SDMZvpTFE5tNuegF1ULtuJvA9EVzsJv/go-unixfs/importer/trickle"

	ipld "gx/ipfs/QmUSyMZ8Vt4vTZr5HdDEgEfpwAXfQRuDdfCFTt7XBzhxpQ/go-ipld-format"
	chunker "gx/ipfs/Qme4ThG6LN6EMrMYyf2AMywAZaGbTYxQu4njfcSSkcisLi/go-ipfs-chunker"
)

// BuildDagFromReader creates a DAG given a DAGService and a Splitter
// implementation (Splitters are io.Readers), using a Balanced layout.
func BuildDagFromReader(ds ipld.DAGService, spl chunker.Splitter) (ipld.Node, error) {
	dbp := h.DagBuilderParams{
		Dagserv:  ds,
		Maxlinks: h.DefaultLinksPerBlock,
	}

	return bal.Layout(dbp.New(spl))
}

// BuildTrickleDagFromReader creates a DAG given a DAGService and a Splitter
// implementation (Splitters are io.Readers), using a Trickle Layout.
func BuildTrickleDagFromReader(ds ipld.DAGService, spl chunker.Splitter) (ipld.Node, error) {
	dbp := h.DagBuilderParams{
		Dagserv:  ds,
		Maxlinks: h.DefaultLinksPerBlock,
	}

	return trickle.Layout(dbp.New(spl))
}
