package main

import (
	"fmt"
	"time"

	"versioncontrolsystem_rp/engine"
	"versioncontrolsystem_rp/models"
)

type simpleIDGenerator struct {
	counter int
}

func (g *simpleIDGenerator) NewCommitID(tree map[string]string, parents []models.CommitID, t time.Time) models.CommitID {
	g.counter++
	return models.CommitID(fmt.Sprintf("commit-%d-%d", t.Unix(), g.counter))
}

type realClock struct{}

func (c *realClock) Now() time.Time {
	return time.Now()
}

func main() {
	repo := models.NewRepo()
	idGen := &simpleIDGenerator{}
	clock := &realClock{}
	eng := engine.NewInMemoryEngine(repo, idGen, clock)

	fmt.Println("=== Initializing repository ===")
	eng.AddFileToWorkingDirectory("file1.txt", "content1")
	eng.AddFileToWorkingDirectory("file2.txt", "content2")
	eng.AddFileToWorkingDirectory("src/main.go", "package main\n\nfunc main() {}")

	fmt.Println("\n=== Staging files ===")
	eng.Add("file1.txt")
	eng.Add("file2.txt")
	eng.Add("src/main.go")

	status, _ := eng.Status()
	fmt.Printf("Staged %d files\n", len(status.Files))

	fmt.Println("\n=== Making first commit ===")
	commitID1, _ := eng.Commit("Initial commit")
	fmt.Printf("Committed: %s\n", commitID1)

	fmt.Println("\n=== Adding more files and committing ===")
	eng.AddFileToWorkingDirectory("README.md", "# My Project\n\nA sample project.")
	eng.Add("README.md")
	commitID2, _ := eng.Commit("Add README")
	fmt.Printf("Committed: %s\n", commitID2)

	fmt.Println("\n=== Viewing log ===")
	log, _ := eng.Log()
	fmt.Printf("Commit history (%d commits):\n", len(log))
	for i, entry := range log {
		fmt.Printf("  %d. %s - %s\n", i+1, entry.ID, entry.Message)
	}

	fmt.Println("\n=== Creating and switching to feature branch ===")
	eng.CreateBranch("feature/new-feature")
	eng.CheckoutBranch("feature/new-feature")
	fmt.Println("Switched to branch: feature/new-feature")

	fmt.Println("\n=== Making changes on feature branch ===")
	eng.AddFileToWorkingDirectory("feature.go", "package main\n\nfunc newFeature() {}")
	eng.Add("feature.go")
	commitID3, _ := eng.Commit("Add new feature")
	fmt.Printf("Committed on feature branch: %s\n", commitID3)

	fmt.Println("\n=== Switching back to main ===")
	eng.CheckoutBranch("main")
	status, _ = eng.Status()
	fmt.Printf("On branch: main\n")
	fmt.Printf("Files in working directory: %d\n", len(status.Files))

	fmt.Println("\n=== Viewing final log ===")
	log, _ = eng.Log()
	fmt.Printf("Commit history on main (%d commits):\n", len(log))
	for i, entry := range log {
		fmt.Printf("  %d. %s - %s\n", i+1, entry.ID, entry.Message)
	}

	fmt.Println("\n=== Final status ===")
	status, _ = eng.Status()
	fmt.Printf("Head: %s\n", status.HeadDescription)
	fmt.Printf("Files: %d\n", len(status.Files))
	for _, file := range status.Files {
		statusStr := ""
		if file.Untracked {
			statusStr += "untracked "
		}
		if file.Staged {
			statusStr += "staged "
		}
		if file.Modified {
			statusStr += "modified "
		}
		if statusStr == "" {
			statusStr = "committed"
		}
		fmt.Printf("  %s: %s\n", file.Path, statusStr)
	}
}
