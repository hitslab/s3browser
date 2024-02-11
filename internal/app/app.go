package app

import (
	"context"

	"github.com/hitslab/s3browser/internal/config"
	"github.com/hitslab/s3browser/internal/s3"
)

type S3Settings struct {
	BaseUrl    string `json:"base_url"`
	Region     string `json:"region"`
	Bucket     string `json:"bucket"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	PathStyle  bool   `json:"path_style"`
	DisableSSL bool   `json:"disable_ssl"`
}

// App struct
type App struct {
	ctx context.Context
	s3  s3.S3
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.s3 = s3.S3{}
}

func (a *App) GetConfig() (config.Config, error) {
	return config.GetConfig()
}

func (a *App) Connect(s config.S3Settings) error {
	cfg := config.Config{}
	cfg.S3 = s

	err := config.SetConfig(a.ctx, cfg)
	if err != nil {
		return err
	}

	err = a.s3.Connect(s)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) List(prefix string) (s3.List, error) {
	return a.s3.List(prefix)
}
