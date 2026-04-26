package skill

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloud-ru/evo-ai-agents-skills-cli/internal/di"
	"github.com/spf13/cobra"
)

var addOptions struct {
	projectID  string
	skillID    string
	outputDir  string
	authHeader string
}

var AddCMD = &cobra.Command{
	Use:   "add",
	Short: "Download and extract a skill archive from BFF",
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

		projectID := strings.TrimSpace(addOptions.projectID)
		if projectID == "" {
			projectID = strings.TrimSpace(config.ProjectID)
		}
		if projectID == "" {
			return fmt.Errorf("project ID is required: pass --project-id or set PROJECT_ID")
		}
		authHeader := strings.TrimSpace(addOptions.authHeader)
		if authHeader == "" {
			authHeader = strings.TrimSpace(config.AuthHeader)
		}

		targetDir, err := filepath.Abs(addOptions.outputDir)
		if err != nil {
			return fmt.Errorf("resolve output directory: %w", err)
		}

		if err = apiClient.Skills.DownloadAndExtract(projectID, addOptions.skillID, authHeader, targetDir); err != nil {
			return err
		}

		fmt.Printf("Skill extracted to %s\n", targetDir)
		return nil
	},
}

func init() {
	AddCMD.Flags().StringVar(&addOptions.projectID, "project-id", "", "Project ID for scoped route")
	AddCMD.Flags().StringVar(&addOptions.skillID, "skill-id", "", "Skill ID")
	AddCMD.Flags().StringVar(&addOptions.outputDir, "output", ".", "Directory where the archive will be extracted")
	AddCMD.Flags().StringVar(&addOptions.authHeader, "auth-header", "", "Authorization header value, e.g. 'Bearer <token>'")
}
