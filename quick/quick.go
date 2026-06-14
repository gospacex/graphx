// Package quick provides one-liner constructors (4 per backend: {Get,New} × {Config,Path}).
//
// Naming convention: <Letter><Mode><Source>
//
//	Letter: N4=Neo4j, M=Memgraph, A=AGE, G=Dgraph, J=Janus, V=Nebula, T=TigerGraph, U=TuGraph
//	Mode:   O=Open(Get), P=Path
//	Source: S=Singleton, C=Construct(New)
package quick

import (
	"context"

	"github.com/gospacex/graphx"
	"github.com/gospacex/graphx/agex"
	"github.com/gospacex/graphx/dgraphx"
	"github.com/gospacex/graphx/janusx"
	"github.com/gospacex/graphx/memgraphx"
	"github.com/gospacex/graphx/nebulax"
	"github.com/gospacex/graphx/neo4jx"
	"github.com/gospacex/graphx/tigergraphx"
	"github.com/gospacex/graphx/tugraphx"
)

// Neo4j shortcuts: N4{O|P}{S|C}

func N4OS(cfg graphx.Config) (*neo4jx.Neo4j, error)   { return neo4jx.Get(context.Background(), cfg) }
func N4PS(path string) (*neo4jx.Neo4j, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return neo4jx.Get(context.Background(), cfg) }
func N4OC(cfg graphx.Config) (*neo4jx.Neo4j, error)   { return neo4jx.New(context.Background(), cfg) }
func N4PC(path string) (*neo4jx.Neo4j, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return neo4jx.New(context.Background(), cfg) }

// Memgraph shortcuts: M{O|P}{S|C}

func MOS(cfg graphx.Config) (*memgraphx.Memgraph, error)   { return memgraphx.Get(context.Background(), cfg) }
func MPS(path string) (*memgraphx.Memgraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return memgraphx.Get(context.Background(), cfg) }
func MOC(cfg graphx.Config) (*memgraphx.Memgraph, error)   { return memgraphx.New(context.Background(), cfg) }
func MPC(path string) (*memgraphx.Memgraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return memgraphx.New(context.Background(), cfg) }

// AGE shortcuts: A{O|P}{S|C}

func AOS(cfg graphx.Config) (*agex.AGE, error)   { return agex.Get(context.Background(), cfg) }
func APS(path string) (*agex.AGE, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return agex.Get(context.Background(), cfg) }
func AOC(cfg graphx.Config) (*agex.AGE, error)   { return agex.New(context.Background(), cfg) }
func APC(path string) (*agex.AGE, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return agex.New(context.Background(), cfg) }

// Dgraph shortcuts: G{O|P}{S|C}

func GOS(cfg graphx.Config) (*dgraphx.Dgraph, error)   { return dgraphx.Get(context.Background(), cfg) }
func GPS(path string) (*dgraphx.Dgraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return dgraphx.Get(context.Background(), cfg) }
func GOC(cfg graphx.Config) (*dgraphx.Dgraph, error)   { return dgraphx.New(context.Background(), cfg) }
func GPC(path string) (*dgraphx.Dgraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return dgraphx.New(context.Background(), cfg) }

// JanusGraph shortcuts: J{O|P}{S|C}

func JOS(cfg graphx.Config) (*janusx.JanusGraph, error)   { return janusx.Get(context.Background(), cfg) }
func JPS(path string) (*janusx.JanusGraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return janusx.Get(context.Background(), cfg) }
func JOC(cfg graphx.Config) (*janusx.JanusGraph, error)   { return janusx.New(context.Background(), cfg) }
func JPC(path string) (*janusx.JanusGraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return janusx.New(context.Background(), cfg) }

// Nebula shortcuts: V{O|P}{S|C}

func VOS(cfg graphx.Config) (*nebulax.Nebula, error)   { return nebulax.Get(context.Background(), cfg) }
func VPS(path string) (*nebulax.Nebula, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return nebulax.Get(context.Background(), cfg) }
func VOC(cfg graphx.Config) (*nebulax.Nebula, error)   { return nebulax.New(context.Background(), cfg) }
func VPC(path string) (*nebulax.Nebula, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return nebulax.New(context.Background(), cfg) }

// TigerGraph shortcuts: T{O|P}{S|C}

func TOS(cfg graphx.Config) (*tigergraphx.TigerGraph, error)   { return tigergraphx.Get(context.Background(), cfg) }
func TPS(path string) (*tigergraphx.TigerGraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return tigergraphx.Get(context.Background(), cfg) }
func TOC(cfg graphx.Config) (*tigergraphx.TigerGraph, error)   { return tigergraphx.New(context.Background(), cfg) }
func TPC(path string) (*tigergraphx.TigerGraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return tigergraphx.New(context.Background(), cfg) }

// TuGraph shortcuts: U{O|P}{S|C}

func UOS(cfg graphx.Config) (*tugraphx.TuGraph, error)   { return tugraphx.Get(context.Background(), cfg) }
func UPS(path string) (*tugraphx.TuGraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return tugraphx.Get(context.Background(), cfg) }
func UOC(cfg graphx.Config) (*tugraphx.TuGraph, error)   { return tugraphx.New(context.Background(), cfg) }
func UPC(path string) (*tugraphx.TuGraph, error)          { cfg, err := graphx.Load(path); if err != nil { return nil, err }; return tugraphx.New(context.Background(), cfg) }
