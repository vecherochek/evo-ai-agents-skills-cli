package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-skills-cli/cmd"
	"github.com/cloud-ru/evo-ai-agents-skills-cli/internal/di"
)

func main() {
	ctx := context.Background()
	container := di.GetContainer()
	defer func() {
		if err := container.Close(); err != nil {
			log.Error("failed to close DI container", "error", err)
		}
	}()

	if err := cmd.RootCMD.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
