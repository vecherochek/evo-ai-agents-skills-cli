package skill_marketplace

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vecherochek/evo-ai-agents-skills-cli/internal/di"
)

var addOptions struct {
	skillID    string
	outputDir  string
	authHeader string
}

var AddCMD = &cobra.Command{
	Use:   "add",
	Short: "Download and extract a marketplace skill archive from BFF",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := di.GetContainer()
		config, err := container.GetConfig()
		if err != nil {
			return err
		}
		apiClient, err := container.GetAPI()
		if err != nil {
			return err
		}

		if strings.TrimSpace(addOptions.skillID) == "" {
			return fmt.Errorf("flag --skill-id is required")
		}
		authHeader := strings.TrimSpace(addOptions.authHeader)
		if authHeader == "" {
			authHeader = strings.TrimSpace(config.AuthHeader)
		}

		targetDir, err := filepath.Abs(addOptions.outputDir)
		if err != nil {
			return fmt.Errorf("resolve output directory: %w", err)
		}

		if err = apiClient.Skills.DownloadMarketplaceAndExtract(addOptions.skillID, authHeader, targetDir); err != nil {
			return err
		}

		fmt.Printf("Marketplace skill extracted to %s\n", targetDir)
		return nil
	},
}

func init() {
	AddCMD.Flags().StringVar(&addOptions.skillID, "skill-id", "", "Marketplace skill ID")
	AddCMD.Flags().StringVar(&addOptions.outputDir, "output", ".", "Directory where the archive will be extracted")
	AddCMD.Flags().StringVar(&addOptions.authHeader, "auth-header", "", "Authorization header value, e.g. 'Bearer <token>'")
}
