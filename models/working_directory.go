package models

// WorkingDirectory represents the in-memory view of editable files.
type WorkingDirectory struct {
	files map[string]string
}

// NewWorkingDirectory constructs an empty working directory.
func NewWorkingDirectory() WorkingDirectory {
	return WorkingDirectory{
		files: make(map[string]string),
	}
}

func (wd *WorkingDirectory) GetFiles() map[string]string {
	return wd.files
}

func (wd *WorkingDirectory) AddFile(path string, content string) {
	if wd.files == nil {
		wd.files = make(map[string]string)
	}
	wd.files[path] = content
}

func (wd *WorkingDirectory) RemoveFile(path string) {
	delete(wd.files, path)
}

func (wd *WorkingDirectory) Clear() {
	for k := range wd.files {
		delete(wd.files, k)
	}
}
