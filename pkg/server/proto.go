package server

import _ "embed"

// GetProtos returns the proto files.
// Key is filename, value is the content.
func GetProtos() map[string]string {
	return map[string]string{
		"server.proto": protoServer,
	}
}

//go:embed server.proto
var protoServer string
