package main

import (
	"fmt"
	"log"

	"github.com/gospacex/graphx/neo4jx"
)

func main() {
	db, err := neo4jx.New(nil, config())
	if err != nil {
		log.Fatalf("neo4jx.New: %v", err)
	}
	defer db.Close(nil)

	fmt.Println("connected to Neo4j, driver:", db.Driver())
}
