package executor

type GitExecutor struct {
	ProgramExecutor
}

func NewGitExecutor() *GitExecutor {
	return &GitExecutor{
		newFileExecutor("git"),
	}
}

func (r *GitExecutor) Checkout(issueBranch string) error {
	return r.execute("checkout", issueBranch)
}

func (r *GitExecutor) Stash(pop bool) error {
	arg := "stash"
	if pop {
		if !r.stashHasEntries() {
			logger.Print("no stash entries in repo")
			return nil
		} else {
			return r.execute(arg, "pop")
		}
	} else {
		return r.execute(arg)
	}
}

func (r *GitExecutor) stashHasEntries() bool {
	stdout, err := r.executeWithResult("stash", "list")
	if err != nil {
		logger.Fatal(err.Error())
		return false
	}
	if stdout.Len() < 1 {
		return false
	}
	return true
}
