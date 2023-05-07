package executor

import (
	"dh/internal/config"
)

type GitExecutor struct {
	BaseExecutor
}

func NewGitExecutor() *GitExecutor {
	res := &GitExecutor{
		newBaseExecutor("git"),
	}
	return res
}

func (r *GitExecutor) Checkout(issueBranch string) error {
	return r.execute("checkout", issueBranch)
}

func (r *GitExecutor) Stash(pop bool) error {
	arg := "stash"
	if pop {
		if !r.stashHasEntries() {
			config.WarnLog.Print("no stash entries in repo")
			return nil
		} else {
			return r.execute(arg, "pop")
		}
	} else {
		return r.execute(arg)
	}
}

func (r *GitExecutor) stashHasEntries() bool {
	err, stdout, _ := r.executeWithResult("stash", "list")
	if err != nil {
		config.ErrLog.Fatal(err.Error())
		return false
	}
	if stdout.Len() < 1 {
		return false
	}
	return true
}
