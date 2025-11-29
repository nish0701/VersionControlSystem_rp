package engine

import (
	"errors"
	"regexp"
	"time"

	"versioncontrolsystem_rp/models"
)

type FileStatus struct {
	Path      string
	Staged    bool
	Modified  bool
	Untracked bool
}

type StatusResult struct {
	HeadDescription string
	Files           []FileStatus
}

type LogEntry struct {
	ID        models.CommitID
	Message   string
	Timestamp time.Time
	Parents   []models.CommitID
}

type Engine interface {
	AddFileToWorkingDirectory(path string, content string)
	Add(pattern string) error
	Commit(message string) (models.CommitID, error)
	CheckoutBranch(name models.BranchName) error
	CreateBranch(name models.BranchName) error
	Status() (StatusResult, error)
	Log() ([]LogEntry, error)
	Reset(commitID models.CommitID) error
}

type CommitIDGenerator interface {
	NewCommitID(tree map[string]string, parents []models.CommitID, t time.Time) models.CommitID
}

type Clock interface {
	Now() time.Time
}

type InMemoryEngine struct {
	repo  *models.Repo
	idGen CommitIDGenerator
	clock Clock
}

const (
	errMsgNoSuchBranch      = "no such branch"
	errMsgNoSuchCommit      = "no such commit"
	errMsgNothingToCommit   = "nothing to commit"
	errMsgNoMatchesForAdd   = "no files matched add pattern"
	errMsgBranchExists      = "branch already exists"
	errMsgNoCurrentCommit   = "no current commit to base branch on"
	errMsgResetDetachedHead = "cannot reset in detached HEAD state"
)

func NewInMemoryEngine(repo *models.Repo, idGen CommitIDGenerator, clock Clock) *InMemoryEngine {
	return &InMemoryEngine{
		repo:  repo,
		idGen: idGen,
		clock: clock,
	}
}

func (e *InMemoryEngine) AddFileToWorkingDirectory(path string, content string) {
	repo := e.repo
	if repo == nil {
		return
	}
	wd := repo.GetWorkingDirectory()
	wd.AddFile(path, content)
}

func (e *InMemoryEngine) Add(pattern string) error {
	repo := e.repo
	if repo == nil {
		return nil
	}
	wd := repo.GetWorkingDirectory()
	index := repo.GetIndex()
	wdFiles := wd.GetFiles()
	matched := false
	if content, ok := wdFiles[pattern]; ok {
		index.AddEntry(pattern, content)
		matched = true
	} else {
		// for the suupport of git add *.go* etc.
		r, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		for path, content := range wdFiles {
			if r.MatchString(path) {
				index.AddEntry(path, content)
				matched = true
			}
		}
	}
	if !matched {
		return errors.New(errMsgNoMatchesForAdd)
	}
	return nil
}

func (e *InMemoryEngine) Commit(message string) (models.CommitID, error) {
	repo := e.repo
	if repo == nil {
		return "", nil
	}
	index := repo.GetIndex()
	indexEntries := index.GetEntries()
	if len(indexEntries) == 0 {
		return "", errors.New(errMsgNothingToCommit)
	}
	parentID, parentTree := currentCommit(repo)
	tree := make(map[string]string, len(parentTree)+len(indexEntries))
	for k, v := range parentTree {
		tree[k] = v
	}
	for path, content := range indexEntries {
		tree[path] = content
	}
	now := e.clock.Now()
	var parents []models.CommitID
	if parentID != "" {
		parents = []models.CommitID{parentID}
	}
	id := e.idGen.NewCommitID(tree, parents, now)
	commit := &models.Commit{}
	commit.SetID(id)
	commit.SetParents(parents)
	commit.SetMessage(message)
	commit.SetTimestamp(now)
	commit.SetTree(tree)
	repo.AddCommit(id, commit)
	index.Clear()
	head := repo.GetHead()
	if head.GetBranchName() != nil {
		branch, ok := repo.GetBranch(*head.GetBranchName())
		if ok {
			branch.SetTarget(id)
		}
	}
	return id, nil
}

func (e *InMemoryEngine) CheckoutBranch(name models.BranchName) error {
	repo := e.repo
	if repo == nil {
		return nil
	}
	branch, ok := repo.GetBranch(name)
	if !ok {
		return errors.New(errMsgNoSuchBranch)
	}
	applyCommitToWorkingDirectory(repo, branch.GetTarget())
	nameCopy := name
	head := repo.GetHead()
	head.SetBranchName(&nameCopy)
	head.SetCommitID(nil)
	repo.GetIndex().Clear()
	return nil
}

