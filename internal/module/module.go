package module

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/geektheripper/go-gutils/git/virtual_repo"
)

type Module struct {
	Name     string
	Versions []*semver.Version
}

func (m *Module) GetLatestVersion() *semver.Version {
	sort.Sort(semver.Collection(m.Versions))
	return m.Versions[len(m.Versions)-1]
}

func (m *Module) NextVersion(upgradeType ...string) *semver.Version {
	if len(upgradeType) > 1 {
		panic("only one upgrade type is allowed")
	}

	latestVersion := m.GetLatestVersion()

	nextVersion := latestVersion.IncPatch()

	switch upgradeType[0] {
	case "patch":
		nextVersion = latestVersion.IncPatch()
	case "minor":
		nextVersion = latestVersion.IncMinor()
	case "major":
		nextVersion = latestVersion.IncMajor()
	}

	return &nextVersion
}

var moduleRefRegex = regexp.MustCompile(`^refs/tags/(?P<module>[^/]*)@v(?P<version>.*)$`)

func ResolveModules(v *virtual_repo.VirtualRepo) (map[string]*Module, error) {
	refs, err := v.FilterRefs("refs/tags/")
	if err != nil {
		return nil, err
	}

	modules := map[string]*Module{}
	for _, ref := range refs {
		tag := ref.Name().String()
		matches := moduleRefRegex.FindStringSubmatch(tag)

		if len(matches) == 0 {
			continue
		}

		moduleName := matches[moduleRefRegex.SubexpIndex("module")]
		versionText := matches[moduleRefRegex.SubexpIndex("version")]

		if _, ok := modules[moduleName]; !ok {
			modules[moduleName] = &Module{
				Name:     moduleName,
				Versions: []*semver.Version{},
			}
		}

		if version, err := semver.NewVersion(versionText); err != nil {
			continue
		} else {
			modules[moduleName].Versions = append(modules[moduleName].Versions, version)
		}
	}

	return modules, nil

}

func Tag(moduleName string, version string) string {
	return fmt.Sprintf("%s@v%s", moduleName, version)
}
