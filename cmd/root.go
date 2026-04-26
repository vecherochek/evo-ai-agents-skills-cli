package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-skills-cli/cmd/auth"
	"github.com/cloud-ru/evo-ai-agents-skills-cli/cmd/skill"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

var RootCMD = &cobra.Command{
	Use:   "skills-cli",
	Short: "CLI for managing AI assistant skills",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		logger := log.New(os.Stderr)
		logger.SetReportTimestamp(true)
		logger.SetReportCaller(true)
		if verbose {
			logger.SetLevel(log.DebugLevel)
		} else {
			logger.SetLevel(log.InfoLevel)
		}
		log.SetDefault(logger)
	},
}

func init() {
	RootCMD.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Enable verbose logs")
	RootCMD.AddCommand(auth.RootCMD)
	RootCMD.AddCommand(skill.RootCMD)
}
