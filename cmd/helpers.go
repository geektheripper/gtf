package cmd

import (
	"path/filepath"

	"github.com/geektheripper/go-gutils/git/git_utils"
	"github.com/spf13/viper"
)

func MustGetRepo() string {
	repo := viper.GetString("repo")
	if repo == "" {
		logger.Fatalf("repo is not set")
	}

	ok, err := git_utils.IsGitRepo(repo)
	if err != nil {
		logger.Fatalf("failed to check if repo is a git repository: %v", err)
	}

	if !ok {
		logger.Fatalf("repo is not a git repository: %s", repo)
	}

	return repo
}

func MustGetModuleNamePath(args []string) (string, string) {
	packageName := args[0]
	packagePath := args[0]
	if len(args) > 1 {
		packagePath = args[1]
	}

	if packageName == "" {
		logger.Fatalf("package name is required")
	}

	if filepath.IsAbs(packageName) {
		logger.Fatalf("package name must be a relative path: %s", packageName)
	}

	if !filepath.IsAbs(packagePath) {
		packagePath = filepath.Join(MustGetRepo(), packagePath)
	}

	return filepath.Clean(packageName), packagePath
}
