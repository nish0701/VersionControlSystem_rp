package models

// BranchName is the name of a branch (e.g., "main", "feature/x").
type BranchName string

// Branch is a named reference to a commit.
type Branch struct {
	name   BranchName
	target CommitID
}

func (b *Branch) GetName() BranchName {
	return b.name
}

func (b *Branch) SetName(name BranchName) {
	b.name = name
}

func (b *Branch) GetTarget() CommitID {
	return b.target
}

func (b *Branch) SetTarget(target CommitID) {
	b.target = target
}

// DefaultBranchName is the branch created when a new repository is initialized.
const DefaultBranchName BranchName = "main"
