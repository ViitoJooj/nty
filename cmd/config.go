package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/ViitoJooj/nty/internal/tui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		langFlag, err := cmd.Flags().GetBool("lang")
		if err != nil {
			fmt.Println(err)
		}

		githubFlag, err := cmd.Flags().GetBool("github")
		if err != nil {
			fmt.Println(err)
		}

		containerEngineFlag, err := cmd.Flags().GetBool("container-engine")
		if err != nil {
			fmt.Println(err)
		}

		aiFlag, err := cmd.Flags().GetBool("ai")
		if err != nil {
			fmt.Println(err)
		}

		// config --github
		// this if will run the github config
		// the user need enter the github username and email
		// used for nty commit and push to github
		if githubFlag {
			application := tea.NewProgram(tui.NewApp())

			_, err = application.Run()
			if err != nil {
				fmt.Println(err)
			}
		}

		if langFlag {
			fmt.Println("Lang")
		}

		if containerEngineFlag {
			fmt.Println("container engine")
		}

		if aiFlag {
			fmt.Println("ai")
		}

		application := tea.NewProgram(tui.NewApp())
		_, err = application.Run()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	configCmd.Flags().Bool("github", false, "Configure github accout.")
	configCmd.Flags().Bool("lang", false, "Configure the language.")
	configCmd.Flags().Bool("container-engine", false, "Configure your container manager.")
	configCmd.Flags().Bool("ai", false, "Configure your artificial intelligence.")

	rootCmd.AddCommand(configCmd)
}
