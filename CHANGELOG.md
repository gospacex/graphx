# Changelog

## v1.0.0 (unreleased)

### Added

- 7 backend subpackages: neo4jx, memgraphx, agex, dgraphx, janusx, nebulax, tigergraphx
- Each subpackage: New/Close/Ping + native driver accessor + Cfg() + singleton Get/Reset
- graphx.Singleton[T] - generic, double-checked locking, per-subpackage
- observability/ - OTel tracing integration with Jaeger/Kafka/Redis exporters
- quick/ - 28 cross-backend shortcut functions
- internal/mqxbinding/ - shared Kafka/Redis handle factory (mqx ecosystem)
