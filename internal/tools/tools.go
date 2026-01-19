//go:build tools
// +build tools

package tools

import (
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/swaggo/swag/cmd/swag"
	_ "sigs.k8s.io/yaml"
)
