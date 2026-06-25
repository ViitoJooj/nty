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
		files, err := services.ChangedFiles()
		if err != nil {
			fmt.Println("erro:", err)
			return
		}
		if len(files) == 0 {
			fmt.Println("nada para commitar.")
			return
		}

		rawBranch, err := services.BranchName(files)
		if err != nil {
			fmt.Println("erro ao gerar nome da branch:", err)
			return
		}
		branch := services.SanitizeBranch(rawBranch)
		if err := services.CreateBranch(branch); err != nil {
			fmt.Println("erro ao criar branch:", err)
			return
		}
		fmt.Printf("branch: %s\n", branch)

		for _, f := range files {
			diff, err := services.StageAndDiff(f)
			if err != nil {
				fmt.Printf("erro ao preparar %s: %v\n", f, err)
				return
			}
			msg, err := services.CommitMessage(f, diff)
			if err != nil {
				fmt.Printf("erro ao gerar mensagem para %s: %v\n", f, err)
				return
			}
			if err := services.Commit(msg); err != nil {
				fmt.Printf("erro ao commitar %s: %v\n", f, err)
				return
			}
			fmt.Printf("  OK %s — %s\n", f, msg)
		}

		if err := services.Push(branch); err != nil {
			fmt.Println("erro no push:", err)
			return
		}
		fmt.Println("push concluído.")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
