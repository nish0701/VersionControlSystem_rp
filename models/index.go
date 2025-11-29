package models

// Index represents the staging area, tracking which file contents are staged for commit.
type Index struct {
	entries map[string]string
}

// NewIndex constructs an empty Index.
func NewIndex() Index {
	return Index{
		entries: make(map[string]string),
	}
}

func (i *Index) GetEntries() map[string]string {
	return i.entries
}

func (i *Index) AddEntry(path string, content string) {
	if i.entries == nil {
		i.entries = make(map[string]string)
	}
	i.entries[path] = content
}

func (i *Index) RemoveEntry(path string) {
	delete(i.entries, path)
}

func (i *Index) Clear() {
	for k := range i.entries {
		delete(i.entries, k)
	}
}
