package utils

import (
	"os"

	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/ipld/unixfs/importer/balanced"
	"github.com/ipfs/boxo/ipld/unixfs/importer/helpers"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/sync"
	chunker "github.com/ipfs/go-ipfs-chunker"
	"github.com/multiformats/go-multihash"
)

// CalculateFileCID computes the IPFS CID of a file
func CalculateFileCID(filePath string) (cid.Cid, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return cid.Cid{}, err
	}
	defer file.Close()

	// Create a thread-safe in-memory datastore
	ds := sync.MutexWrap(datastore.NewMapDatastore())

	// Create a blockstore using the datastore
	bs := blockstore.NewBlockstore(ds)

	// Create a blockservice with the blockstore
	bsvc := blockservice.New(bs, nil)

	// Create a DAG service using the blockservice
	dagService := merkledag.NewDAGService(bsvc)

	// Set up parameters for the DAG builder
	params := helpers.DagBuilderParams{
		Maxlinks:  helpers.DefaultLinksPerBlock,
		RawLeaves: true,
		CidBuilder: cid.V1Builder{
			Codec:    cid.DagProtobuf,
			MhType:   multihash.SHA2_256,
			MhLength: -1, // Default length
		},
		Dagserv: dagService,
	}

	// Create a chunker to split the file into appropriate chunks
	db, err := params.New(chunker.NewSizeSplitter(file, 1048576))
	if err != nil {
		return cid.Cid{}, err
	}

	// Create a balanced DAG layout from the file chunks
	node, err := balanced.Layout(db)
	if err != nil {
		return cid.Cid{}, err
	}

	return node.Cid(), nil
}
