package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-skills-cli/internal/auth"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	PublicBFFURL string `env:"PUBLIC_BFF_URL" envDefault:""`
	ProjectID    string `env:"PROJECT_ID" envDefault:""`
	AuthHeader   string `env:"AUTH_HEADER" envDefault:""`
	TimeoutSec   int    `env:"HTTP_TIMEOUT_SEC" envDefault:"60"`
	IAMKeyID     string `env:"IAM_KEY_ID" envDefault:""`
	IAMSecret    string `env:"IAM_SECRET" envDefault:""`
	IAMEndpoint  string `env:"IAM_ENDPOINT" envDefault:"https://iam.api.cloud.ru"`
	CustomerID   string `env:"CUSTOMER_ID" envDefault:""`
}

func Load() (*Config, error) {
	auth.InitCredentials()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Errorf("failed to parse environment variables: %+v", err)
		return nil, err
	}
	return cfg, nil
}
