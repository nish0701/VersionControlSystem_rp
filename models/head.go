package models

// Head represents the current position in the commit graph.
type Head struct {
	branchName *BranchName
	commitID   *CommitID
}

func (h *Head) GetBranchName() *BranchName {
	return h.branchName
}

func (h *Head) SetBranchName(name *BranchName) {
	h.branchName = name
}

func (h *Head) GetCommitID() *CommitID {
	return h.commitID
}

func (h *Head) SetCommitID(id *CommitID) {
	h.commitID = id
}
