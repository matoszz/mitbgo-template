package version

import (
	"github.com/spf13/cobra"

	"github.com/datumforge/datum/pkg/utils/cli/useragent"

	template "github.com/datumforge/go-template/cmd/cli/cmd"
	"github.com/datumforge/go-template/internal/constants"
)

// VersionCmd is the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print template CLI version",
	Long:  `The template version command prints the version of the template CLI`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Println(constants.VerboseCLIVersion)
		cmd.Printf("User Agent: %s\n", useragent.GetUserAgent())
	},
}

func init() {
	template.RootCmd.AddCommand(versionCmd)
}
