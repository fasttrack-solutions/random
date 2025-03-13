package config

import (
	"flag"
	"github.com/fasttrack-solutions/envs"
)

var (
	GRPCPort = flag.Int("grpc-port", 3401, "Port for gRPC server")
	HTTPPort = flag.Int("http-port", 3402, "Port for HTTP server")
)

func init() {
	// Parse flags if not parsed already.
	if !flag.Parsed() {
		flag.Parse()
	}

	// Determine and read environment variables.
	flagsErr := envs.GetAllFlags()
	if flagsErr != nil {
		panic(flagsErr)
	}
}
