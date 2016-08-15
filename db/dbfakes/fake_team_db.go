// This file was generated by counterfeiter
package dbfakes

import (
	"sync"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

type FakeTeamDB struct {
	GetPipelinesStub        func() ([]db.SavedPipeline, error)
	getPipelinesMutex       sync.RWMutex
	getPipelinesArgsForCall []struct{}
	getPipelinesReturns     struct {
		result1 []db.SavedPipeline
		result2 error
	}
	GetPublicPipelinesStub        func() ([]db.SavedPipeline, error)
	getPublicPipelinesMutex       sync.RWMutex
	getPublicPipelinesArgsForCall []struct{}
	getPublicPipelinesReturns     struct {
		result1 []db.SavedPipeline
		result2 error
	}
	GetPrivateAndAllPublicPipelinesStub        func() ([]db.SavedPipeline, error)
	getPrivateAndAllPublicPipelinesMutex       sync.RWMutex
	getPrivateAndAllPublicPipelinesArgsForCall []struct{}
	getPrivateAndAllPublicPipelinesReturns     struct {
		result1 []db.SavedPipeline
		result2 error
	}
	GetPipelineByNameStub        func(pipelineName string) (db.SavedPipeline, bool, error)
	getPipelineByNameMutex       sync.RWMutex
	getPipelineByNameArgsForCall []struct {
		pipelineName string
	}
	getPipelineByNameReturns struct {
		result1 db.SavedPipeline
		result2 bool
		result3 error
	}
	GetPublicPipelineByNameStub        func(pipelineName string) (db.SavedPipeline, bool, error)
	getPublicPipelineByNameMutex       sync.RWMutex
	getPublicPipelineByNameArgsForCall []struct {
		pipelineName string
	}
	getPublicPipelineByNameReturns struct {
		result1 db.SavedPipeline
		result2 bool
		result3 error
	}
	OrderPipelinesStub        func([]string) error
	orderPipelinesMutex       sync.RWMutex
	orderPipelinesArgsForCall []struct {
		arg1 []string
	}
	orderPipelinesReturns struct {
		result1 error
	}
	GetTeamStub        func() (db.SavedTeam, bool, error)
	getTeamMutex       sync.RWMutex
	getTeamArgsForCall []struct{}
	getTeamReturns     struct {
		result1 db.SavedTeam
		result2 bool
		result3 error
	}
	UpdateBasicAuthStub        func(basicAuth *db.BasicAuth) (db.SavedTeam, error)
	updateBasicAuthMutex       sync.RWMutex
	updateBasicAuthArgsForCall []struct {
		basicAuth *db.BasicAuth
	}
	updateBasicAuthReturns struct {
		result1 db.SavedTeam
		result2 error
	}
	UpdateGitHubAuthStub        func(gitHubAuth *db.GitHubAuth) (db.SavedTeam, error)
	updateGitHubAuthMutex       sync.RWMutex
	updateGitHubAuthArgsForCall []struct {
		gitHubAuth *db.GitHubAuth
	}
	updateGitHubAuthReturns struct {
		result1 db.SavedTeam
		result2 error
	}
	UpdateUAAAuthStub        func(uaaAuth *db.UAAAuth) (db.SavedTeam, error)
	updateUAAAuthMutex       sync.RWMutex
	updateUAAAuthArgsForCall []struct {
		uaaAuth *db.UAAAuth
	}
	updateUAAAuthReturns struct {
		result1 db.SavedTeam
		result2 error
	}
	UpdateGenericOAuthStub        func(genericOAuth *db.GenericOAuth) (db.SavedTeam, error)
	updateGenericOAuthMutex       sync.RWMutex
	updateGenericOAuthArgsForCall []struct {
		genericOAuth *db.GenericOAuth
	}
	updateGenericOAuthReturns struct {
		result1 db.SavedTeam
		result2 error
	}

	GetConfigStub        func(pipelineName string) (atc.Config, atc.RawConfig, db.ConfigVersion, error)
	getConfigMutex       sync.RWMutex
	getConfigArgsForCall []struct {
		pipelineName string
	}
	getConfigReturns struct {
		result1 atc.Config
		result2 atc.RawConfig
		result3 db.ConfigVersion
		result4 error
	}
	SaveConfigStub        func(string, atc.Config, db.ConfigVersion, db.PipelinePausedState) (db.SavedPipeline, bool, error)
	saveConfigMutex       sync.RWMutex
	saveConfigArgsForCall []struct {
		arg1 string
		arg2 atc.Config
		arg3 db.ConfigVersion
		arg4 db.PipelinePausedState
	}
	saveConfigReturns struct {
		result1 db.SavedPipeline
		result2 bool
		result3 error
	}
	CreateOneOffBuildStub        func() (db.Build, error)
	createOneOffBuildMutex       sync.RWMutex
	createOneOffBuildArgsForCall []struct{}
	createOneOffBuildReturns     struct {
		result1 db.Build
		result2 error
	}
	GetBuildsStub        func(page db.Page, publicOnly bool) ([]db.Build, db.Pagination, error)
	getBuildsMutex       sync.RWMutex
	getBuildsArgsForCall []struct {
		page       db.Page
		publicOnly bool
	}
	getBuildsReturns struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}
	GetBuildStub        func(buildID int) (db.Build, bool, error)
	getBuildMutex       sync.RWMutex
	getBuildArgsForCall []struct {
		buildID int
	}
	getBuildReturns struct {
		result1 db.Build
		result2 bool
		result3 error
	}
	WorkersStub        func() ([]db.SavedWorker, error)
	workersMutex       sync.RWMutex
	workersArgsForCall []struct{}
	workersReturns     struct {
		result1 []db.SavedWorker
		result2 error
	}
	GetContainerStub        func(handle string) (db.SavedContainer, bool, error)
	getContainerMutex       sync.RWMutex
	getContainerArgsForCall []struct {
		handle string
	}
	getContainerReturns struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}
	FindContainersByDescriptorsStub        func(id db.Container) ([]db.SavedContainer, error)
	findContainersByDescriptorsMutex       sync.RWMutex
	findContainersByDescriptorsArgsForCall []struct {
		id db.Container
	}
	findContainersByDescriptorsReturns struct {
		result1 []db.SavedContainer
		result2 error
	}
	GetVolumesStub        func() ([]db.SavedVolume, error)
	getVolumesMutex       sync.RWMutex
	getVolumesArgsForCall []struct{}
	getVolumesReturns     struct {
		result1 []db.SavedVolume
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTeamDB) GetPipelines() ([]db.SavedPipeline, error) {
	fake.getPipelinesMutex.Lock()
	fake.getPipelinesArgsForCall = append(fake.getPipelinesArgsForCall, struct{}{})
	fake.recordInvocation("GetPipelines", []interface{}{})
	fake.getPipelinesMutex.Unlock()
	if fake.GetPipelinesStub != nil {
		return fake.GetPipelinesStub()
	} else {
		return fake.getPipelinesReturns.result1, fake.getPipelinesReturns.result2
	}
}

