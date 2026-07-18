// Package middleware provides HTTP middleware components for the CLI Proxy API server.
// This file re-exports model access types and functions from the shared modelaccess package.
package middleware

import (
	"context"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/interfaces"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/modelaccess"
)

// ModelAccess is re-exported from the shared modelaccess package.
type ModelAccess = modelaccess.ModelAccess

// ContextWithModelAccess is re-exported from the shared modelaccess package.
var ContextWithModelAccess = modelaccess.ContextWithModelAccess

// ModelAccessFromContext is re-exported from the shared modelaccess package.
var ModelAccessFromContext = modelaccess.ModelAccessFromContext

// ValidateModelAccess is re-exported from the shared modelaccess package.
func ValidateModelAccess(ctx context.Context, requestedModel string) *interfaces.ErrorMessage {
	return modelaccess.ValidateModelAccess(ctx, requestedModel)
}
