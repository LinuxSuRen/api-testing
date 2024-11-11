package ui

import (
	_ "embed"
	"encoding/json"
)

//go:embed package.json
var packageJSON []byte

type JSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func GetPackageJSON() (data JSON) {
	data = JSON{}
	_ = json.Unmarshal(packageJSON, &data)
	return
}