func (fake *FakeTeamDB) GetPipelinesCallCount() int {
	fake.getPipelinesMutex.RLock()
	defer fake.getPipelinesMutex.RUnlock()
	return len(fake.getPipelinesArgsForCall)
}

func (fake *FakeTeamDB) GetPipelinesReturns(result1 []db.SavedPipeline, result2 error) {
	fake.GetPipelinesStub = nil
	fake.getPipelinesReturns = struct {
		result1 []db.SavedPipeline
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) GetPublicPipelines() ([]db.SavedPipeline, error) {
	fake.getPublicPipelinesMutex.Lock()
	fake.getPublicPipelinesArgsForCall = append(fake.getPublicPipelinesArgsForCall, struct{}{})
	fake.recordInvocation("GetPublicPipelines", []interface{}{})
	fake.getPublicPipelinesMutex.Unlock()
	if fake.GetPublicPipelinesStub != nil {
		return fake.GetPublicPipelinesStub()
	} else {
		return fake.getPublicPipelinesReturns.result1, fake.getPublicPipelinesReturns.result2
	}
}

func (fake *FakeTeamDB) GetPublicPipelinesCallCount() int {
	fake.getPublicPipelinesMutex.RLock()
	defer fake.getPublicPipelinesMutex.RUnlock()
	return len(fake.getPublicPipelinesArgsForCall)
}

func (fake *FakeTeamDB) GetPublicPipelinesReturns(result1 []db.SavedPipeline, result2 error) {
	fake.GetPublicPipelinesStub = nil
	fake.getPublicPipelinesReturns = struct {
		result1 []db.SavedPipeline
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) GetPrivateAndAllPublicPipelines() ([]db.SavedPipeline, error) {
	fake.getPrivateAndAllPublicPipelinesMutex.Lock()
	fake.getPrivateAndAllPublicPipelinesArgsForCall = append(fake.getPrivateAndAllPublicPipelinesArgsForCall, struct{}{})
	fake.recordInvocation("GetPrivateAndAllPublicPipelines", []interface{}{})
	fake.getPrivateAndAllPublicPipelinesMutex.Unlock()
	if fake.GetPrivateAndAllPublicPipelinesStub != nil {
		return fake.GetPrivateAndAllPublicPipelinesStub()
	} else {
		return fake.getPrivateAndAllPublicPipelinesReturns.result1, fake.getPrivateAndAllPublicPipelinesReturns.result2
	}
}

func (fake *FakeTeamDB) GetPrivateAndAllPublicPipelinesCallCount() int {
	fake.getPrivateAndAllPublicPipelinesMutex.RLock()
	defer fake.getPrivateAndAllPublicPipelinesMutex.RUnlock()
	return len(fake.getPrivateAndAllPublicPipelinesArgsForCall)
}

func (fake *FakeTeamDB) GetPrivateAndAllPublicPipelinesReturns(result1 []db.SavedPipeline, result2 error) {
	fake.GetPrivateAndAllPublicPipelinesStub = nil
	fake.getPrivateAndAllPublicPipelinesReturns = struct {
		result1 []db.SavedPipeline
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) GetPipelineByName(pipelineName string) (db.SavedPipeline, bool, error) {
	fake.getPipelineByNameMutex.Lock()
	fake.getPipelineByNameArgsForCall = append(fake.getPipelineByNameArgsForCall, struct {
		pipelineName string
	}{pipelineName})
	fake.recordInvocation("GetPipelineByName", []interface{}{pipelineName})
	fake.getPipelineByNameMutex.Unlock()
	if fake.GetPipelineByNameStub != nil {
		return fake.GetPipelineByNameStub(pipelineName)
	} else {
		return fake.getPipelineByNameReturns.result1, fake.getPipelineByNameReturns.result2, fake.getPipelineByNameReturns.result3
	}
}

func (fake *FakeTeamDB) GetPipelineByNameCallCount() int {
	fake.getPipelineByNameMutex.RLock()
	defer fake.getPipelineByNameMutex.RUnlock()
	return len(fake.getPipelineByNameArgsForCall)
}

func (fake *FakeTeamDB) GetPipelineByNameArgsForCall(i int) string {
	fake.getPipelineByNameMutex.RLock()
	defer fake.getPipelineByNameMutex.RUnlock()
	return fake.getPipelineByNameArgsForCall[i].pipelineName
}

func (fake *FakeTeamDB) GetPipelineByNameReturns(result1 db.SavedPipeline, result2 bool, result3 error) {
	fake.GetPipelineByNameStub = nil
	fake.getPipelineByNameReturns = struct {
		result1 db.SavedPipeline
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamDB) GetPublicPipelineByName(pipelineName string) (db.SavedPipeline, bool, error) {
	fake.getPublicPipelineByNameMutex.Lock()
	fake.getPublicPipelineByNameArgsForCall = append(fake.getPublicPipelineByNameArgsForCall, struct {
		pipelineName string
	}{pipelineName})
	fake.recordInvocation("GetPublicPipelineByName", []interface{}{pipelineName})
	fake.getPublicPipelineByNameMutex.Unlock()
	if fake.GetPublicPipelineByNameStub != nil {
		return fake.GetPublicPipelineByNameStub(pipelineName)
	} else {
		return fake.getPublicPipelineByNameReturns.result1, fake.getPublicPipelineByNameReturns.result2, fake.getPublicPipelineByNameReturns.result3
	}
}

func (fake *FakeTeamDB) GetPublicPipelineByNameCallCount() int {
	fake.getPublicPipelineByNameMutex.RLock()
	defer fake.getPublicPipelineByNameMutex.RUnlock()
	return len(fake.getPublicPipelineByNameArgsForCall)
}

func (fake *FakeTeamDB) GetPublicPipelineByNameArgsForCall(i int) string {
	fake.getPublicPipelineByNameMutex.RLock()
	defer fake.getPublicPipelineByNameMutex.RUnlock()
	return fake.getPublicPipelineByNameArgsForCall[i].pipelineName
}

func (fake *FakeTeamDB) GetPublicPipelineByNameReturns(result1 db.SavedPipeline, result2 bool, result3 error) {
	fake.GetPublicPipelineByNameStub = nil
	fake.getPublicPipelineByNameReturns = struct {
		result1 db.SavedPipeline
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamDB) OrderPipelines(arg1 []string) error {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.orderPipelinesMutex.Lock()
	fake.orderPipelinesArgsForCall = append(fake.orderPipelinesArgsForCall, struct {
		arg1 []string
	}{arg1Copy})
	fake.recordInvocation("OrderPipelines", []interface{}{arg1Copy})
	fake.orderPipelinesMutex.Unlock()
	if fake.OrderPipelinesStub != nil {
		return fake.OrderPipelinesStub(arg1)
	} else {
		return fake.orderPipelinesReturns.result1
	}
}

func (fake *FakeTeamDB) OrderPipelinesCallCount() int {
	fake.orderPipelinesMutex.RLock()
	defer fake.orderPipelinesMutex.RUnlock()
	return len(fake.orderPipelinesArgsForCall)
}

func (fake *FakeTeamDB) OrderPipelinesArgsForCall(i int) []string {
	fake.orderPipelinesMutex.RLock()
	defer fake.orderPipelinesMutex.RUnlock()
	return fake.orderPipelinesArgsForCall[i].arg1
}

func (fake *FakeTeamDB) OrderPipelinesReturns(result1 error) {
	fake.OrderPipelinesStub = nil
	fake.orderPipelinesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTeamDB) GetTeam() (db.SavedTeam, bool, error) {
	fake.getTeamMutex.Lock()
	fake.getTeamArgsForCall = append(fake.getTeamArgsForCall, struct{}{})
	fake.recordInvocation("GetTeam", []interface{}{})
	fake.getTeamMutex.Unlock()
	if fake.GetTeamStub != nil {
		return fake.GetTeamStub()
	} else {
		return fake.getTeamReturns.result1, fake.getTeamReturns.result2, fake.getTeamReturns.result3
	}
}

func (fake *FakeTeamDB) GetTeamCallCount() int {
	fake.getTeamMutex.RLock()
	defer fake.getTeamMutex.RUnlock()
	return len(fake.getTeamArgsForCall)
}

func (fake *FakeTeamDB) GetTeamReturns(result1 db.SavedTeam, result2 bool, result3 error) {
	fake.GetTeamStub = nil
	fake.getTeamReturns = struct {
		result1 db.SavedTeam
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamDB) UpdateBasicAuth(basicAuth *db.BasicAuth) (db.SavedTeam, error) {
	fake.updateBasicAuthMutex.Lock()
	fake.updateBasicAuthArgsForCall = append(fake.updateBasicAuthArgsForCall, struct {
		basicAuth *db.BasicAuth
	}{basicAuth})
	fake.recordInvocation("UpdateBasicAuth", []interface{}{basicAuth})
	fake.updateBasicAuthMutex.Unlock()
	if fake.UpdateBasicAuthStub != nil {
		return fake.UpdateBasicAuthStub(basicAuth)
	} else {
		return fake.updateBasicAuthReturns.result1, fake.updateBasicAuthReturns.result2
	}
}

func (fake *FakeTeamDB) UpdateBasicAuthCallCount() int {
	fake.updateBasicAuthMutex.RLock()
	defer fake.updateBasicAuthMutex.RUnlock()
	return len(fake.updateBasicAuthArgsForCall)
}

func (fake *FakeTeamDB) UpdateBasicAuthArgsForCall(i int) *db.BasicAuth {
	fake.updateBasicAuthMutex.RLock()
	defer fake.updateBasicAuthMutex.RUnlock()
	return fake.updateBasicAuthArgsForCall[i].basicAuth
}

func (fake *FakeTeamDB) UpdateBasicAuthReturns(result1 db.SavedTeam, result2 error) {
	fake.UpdateBasicAuthStub = nil
	fake.updateBasicAuthReturns = struct {
		result1 db.SavedTeam
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) UpdateGitHubAuth(gitHubAuth *db.GitHubAuth) (db.SavedTeam, error) {
	fake.updateGitHubAuthMutex.Lock()
	fake.updateGitHubAuthArgsForCall = append(fake.updateGitHubAuthArgsForCall, struct {
		gitHubAuth *db.GitHubAuth
	}{gitHubAuth})
	fake.recordInvocation("UpdateGitHubAuth", []interface{}{gitHubAuth})
	fake.updateGitHubAuthMutex.Unlock()
	if fake.UpdateGitHubAuthStub != nil {
		return fake.UpdateGitHubAuthStub(gitHubAuth)
	} else {
		return fake.updateGitHubAuthReturns.result1, fake.updateGitHubAuthReturns.result2
	}
}

func (fake *FakeTeamDB) UpdateGitHubAuthCallCount() int {
	fake.updateGitHubAuthMutex.RLock()
	defer fake.updateGitHubAuthMutex.RUnlock()
	return len(fake.updateGitHubAuthArgsForCall)
}

func (fake *FakeTeamDB) UpdateGitHubAuthArgsForCall(i int) *db.GitHubAuth {
	fake.updateGitHubAuthMutex.RLock()
	defer fake.updateGitHubAuthMutex.RUnlock()
	return fake.updateGitHubAuthArgsForCall[i].gitHubAuth
}

func (fake *FakeTeamDB) UpdateGitHubAuthReturns(result1 db.SavedTeam, result2 error) {
	fake.UpdateGitHubAuthStub = nil
	fake.updateGitHubAuthReturns = struct {
		result1 db.SavedTeam
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) UpdateUAAAuth(uaaAuth *db.UAAAuth) (db.SavedTeam, error) {
	fake.updateUAAAuthMutex.Lock()
	fake.updateUAAAuthArgsForCall = append(fake.updateUAAAuthArgsForCall, struct {
		uaaAuth *db.UAAAuth
	}{uaaAuth})
	fake.recordInvocation("UpdateUAAAuth", []interface{}{uaaAuth})
	fake.updateUAAAuthMutex.Unlock()
	if fake.UpdateUAAAuthStub != nil {
		return fake.UpdateUAAAuthStub(uaaAuth)
	} else {
		return fake.updateUAAAuthReturns.result1, fake.updateUAAAuthReturns.result2
	}
}

func (fake *FakeTeamDB) UpdateUAAAuthCallCount() int {
	fake.updateUAAAuthMutex.RLock()
	defer fake.updateUAAAuthMutex.RUnlock()
	return len(fake.updateUAAAuthArgsForCall)
}

func (fake *FakeTeamDB) UpdateUAAAuthArgsForCall(i int) *db.UAAAuth {
	fake.updateUAAAuthMutex.RLock()
	defer fake.updateUAAAuthMutex.RUnlock()
	return fake.updateUAAAuthArgsForCall[i].uaaAuth
}

func (fake *FakeTeamDB) UpdateUAAAuthReturns(result1 db.SavedTeam, result2 error) {
	fake.UpdateUAAAuthStub = nil
	fake.updateUAAAuthReturns = struct {
		result1 db.SavedTeam
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) UpdateGenericOAuth(genericOAuth *db.GenericOAuth) (db.SavedTeam, error) {
	fake.updateGenericOAuthMutex.Lock()
	fake.updateGenericOAuthArgsForCall = append(fake.updateGenericOAuthArgsForCall, struct {
		genericOAuth *db.GenericOAuth
	}{genericOAuth})
	fake.recordInvocation("UpdateGenericOAuth", []interface{}{genericOAuth})
	fake.updateGenericOAuthMutex.Unlock()
	if fake.UpdateGenericOAuthStub != nil {
		return fake.UpdateGenericOAuthStub(genericOAuth)
	} else {
		return fake.updateGenericOAuthReturns.result1, fake.updateGenericOAuthReturns.result2
	}
}

func (fake *FakeTeamDB) UpdateGenericOAuthCallCount() int {
	fake.updateGenericOAuthMutex.RLock()
	defer fake.updateGenericOAuthMutex.RUnlock()
	return len(fake.updateGenericOAuthArgsForCall)
}

func (fake *FakeTeamDB) UpdateGenericOAuthArgsForCall(i int) *db.GenericOAuth {
	fake.updateGenericOAuthMutex.RLock()
	defer fake.updateGenericOAuthMutex.RUnlock()
	return fake.updateGenericOAuthArgsForCall[i].genericOAuth
}

func (fake *FakeTeamDB) UpdateGenericOAuthReturns(result1 db.SavedTeam, result2 error) {
	fake.UpdateGenericOAuthStub = nil
	fake.updateGenericOAuthReturns = struct {
		result1 db.SavedTeam
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) GetConfig(pipelineName string) (atc.Config, atc.RawConfig, db.ConfigVersion, error) {
	fake.getConfigMutex.Lock()
	fake.getConfigArgsForCall = append(fake.getConfigArgsForCall, struct {
		pipelineName string
	}{pipelineName})
	fake.recordInvocation("GetConfig", []interface{}{pipelineName})
	fake.getConfigMutex.Unlock()
	if fake.GetConfigStub != nil {
		return fake.GetConfigStub(pipelineName)
	} else {
		return fake.getConfigReturns.result1, fake.getConfigReturns.result2, fake.getConfigReturns.result3, fake.getConfigReturns.result4
	}
}

func (fake *FakeTeamDB) GetConfigCallCount() int {
	fake.getConfigMutex.RLock()
	defer fake.getConfigMutex.RUnlock()
	return len(fake.getConfigArgsForCall)
}

func (fake *FakeTeamDB) GetConfigArgsForCall(i int) string {
	fake.getConfigMutex.RLock()
	defer fake.getConfigMutex.RUnlock()
	return fake.getConfigArgsForCall[i].pipelineName
}

func (fake *FakeTeamDB) GetConfigReturns(result1 atc.Config, result2 atc.RawConfig, result3 db.ConfigVersion, result4 error) {
	fake.GetConfigStub = nil
	fake.getConfigReturns = struct {
		result1 atc.Config
		result2 atc.RawConfig
		result3 db.ConfigVersion
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeTeamDB) SaveConfig(arg1 string, arg2 atc.Config, arg3 db.ConfigVersion, arg4 db.PipelinePausedState) (db.SavedPipeline, bool, error) {
	fake.saveConfigMutex.Lock()
	fake.saveConfigArgsForCall = append(fake.saveConfigArgsForCall, struct {
		arg1 string
		arg2 atc.Config
		arg3 db.ConfigVersion
		arg4 db.PipelinePausedState
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("SaveConfig", []interface{}{arg1, arg2, arg3, arg4})
	fake.saveConfigMutex.Unlock()
	if fake.SaveConfigStub != nil {
		return fake.SaveConfigStub(arg1, arg2, arg3, arg4)
	} else {
		return fake.saveConfigReturns.result1, fake.saveConfigReturns.result2, fake.saveConfigReturns.result3
	}
}

func (fake *FakeTeamDB) SaveConfigCallCount() int {
	fake.saveConfigMutex.RLock()
	defer fake.saveConfigMutex.RUnlock()
	return len(fake.saveConfigArgsForCall)
}

func (fake *FakeTeamDB) SaveConfigArgsForCall(i int) (string, atc.Config, db.ConfigVersion, db.PipelinePausedState) {
	fake.saveConfigMutex.RLock()
	defer fake.saveConfigMutex.RUnlock()
	return fake.saveConfigArgsForCall[i].arg1, fake.saveConfigArgsForCall[i].arg2, fake.saveConfigArgsForCall[i].arg3, fake.saveConfigArgsForCall[i].arg4
}

func (fake *FakeTeamDB) SaveConfigReturns(result1 db.SavedPipeline, result2 bool, result3 error) {
	fake.SaveConfigStub = nil
	fake.saveConfigReturns = struct {
		result1 db.SavedPipeline
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamDB) CreateOneOffBuild() (db.Build, error) {
	fake.createOneOffBuildMutex.Lock()
	fake.createOneOffBuildArgsForCall = append(fake.createOneOffBuildArgsForCall, struct{}{})
	fake.recordInvocation("CreateOneOffBuild", []interface{}{})
	fake.createOneOffBuildMutex.Unlock()
	if fake.CreateOneOffBuildStub != nil {
		return fake.CreateOneOffBuildStub()
	} else {
		return fake.createOneOffBuildReturns.result1, fake.createOneOffBuildReturns.result2
	}
}

func (fake *FakeTeamDB) CreateOneOffBuildCallCount() int {
	fake.createOneOffBuildMutex.RLock()
	defer fake.createOneOffBuildMutex.RUnlock()
	return len(fake.createOneOffBuildArgsForCall)
}

func (fake *FakeTeamDB) CreateOneOffBuildReturns(result1 db.Build, result2 error) {
	fake.CreateOneOffBuildStub = nil
	fake.createOneOffBuildReturns = struct {
		result1 db.Build
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) GetBuilds(page db.Page, publicOnly bool) ([]db.Build, db.Pagination, error) {
	fake.getBuildsMutex.Lock()
	fake.getBuildsArgsForCall = append(fake.getBuildsArgsForCall, struct {
		page       db.Page
		publicOnly bool
	}{page, publicOnly})
	fake.recordInvocation("GetBuilds", []interface{}{page, publicOnly})
	fake.getBuildsMutex.Unlock()
	if fake.GetBuildsStub != nil {
		return fake.GetBuildsStub(page, publicOnly)
	} else {
		return fake.getBuildsReturns.result1, fake.getBuildsReturns.result2, fake.getBuildsReturns.result3
	}
}

func (fake *FakeTeamDB) GetBuildsCallCount() int {
	fake.getBuildsMutex.RLock()
	defer fake.getBuildsMutex.RUnlock()
	return len(fake.getBuildsArgsForCall)
}

func (fake *FakeTeamDB) GetBuildsArgsForCall(i int) (db.Page, bool) {
	fake.getBuildsMutex.RLock()
	defer fake.getBuildsMutex.RUnlock()
	return fake.getBuildsArgsForCall[i].page, fake.getBuildsArgsForCall[i].publicOnly
}

func (fake *FakeTeamDB) GetBuildsReturns(result1 []db.Build, result2 db.Pagination, result3 error) {
	fake.GetBuildsStub = nil
	fake.getBuildsReturns = struct {
		result1 []db.Build
		result2 db.Pagination
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamDB) GetBuild(buildID int) (db.Build, bool, error) {
	fake.getBuildMutex.Lock()
	fake.getBuildArgsForCall = append(fake.getBuildArgsForCall, struct {
		buildID int
	}{buildID})
	fake.recordInvocation("GetBuild", []interface{}{buildID})
	fake.getBuildMutex.Unlock()
	if fake.GetBuildStub != nil {
		return fake.GetBuildStub(buildID)
	} else {
		return fake.getBuildReturns.result1, fake.getBuildReturns.result2, fake.getBuildReturns.result3
	}
}

func (fake *FakeTeamDB) GetBuildCallCount() int {
	fake.getBuildMutex.RLock()
	defer fake.getBuildMutex.RUnlock()
	return len(fake.getBuildArgsForCall)
}

func (fake *FakeTeamDB) GetBuildArgsForCall(i int) int {
	fake.getBuildMutex.RLock()
	defer fake.getBuildMutex.RUnlock()
	return fake.getBuildArgsForCall[i].buildID
}

func (fake *FakeTeamDB) GetBuildReturns(result1 db.Build, result2 bool, result3 error) {
	fake.GetBuildStub = nil
	fake.getBuildReturns = struct {
		result1 db.Build
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamDB) Workers() ([]db.SavedWorker, error) {
	fake.workersMutex.Lock()
	fake.workersArgsForCall = append(fake.workersArgsForCall, struct{}{})
	fake.recordInvocation("Workers", []interface{}{})
	fake.workersMutex.Unlock()
	if fake.WorkersStub != nil {
		return fake.WorkersStub()
	} else {
		return fake.workersReturns.result1, fake.workersReturns.result2
	}
}

func (fake *FakeTeamDB) WorkersCallCount() int {
	fake.workersMutex.RLock()
	defer fake.workersMutex.RUnlock()
	return len(fake.workersArgsForCall)
}

func (fake *FakeTeamDB) WorkersReturns(result1 []db.SavedWorker, result2 error) {
	fake.WorkersStub = nil
	fake.workersReturns = struct {
		result1 []db.SavedWorker
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) GetContainer(handle string) (db.SavedContainer, bool, error) {
	fake.getContainerMutex.Lock()
	fake.getContainerArgsForCall = append(fake.getContainerArgsForCall, struct {
		handle string
	}{handle})
	fake.recordInvocation("GetContainer", []interface{}{handle})
	fake.getContainerMutex.Unlock()
	if fake.GetContainerStub != nil {
		return fake.GetContainerStub(handle)
	} else {
		return fake.getContainerReturns.result1, fake.getContainerReturns.result2, fake.getContainerReturns.result3
	}
}

func (fake *FakeTeamDB) GetContainerCallCount() int {
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	return len(fake.getContainerArgsForCall)
}

func (fake *FakeTeamDB) GetContainerArgsForCall(i int) string {
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	return fake.getContainerArgsForCall[i].handle
}

func (fake *FakeTeamDB) GetContainerReturns(result1 db.SavedContainer, result2 bool, result3 error) {
	fake.GetContainerStub = nil
	fake.getContainerReturns = struct {
		result1 db.SavedContainer
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamDB) FindContainersByDescriptors(id db.Container) ([]db.SavedContainer, error) {
	fake.findContainersByDescriptorsMutex.Lock()
	fake.findContainersByDescriptorsArgsForCall = append(fake.findContainersByDescriptorsArgsForCall, struct {
		id db.Container
	}{id})
	fake.recordInvocation("FindContainersByDescriptors", []interface{}{id})
	fake.findContainersByDescriptorsMutex.Unlock()
	if fake.FindContainersByDescriptorsStub != nil {
		return fake.FindContainersByDescriptorsStub(id)
	} else {
		return fake.findContainersByDescriptorsReturns.result1, fake.findContainersByDescriptorsReturns.result2
	}
}

func (fake *FakeTeamDB) FindContainersByDescriptorsCallCount() int {
	fake.findContainersByDescriptorsMutex.RLock()
	defer fake.findContainersByDescriptorsMutex.RUnlock()
	return len(fake.findContainersByDescriptorsArgsForCall)
}

func (fake *FakeTeamDB) FindContainersByDescriptorsArgsForCall(i int) db.Container {
	fake.findContainersByDescriptorsMutex.RLock()
	defer fake.findContainersByDescriptorsMutex.RUnlock()
	return fake.findContainersByDescriptorsArgsForCall[i].id
}

func (fake *FakeTeamDB) FindContainersByDescriptorsReturns(result1 []db.SavedContainer, result2 error) {
	fake.FindContainersByDescriptorsStub = nil
	fake.findContainersByDescriptorsReturns = struct {
		result1 []db.SavedContainer
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) GetVolumes() ([]db.SavedVolume, error) {
	fake.getVolumesMutex.Lock()
	fake.getVolumesArgsForCall = append(fake.getVolumesArgsForCall, struct{}{})
	fake.recordInvocation("GetVolumes", []interface{}{})
	fake.getVolumesMutex.Unlock()
	if fake.GetVolumesStub != nil {
		return fake.GetVolumesStub()
	} else {
		return fake.getVolumesReturns.result1, fake.getVolumesReturns.result2
	}
}

func (fake *FakeTeamDB) GetVolumesCallCount() int {
	fake.getVolumesMutex.RLock()
	defer fake.getVolumesMutex.RUnlock()
	return len(fake.getVolumesArgsForCall)
}

func (fake *FakeTeamDB) GetVolumesReturns(result1 []db.SavedVolume, result2 error) {
	fake.GetVolumesStub = nil
	fake.getVolumesReturns = struct {
		result1 []db.SavedVolume
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamDB) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getPipelinesMutex.RLock()
	defer fake.getPipelinesMutex.RUnlock()
	fake.getPublicPipelinesMutex.RLock()
	defer fake.getPublicPipelinesMutex.RUnlock()
	fake.getPrivateAndAllPublicPipelinesMutex.RLock()
	defer fake.getPrivateAndAllPublicPipelinesMutex.RUnlock()
	fake.getPipelineByNameMutex.RLock()
	defer fake.getPipelineByNameMutex.RUnlock()
	fake.getPublicPipelineByNameMutex.RLock()
	defer fake.getPublicPipelineByNameMutex.RUnlock()
	fake.orderPipelinesMutex.RLock()
	defer fake.orderPipelinesMutex.RUnlock()
	fake.getTeamMutex.RLock()
	defer fake.getTeamMutex.RUnlock()
	fake.updateBasicAuthMutex.RLock()
	defer fake.updateBasicAuthMutex.RUnlock()
	fake.updateGitHubAuthMutex.RLock()
	defer fake.updateGitHubAuthMutex.RUnlock()
	fake.updateUAAAuthMutex.RLock()
	defer fake.updateUAAAuthMutex.RUnlock()
	fake.getConfigMutex.RLock()
	defer fake.getConfigMutex.RUnlock()
	fake.saveConfigMutex.RLock()
	defer fake.saveConfigMutex.RUnlock()
	fake.createOneOffBuildMutex.RLock()
	defer fake.createOneOffBuildMutex.RUnlock()
	fake.getBuildsMutex.RLock()
	defer fake.getBuildsMutex.RUnlock()
	fake.getBuildMutex.RLock()
	defer fake.getBuildMutex.RUnlock()
	fake.workersMutex.RLock()
	defer fake.workersMutex.RUnlock()
	fake.getContainerMutex.RLock()
	defer fake.getContainerMutex.RUnlock()
	fake.findContainersByDescriptorsMutex.RLock()
	defer fake.findContainersByDescriptorsMutex.RUnlock()
	fake.getVolumesMutex.RLock()
	defer fake.getVolumesMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeTeamDB) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ db.TeamDB = new(FakeTeamDB)
