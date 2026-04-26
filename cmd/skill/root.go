package skill

import "github.com/spf13/cobra"

var RootCMD = &cobra.Command{
	Use:   "skill",
	Short: "Skill operations",
}

func init() {
	RootCMD.AddCommand(AddCMD)
}
