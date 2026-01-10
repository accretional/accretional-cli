module github.com/accretional/accretional-cli

go 1.24.0

require (
	github.com/accretional/collector v0.0.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/spf13/cobra v1.8.1
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
)

replace github.com/accretional/collector => ../collector
