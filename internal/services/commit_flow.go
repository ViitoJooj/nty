package services

import "fmt"

func RunCommit() error {
	files, err := ChangedFiles()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("nada para commitar.")
		return nil
	}

	branch, err := resolveBranch(files)
	if err != nil {
		return err
	}

	for _, f := range files {
		diff, err := StageAndDiff(f)
		if err != nil {
			return fmt.Errorf("preparar %s: %w", f, err)
		}
		msg, err := CommitMessage(f, diff)
		if err != nil {
			return fmt.Errorf("gerar mensagem para %s: %w", f, err)
		}
		if err := Commit(msg); err != nil {
			return fmt.Errorf("commitar %s: %w", f, err)
		}
		fmt.Printf("  OK %s — %s\n", f, msg)
	}

	if err := Push(branch); err != nil {
		return fmt.Errorf("push: %w", err)
	}
	fmt.Println("push concluído.")
	return nil
}

func resolveBranch(files []string) (string, error) {
	current, err := CurrentBranch()
	if err != nil {
		return "", err
	}

	if !IsDefaultBranch(current) {
		fmt.Printf("branch (reaproveitada): %s\n", current)
		return current, nil
	}

	if deleted, err := CleanMergedBranches(); err != nil {
		fmt.Println("aviso ao limpar branches mergiadas:", err)
	} else if len(deleted) > 0 {
		fmt.Printf("branches mergiadas removidas: %v\n", deleted)
	}

	rawBranch, err := BranchName(files)
	if err != nil {
		return "", fmt.Errorf("gerar nome da branch: %w", err)
	}
	branch := SanitizeBranch(rawBranch)
	if err := CreateBranch(branch); err != nil {
		return "", fmt.Errorf("criar branch: %w", err)
	}
	fmt.Printf("branch: %s\n", branch)
	return branch, nil
}
