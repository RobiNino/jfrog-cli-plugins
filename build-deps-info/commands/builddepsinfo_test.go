package builddepsinfo

import (
	"github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ServiceManagerMock struct {
	artifactory.EmptyArtifactoryServicesManager
}

func (smm *ServiceManagerMock) SearchFiles(params services.SearchParams) (*content.ContentReader, error) {
	cw, err := content.NewContentWriter(content.DefaultKey, true, false)
	if err != nil {
		return nil, err
	}
	defer cw.Close()
	item := utils.ResultItem{
		Actual_Sha1: "456",
		Properties: []utils.Property{
			{Key: "build.name", Value: "Build-Name"},
			{Key: "build.number", Value: "Build-Number"},
			{Key: "vcs.url", Value: "www.vcs.com"},
			{Key: "vcs.revision", Value: "248"},
		},
	}
	cw.Write(item)
	return content.NewContentReader(cw.GetFilePath(), cw.GetArrayKey()), nil
}

func TestGetDependenciesDetails(t *testing.T) {
	modules := []entities.Module{{Id: "my-plugin:", Artifacts: []entities.Artifact{
		{Name: "Artifact-name", Type: "Type", Checksum: entities.Checksum{Sha1: "123"}},
	}, Dependencies: []entities.Dependency{
		{Id: "Dependency", Type: "File", Checksum: entities.Checksum{Sha1: "456"}},
	}}}
	smMock := new(ServiceManagerMock)

	sha1ToBuildProps, err := getDependenciesDetails(modules, "repository", smMock)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []map[string]*DependencyProps{sha1ToBuildProps}, []map[string]*DependencyProps{GetFirstSearchResultSortedByAsc()})
}

func GetFirstSearchResultSortedByAsc() map[string]*DependencyProps {
	return map[string]*DependencyProps{
		"456": {Build: "Build-Name/Build-Number", Vcs: entities.Vcs{Url: "www.vcs.com", Revision: "248"}},
	}
}
