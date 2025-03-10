package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/geektheripper/go-gutils/git/git_utils"
	"github.com/geektheripper/go-gutils/git/virtual_repo"
	"github.com/geektheripper/gtf/internal/module"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dustin/go-humanize"
)

var publishCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"p"},
	Short:   "publish a module to remote repository",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := MustGetRepo()

		moduleName, modulePath := MustGetModuleNamePath(args)
		logger.Infof("Module Name: %s", moduleName)
		logger.Infof("Module Path: %s", modulePath)

		includes := viper.GetStringSlice("includes")
		if len(includes) > 0 {
			logger.Infof("Additional Files: %v", strings.Join(includes, ", "))
		}

		version := viper.GetString("version")

		upgradeType := ""
		if viper.GetBool("minor") {
			upgradeType = "minor"
		} else if viper.GetBool("major") {
			upgradeType = "major"
		} else if viper.GetBool("patch") {
			upgradeType = "patch"
		}

		force := viper.GetBool("force")
		if force {
			logger.Warnf("Force: true (override the existing tag if name conflicts)")
		}

		lrepo, err := git.PlainOpen(repoPath)
		if err != nil {
			logger.Fatalf("failed to load local repo: %v", err)
		}

		remote := viper.GetString("remote")
		if !strings.Contains(remote, ":") {
			_remote, err := lrepo.Remote(remote)
			if err != nil {
				logger.Fatalf("failed to get remote: %v", err)
			}
			remote = _remote.Config().URLs[0]
		}
		logger.Infof("Remote: %s", remote)

		if !git_utils.ValidateGitRemoteURL(remote) {
			logger.Fatalf("remote invalid")
		}

		vrepo, err := virtual_repo.NewVirtualRepo(remote, moduleName)
		if err != nil {
			logger.Fatalf("failed to create virtual repo: %v", err)
		}

		moduleMap, err := module.ResolveModules(vrepo)
		if err != nil {
			logger.Fatalf("failed to fetch packages form remote: %v", err)
		}

		pkg, ok := moduleMap[moduleName]

		// if package but try to upgrade
		if !ok && upgradeType != "" {
			logger.Fatalf("failed to apply %s, package not found in remote", upgradeType)
		}

		if ok {
			// if specified version already exists
			for _, v := range pkg.Versions {
				if v.String() == version {
					if force {
						logger.Warnf("version %s already exists, override it", version)
					} else {
						logger.Fatalf("version %s already exists", version)
					}
				}
			}

			if version == "" {
				version = pkg.NextVersion(upgradeType).String()
			}
		}

		if version == "" {
			version = "0.0.1"
		}

		tag := module.Tag(moduleName, version)
		logger.Infof("Tag: %s", tag)
		fmt.Print("\n")

		logger.Infof("collecting files for %s", tag)

		if !viper.GetBool("no-license") {
			_, err := vrepo.Import(
				filepath.Join(repoPath, "LICENSE"),
				"LICENSE",
			)
			if err != nil {
				logger.Warnf("warning: failed to copy license: %v", err)
			}
		}

		for _, file := range includes {
			_, err := vrepo.Import(
				filepath.Join(repoPath, file),
				file,
			)
			if err != nil {
				logger.Fatalf("failed to copy file: %v", err)
			}
		}

		report, err := vrepo.Import(modulePath, "")

		if err != nil {
			logger.Fatalf("failed to copy package files: %v", err)
		}

		logger.Infof("imported %d files, %s", report.Count, humanize.Bytes(report.Size))

		if force {
			err := vrepo.DeleteRemoteTag(tag)
			if err != nil {
				logger.Fatalf("failed to delete remote tag: %v", err)
			}
		}

		if err := vrepo.PushTag(tag, "Publish"); err != nil {
			logger.Fatalf("failed to publish: %v", err)
		}

		logger.Infof("published %s", tag)
	},
}

func init() {
	publishCmd.Flags().String("remote", "origin", "the remote to push to")
	viper.BindEnv("remote")

	publishCmd.Flags().StringP("version", "v", "", "publish a specific version")
	publishCmd.Flags().Bool("minor", false, "publish a minor version")
	publishCmd.Flags().Bool("major", false, "publish a major version")
	publishCmd.Flags().Bool("patch", false, "publish a patch version")
	publishCmd.MarkFlagsMutuallyExclusive("version", "minor", "major", "patch")

	publishCmd.Flags().BoolP("force", "f", false, "force publish (override the existing version)")

	publishCmd.Flags().Bool("no-license", false, "default copy license from root, use this to skip")
	publishCmd.Flags().StringArray("includes", []string{}, "include files to the package from root")

	viper.BindPFlags(publishCmd.Flags())

	moduleCmd.AddCommand(publishCmd)
}
