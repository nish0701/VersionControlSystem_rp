package models

// Repo is the root aggregate for all version control state in memory.
type Repo struct {
	commits          map[CommitID]*Commit
	branches         map[BranchName]*Branch
	head             Head
	index            Index
	workingDirectory WorkingDirectory
}

func (r *Repo) GetCommits() map[CommitID]*Commit {
	return r.commits
}

func (r *Repo) AddCommit(id CommitID, commit *Commit) {
	if r.commits == nil {
		r.commits = make(map[CommitID]*Commit)
	}
	r.commits[id] = commit
}

func (r *Repo) GetCommit(id CommitID) (*Commit, bool) {
	commit, ok := r.commits[id]
	return commit, ok
}

func (r *Repo) GetBranches() map[BranchName]*Branch {
	return r.branches
}

func (r *Repo) AddBranch(name BranchName, branch *Branch) {
	if r.branches == nil {
		r.branches = make(map[BranchName]*Branch)
	}
	r.branches[name] = branch
}

func (r *Repo) GetBranch(name BranchName) (*Branch, bool) {
	branch, ok := r.branches[name]
	return branch, ok
}

func (r *Repo) GetHead() *Head {
	return &r.head
}

func (r *Repo) SetHead(head Head) {
	r.head = head
}

func (r *Repo) GetIndex() *Index {
	return &r.index
}

func (r *Repo) GetWorkingDirectory() *WorkingDirectory {
	return &r.workingDirectory
}

// NewRepo constructs a new in-memory repository with a default branch and HEAD.
func NewRepo() *Repo {
	r := &Repo{
		commits:          make(map[CommitID]*Commit),
		branches:         make(map[BranchName]*Branch),
		index:            NewIndex(),
		workingDirectory: NewWorkingDirectory(),
	}

	defaultBranch := &Branch{}
	defaultBranch.SetName(DefaultBranchName)
	defaultBranch.SetTarget("")
	r.AddBranch(DefaultBranchName, defaultBranch)

	defaultBranchName := DefaultBranchName
	r.head.SetBranchName(&defaultBranchName)
	r.head.SetCommitID(nil)

	return r
}
