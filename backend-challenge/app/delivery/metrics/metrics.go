package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App interface {
	Handler() http.Handler
}

type app struct {
}

func New() (App, error) {
	return &app{}, nil
}

func (m *app) Handler() http.Handler {
	return promhttp.Handler()
}
