package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/tugraphx"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := tugraphx.New(ctx, graphx.Config{
		Address:  "localhost:7688",
		Username: "admin",
		Password: "73@TuGraph",
	})
	if err != nil {
		log.Fatalf("tugraphx.New: %v", err)
	}
	defer db.Close(ctx)

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("tugraphx.Ping: %v", err)
	}
	fmt.Println("TuGraph ping OK")

	driver := db.Driver()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.Run(ctx, "RETURN 1 AS num", nil)
	if err != nil {
		log.Fatalf("session.Run: %v", err)
	}
	rec, err := result.Single(ctx)
	if err != nil {
		log.Fatalf("result.Single: %v", err)
	}
	fmt.Printf("Cypher result: %+v\n", rec.AsMap())
	os.Exit(0)
}
