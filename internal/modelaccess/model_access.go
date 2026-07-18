// Package modelaccess provides model access restriction types and validation —
// checking whether a downstream API key is allowed to request a particular model.
// This package is separate from the HTTP middleware layer so that it can be
// imported by both the HTTP middleware and the auth execution layer without
// creating import cycles.
package modelaccess

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/interfaces"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/thinking"
)

// modelAccessKey is the context key for ModelAccess.
type modelAccessKey struct{}

// ModelAccess holds the per-request model restriction for a downstream API key.
// A nil or empty AllowedModels means no restriction (backward compatible).
type ModelAccess struct {
	AllowedModels []string
}

// ContextWithModelAccess returns a child context carrying ma.
func ContextWithModelAccess(ctx context.Context, ma *ModelAccess) context.Context {
	if ma == nil {
		return ctx
	}
	return context.WithValue(ctx, modelAccessKey{}, ma)
}

// ModelAccessFromContext extracts the ModelAccess stored by ContextWithModelAccess.
// Returns nil if no ModelAccess was stored.
func ModelAccessFromContext(ctx context.Context) *ModelAccess {
	if v, ok := ctx.Value(modelAccessKey{}).(*ModelAccess); ok {
		return v
	}
	return nil
}

// ValidateModelAccess checks whether the requested model is allowed for the
// downstream API key identified in ctx. Returns nil if access is granted,
// or a 403 *interfaces.ErrorMessage if denied.
//
// If ctx carries no ModelAccess, or AllowedModels is nil/empty, access is
// granted (backward compatible — no restriction configured).
func ValidateModelAccess(ctx context.Context, requestedModel string) *interfaces.ErrorMessage {
	ma := ModelAccessFromContext(ctx)
	if ma == nil || len(ma.AllowedModels) == 0 {
		return nil
	}
	normalized := strings.TrimSpace(requestedModel)
	if normalized == "" {
		return nil // can't validate empty model; let downstream handle
	}
	// Strip thinking suffix before matching.
	normalized = strings.ToLower(thinking.ParseSuffix(normalized).ModelName)
	if normalized == "" {
		return nil
	}
	for _, pattern := range ma.AllowedModels {
		if config.MatchModelPattern(pattern, normalized) {
			return nil
		}
	}
	return &interfaces.ErrorMessage{
		StatusCode: http.StatusForbidden,
		Error:      fmt.Errorf("model %q is not allowed for this API key", requestedModel),
	}
}
