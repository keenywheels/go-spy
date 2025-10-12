package http

import gen "github.com/keenywheels/go-spy/internal/api/v1"

var _ gen.Handler = (*Controller)(nil)

// Controller contains http handlers
type Controller struct {
}

// New creates new controller instance
func New() *Controller {
	return &Controller{}
}
