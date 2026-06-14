package main

import (
	"fmt"
	"log"

	"github.com/gospacex/graphx/memgraphx"
)

func main() {
	db, err := memgraphx.New(nil, config())
	if err != nil {
		log.Fatalf("memgraphx.New: %v", err)
	}
	defer db.Close(nil)

	fmt.Println("connected to Memgraph, driver:", db.Driver())
}
