package cmd

import (
	"fmt"

	"github.com/ViitoJooj/nty/internal/services"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Cria branch, commita arquivo por arquivo (mensagens via IA) e dá push.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := services.RunCommit(); err != nil {
			fmt.Println("erro:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
