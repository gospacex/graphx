package main

import "github.com/gospacex/graphx"

func config() graphx.Config {
	return graphx.Config{
		Address:  "localhost:7687",
		Username: "neo4j",
		Password: "password",
		PoolSize: 5,
	}
}
