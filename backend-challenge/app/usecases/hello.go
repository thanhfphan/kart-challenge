package usecases

import (
	"context"
	"fmt"

	"github.com/thanhfphan/kart-challenge/app/repos"
	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
	"github.com/thanhfphan/kart-challenge/pkg/logging"
)

var _ Hello = (*hello)(nil)

type Hello interface {
	SayHello(ctx context.Context, name string) (string, error)
}

type hello struct {
	cfg *config.Config
	env *env.Env
}

func newHello(cfg *config.Config, env *env.Env, repos repos.Repo) (Hello, error) {
	return &hello{
		cfg: cfg,
		env: env,
	}, nil
}

func (u *hello) SayHello(ctx context.Context, name string) (string, error) {
	log := logging.FromContext(ctx)
	log.Infof("Say hello to %s", name)

	msg := fmt.Sprintf("Hello %s", name)
	return msg, nil
}
