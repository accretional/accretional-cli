module github.com/accretional/accretional-cli

go 1.24.0

require (
	github.com/accretional/collector v0.0.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/spf13/cobra v1.8.1
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

replace github.com/accretional/collector => ../collector

