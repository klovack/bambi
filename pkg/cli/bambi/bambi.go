package bambi

import (
	"os"
	"path/filepath"

	"github.com/klovack/bambi/pkg/command"
	"github.com/klovack/bambi/pkg/util"
	"github.com/spf13/cobra"
)

type Project struct {
	Name       string
	Path       string
	Folders    map[string]bool
	IsFlat     bool
	InitGit    bool
	ModuleRepo string
}

type foldersB struct {
	API     bool
	Assets  bool
	Cmd     bool
	Configs bool
	Docs    bool
	Pkg     bool
	Test    bool
}

// NewCommand initialize bambi command
func NewCommand() *cobra.Command {
	fb := &foldersB{}

	project := &Project{}

	c := &cobra.Command{
		Use:   "bambi",
		Short: "Add structure to your go project",
		Long: `Description:
    Focus more on writing program rather than structuring your project folder.
    This command creates project structures for Go program.
    `,
		Example: "bambi .",
		Args:    cobra.ExactArgs(1),
		Run: func(ccmd *cobra.Command, args []string) {
			project.Path = args[0]

			if !project.IsFlat {
				project.Path = filepath.Join(args[0], project.Name)
			}

			project.Folders = map[string]bool{
				"api":     fb.API,
				"assets":  fb.Assets,
				"cmd":     fb.Cmd,
				"configs": fb.Configs,
				"docs":    fb.Docs,
				"pkg":     fb.Pkg,
				"test":    fb.Test,
			}

			// Initialize necessary folders
			initFolders(project)

			err := os.Chdir(project.Path)
			util.CheckErrorP(err)

			// Initialize git repository. It panics if git is not installed or
			// doesn't have permission to modify folder
			initGit(project)

			// Initialize go module. It panics if go mod init can't be executed.
			// It will however check whether the user has entered go module repo
			initGoModule(project)
		},
	}

	c.Flags().StringVarP(&project.Name, "name", "n", "my-app", "Give name to the project")
	c.Flags().BoolVarP(&project.IsFlat, "flat", "f", false, "Don't create subfolder for the project")

	c.Flags().StringVarP(&project.ModuleRepo, "go-mod", "m", "", "Initialize go mod repo if the go version > 1.11")

	// Git
	c.Flags().BoolVarP(&project.InitGit, "git", "g", true, "Initialize git repository")

	// Folders
	c.Flags().BoolVar(&fb.API, "api", true, "Tell the program to also create api folder")
	c.Flags().BoolVar(&fb.Assets, "assets", true, "Tell the program to also create assets folder")
	c.Flags().BoolVar(&fb.Cmd, "cmd", true, "Tell the program to also create cmd folder")
	c.Flags().BoolVar(&fb.Configs, "configs", true, "Tell the program to also create configs folder")
	c.Flags().BoolVar(&fb.Docs, "docs", true, "Tell the program to also create docs folder")
	c.Flags().BoolVar(&fb.Pkg, "pkg", true, "Tell the program to also create pkg folder")
	c.Flags().BoolVar(&fb.Test, "test", true, "Tell the program to also create test folder")

	// Enter sub command here
	c.AddCommand()

	return c
}

func initFolders(project *Project) {
	err := os.MkdirAll(project.Path, os.ModePerm)
	util.CheckErrorP(err)

	// Create structure folders
	for dir := range project.Folders {
		if project.Folders[dir] {
			err = os.MkdirAll(filepath.Join(project.Path, dir), os.ModePerm)
			util.CheckErrorP(err)
			if dir == "cmd" {
				err = os.MkdirAll(filepath.Join(project.Path, dir, project.Name), os.ModePerm)
			}
			util.CheckErrorP(err)
		}
	}
}

func initGit(project *Project) {
	if project.InitGit {
		if !command.IsAvailable("git") {
			util.Exit("Please install git to initialize git repository within the project")
		}

		err := command.Execute("git", "init")
		util.CheckErrorP(err)

		_, err = os.Create(".gitignore")
		util.CheckErrorP(err)
	}
}

func initGoModule(project *Project) {
	if util.HasGoModule() && project.ModuleRepo != "" {
		err := command.Execute("go", "mod", "init", project.ModuleRepo)
		util.CheckErrorP(err)

		_, err = os.Create(".gitignore")
		util.CheckErrorP(err)
	}
}
