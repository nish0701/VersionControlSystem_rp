package models

import "time"

// CommitID uniquely identifies a commit in the repository.
type CommitID string

// Commit represents an immutable snapshot of the repository state at a point in time.
type Commit struct {
	id        CommitID
	parents   []CommitID
	message   string
	timestamp time.Time
	tree      map[string]string
}

func (c *Commit) GetID() CommitID {
	return c.id
}

func (c *Commit) SetID(id CommitID) {
	c.id = id
}

func (c *Commit) GetParents() []CommitID {
	return c.parents
}

func (c *Commit) SetParents(parents []CommitID) {
	c.parents = parents
}

func (c *Commit) GetMessage() string {
	return c.message
}

func (c *Commit) SetMessage(message string) {
	c.message = message
}

func (c *Commit) GetTimestamp() time.Time {
	return c.timestamp
}

func (c *Commit) SetTimestamp(timestamp time.Time) {
	c.timestamp = timestamp
}

func (c *Commit) GetTree() map[string]string {
	return c.tree
}

func (c *Commit) SetTree(tree map[string]string) {
	c.tree = tree
}
