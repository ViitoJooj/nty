package cmd

import (
	"fmt"

	"github.com/ViitoJooj/nty/internal/services"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Gera versão (IA), release notes e cria a release no GitHub. Dry-run sem --publish.",
	Run: func(cmd *cobra.Command, args []string) {
		publish, _ := cmd.Flags().GetBool("publish")
		if err := services.RunRelease(publish); err != nil {
			fmt.Println("erro:", err)
		}
	},
}

func init() {
	releaseCmd.Flags().Bool("publish", false, "Cria a release no GitHub (sem isso é dry-run).")
	rootCmd.AddCommand(releaseCmd)
}
