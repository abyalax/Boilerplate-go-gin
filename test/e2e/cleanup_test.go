package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/bootstrap"
	"github.com/abyalax/Boilerplate-go-gin/src/config/env"
)

func TestMain(m *testing.M) {
	code := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, _ := env.Load()
	if cfg != nil {
		app, _ := bootstrap.NewApp(cfg)
		if app != nil {
			app.Stop(ctx)
		}
	}

	os.Exit(code)
}