func (e *InMemoryEngine) CreateBranch(name models.BranchName) error {
	repo := e.repo
	if repo == nil {
		return nil
	}
	if _, exists := repo.GetBranch(name); exists {
		return errors.New(errMsgBranchExists)
	}
	id, _ := currentCommit(repo)
	if id == "" {
		return errors.New(errMsgNoCurrentCommit)
	}
	branch := &models.Branch{}
	branch.SetName(name)
	branch.SetTarget(id)
	repo.AddBranch(name, branch)
	return nil
}

func (e *InMemoryEngine) Status() (StatusResult, error) {
	repo := e.repo
	if repo == nil {
		return StatusResult{}, nil
	}
	headDescription, baseCommit := describeHead(repo)
	commitTree := map[string]string{}
	if baseCommit != nil {
		for k, v := range baseCommit.GetTree() {
			commitTree[k] = v
		}
	}
	pathsSet := map[string]struct{}{}
	for path := range commitTree {
		pathsSet[path] = struct{}{}
	}
	indexEntries := repo.GetIndex().GetEntries()
	for path := range indexEntries {
		pathsSet[path] = struct{}{}
	}
	wdFiles := repo.GetWorkingDirectory().GetFiles()
	for path := range wdFiles {
		pathsSet[path] = struct{}{}
	}
	var paths []string
	for p := range pathsSet {
		paths = append(paths, p)
	}
	var files []FileStatus
	for _, path := range paths {
		wdContent, inWD := wdFiles[path]
		indexContent, inIndex := indexEntries[path]
		commitContent, inCommit := commitTree[path]
		staged := inIndex
		untracked := inWD && !inCommit && !inIndex
		modified := false
		if inIndex {
			if inWD && wdContent != indexContent {
				modified = true
			}
		} else if inWD {
			if !inCommit || commitContent != wdContent {
				modified = true
			}
		}
		files = append(files, FileStatus{
			Path:      path,
			Staged:    staged,
			Modified:  modified,
			Untracked: untracked,
		})
	}
	return StatusResult{
		HeadDescription: headDescription,
		Files:           files,
	}, nil
}

func (e *InMemoryEngine) Log() ([]LogEntry, error) {
	repo := e.repo
	if repo == nil {
		return nil, nil
	}
	_, baseCommit := describeHead(repo)
	if baseCommit == nil {
		return nil, nil
	}
	var entries []LogEntry
	current := baseCommit
	for current != nil {
		entry := LogEntry{
			ID:        current.GetID(),
			Message:   current.GetMessage(),
			Timestamp: current.GetTimestamp(),
			Parents:   append([]models.CommitID(nil), current.GetParents()...),
		}
		entries = append(entries, entry)
		parents := current.GetParents()
		if len(parents) == 0 {
			break
		}
		next, ok := repo.GetCommit(parents[0])
		if !ok {
			break
		}
		current = next
	}
	return entries, nil
}

func (e *InMemoryEngine) Reset(commitID models.CommitID) error {
	repo := e.repo
	if repo == nil {
		return nil
	}
	head := repo.GetHead()
	if head.GetBranchName() == nil {
		return errors.New(errMsgResetDetachedHead)
	}
	if _, ok := repo.GetCommit(commitID); !ok {
		return errors.New(errMsgNoSuchCommit)
	}
	branch, _ := repo.GetBranch(*head.GetBranchName())
	branch.SetTarget(commitID)
	applyCommitToWorkingDirectory(repo, commitID)
	repo.GetIndex().Clear()
	return nil
}

func currentCommit(repo *models.Repo) (models.CommitID, map[string]string) {
	head := repo.GetHead()
	if head.GetBranchName() == nil {
		return "", map[string]string{}
	}
	branch, ok := repo.GetBranch(*head.GetBranchName())
	if !ok || branch.GetTarget() == "" {
		return "", map[string]string{}
	}
	commit, ok := repo.GetCommit(branch.GetTarget())
	if !ok {
		return "", map[string]string{}
	}
	return commit.GetID(), commit.GetTree()
}

func applyCommitToWorkingDirectory(repo *models.Repo, id models.CommitID) {
	wd := repo.GetWorkingDirectory()
	wd.Clear()
	if id == "" {
		return
	}
	commit, ok := repo.GetCommit(id)
	if !ok {
		return
	}
	for path, content := range commit.GetTree() {
		wd.AddFile(path, content)
	}
}

func describeHead(repo *models.Repo) (string, *models.Commit) {
	head := repo.GetHead()
	if head.GetBranchName() == nil {
		return "", nil
	}
	branchName := head.GetBranchName()
	branch, ok := repo.GetBranch(*branchName)
	if !ok || branch.GetTarget() == "" {
		return string(*branchName), nil
	}
	commit, ok := repo.GetCommit(branch.GetTarget())
	if !ok {
		return string(*branchName), nil
	}
	return string(*branchName), commit
}
