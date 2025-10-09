package security

import (
	gen "github.com/keenywheels/go-spy/internal/api/v1"
)

var _ gen.SecurityHandler = (*Controller)(nil)

// Controller contains http handlers
type Controller struct {
	clients map[string]string
}

// New creates new controller instance
func New(clients map[string]string) *Controller {
	return &Controller{
		clients: clients,
	}
}
