package auth

import "github.com/spf13/cobra"

var RootCMD = &cobra.Command{
	Use:   "auth",
	Short: "Authentication management",
	Long: `Manage authentication for ai-agents-skills-cli.

Available commands:
  login   - Save IAM credentials
  logout  - Remove saved credentials
  status  - Show current authentication status
  config  - Show active auth configuration`,
}

func init() {
	RootCMD.AddCommand(loginCmd)
	RootCMD.AddCommand(logoutCmd)
	RootCMD.AddCommand(statusCmd)
	RootCMD.AddCommand(configCmd)
}
