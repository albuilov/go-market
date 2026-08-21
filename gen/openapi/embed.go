// Package openapi предоставляет сгенерированную спецификацию HTTP API.
package openapi

import _ "embed"

// Specification содержит сгенерированную спецификацию OpenAPI в формате JSON.
//
//go:embed spec/go-market.swagger.json
var Specification []byte
