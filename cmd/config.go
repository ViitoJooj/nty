package cmd

import (
	"fmt"

	"github.com/ViitoJooj/nty/internal/tui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		lang, _ := cmd.Flags().GetBool("lang")
		containerEngine, _ := cmd.Flags().GetBool("container-engine")
		ai, _ := cmd.Flags().GetBool("ai")

		// Placeholders até cada configuração ter sua tela.
		if lang {
			fmt.Println("Lang")
		}
		if containerEngine {
			fmt.Println("container engine")
		}
		if ai {
			fmt.Println("ai")
		}

		// --github abre o TUI; sem flags também cai no menu de config.
		if err := tui.Run(); err != nil {
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
