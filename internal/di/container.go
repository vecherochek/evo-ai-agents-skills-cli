package di

import (
	"strings"
	"sync"

	"github.com/cloud-ru/evo-ai-agents-skills-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-skills-cli/internal/auth"
	"github.com/cloud-ru/evo-ai-agents-skills-cli/internal/config"
	"github.com/samber/do/v2"
	"github.com/samber/oops"
)

type Container struct {
	injector     do.Injector
	errorHandler *ErrorHandler
}

func NewContainer() *Container {
	injector := do.New()

	do.Provide(injector, func(i do.Injector) (*config.Config, error) {
		return config.Load()
	})

	do.Provide(injector, func(i do.Injector) (*api.Client, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, oops.Errorf("failed to get config: %w", err)
		}
		authService, err := do.Invoke[auth.IAMAuthServiceInterface](i)
		if err != nil {
			return nil, oops.Errorf("failed to get auth service: %w", err)
		}
		return api.NewClient(cfg.PublicBFFURL, cfg.TimeoutSec, authService)
	})

	do.Provide(injector, func(i do.Injector) (*api.API, error) {
		client, err := do.Invoke[*api.Client](i)
		if err != nil {
			return nil, oops.Errorf("failed to get API client: %w", err)
		}
		return api.NewAPI(client), nil
	})

	do.Provide(injector, func(i do.Injector) (auth.IAMAuthServiceInterface, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, oops.Errorf("failed to get config: %w", err)
		}

		k := strings.TrimSpace(cfg.IAMKeyID)
		sec := strings.TrimSpace(cfg.IAMSecret)
		ep := strings.TrimSpace(cfg.IAMEndpoint)
		if k == "" || sec == "" {
			// allow fallback to explicit --auth-header / AUTH_HEADER
			return nil, nil
		}
		if ep == "" {
			ep = "https://iam.api.cloud.ru"
		}
		return auth.NewIAMAuthService(k, sec, ep), nil
	})

	return &Container{
		injector:     injector,
		errorHandler: NewErrorHandler(),
	}
}

func (c *Container) GetConfig() (*config.Config, error) {
	cfg, err := do.Invoke[*config.Config](c.injector)
	if err != nil {
		return nil, c.errorHandler.HandleConfigError(err)
	}
	return cfg, nil
}

func (c *Container) GetAPI() (*api.API, error) {
	apiClient, err := do.Invoke[*api.API](c.injector)
	if err != nil {
		return nil, c.errorHandler.HandleAPIError(err)
	}
	return apiClient, nil
}

func (c *Container) Close() error {
	return nil
}

var (
	containerInstance *Container
	once              sync.Once
)

func GetContainer() *Container {
	once.Do(func() {
		containerInstance = NewContainer()
	})
	return containerInstance
}
