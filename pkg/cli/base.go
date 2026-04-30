package cli

import (
	"context"

	"github.com/mirkobrombin/go-cli-builder/v2/pkg/log"
	"github.com/mirkobrombin/go-foundation/pkg/di"
)

// Base is a struct that can be embedded in commands to provide common functionality.
type Base struct {
	Logger    log.Logger      `internal:"ignore"`
	Ctx       context.Context `internal:"ignore"`
	Container *di.Container   `internal:"ignore"`
}