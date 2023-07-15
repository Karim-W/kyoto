package cli

import (
	"fmt"
	"os"

	"github.com/karim-w/kyoto/spotify"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kyoto",
	Short: "kyoto is a CLI that interacts with the Spotify API",
	Long: `kyoto is a CLI that interacts with the Spotify API. Allowing you to
		- Get the currently playing song
	`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Login to Spotify",
		Long:  `Login to Spotify`,
		Run: func(cmd *cobra.Command, args []string) {
			spotify.StartAuth()
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
