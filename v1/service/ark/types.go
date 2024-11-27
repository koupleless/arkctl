/**
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ark

import "github.com/koupleless/arkctl/common/fileutil"

type ArkContainerRunType string

const (
	ArkContainerRunTypeLocal ArkContainerRunType = "local"
	ArkContainerRunTypeVM    ArkContainerRunType = "vm" // the reason why we need vm is we might use scp to copy file to vm server
	ArkContainerRunTypeK8s   ArkContainerRunType = "pod"
)

// ArkResponseData is the response data of ark api.
type ArkResponseData struct {
	Code         string        `json:"code"`
	Message      string        `json:"message"`
	ElapsedSpace int           `json:"elapsedSpace"`
	BizInfos     []interface{} `json:"bizInfos"`
}

// GenericArkResponseBase is the base response of ark api.
type GenericArkResponseBase[T any] struct {
	// Code is the response code
	Code string `json:"code"`

	// Data is the response data
	Data T `json:"data"`

	// Message is the error message
	Message string `json:"message"`
}

// ArkResponseBase is the base response of ark api.
type ArkResponseBase struct {
	// Code is the response code
	Code string `json:"code"`

	// Data is the response data
	Data ArkResponseData `json:"data"`

	// Message is the error message
	Message string `json:"message"`

	// ErrorStackTrace is the error stack trace
	ErrorStackTrace string `json:"errorStackTrace"`
}

// ArkContainerRuntimeInfo contains necessary info of an ark container.
type ArkContainerRuntimeInfo struct {
	// RunType is the type of ark container, like local, vm server, pod, etc.
	RunType ArkContainerRunType `json:"runType"`

	// Coordinate is the exact location of ark container.
	// If the RunType is local, then it's localhost.
	// If the RunType is vm server, then it's ip.
	// If the RunType is pod, then it's the {namespace}/{podName}
	Coordinate string `json:"coordinate"`

	// Port is the ark api port of ark container.
	Port *int `json:"port"`
}

func (info *ArkContainerRuntimeInfo) GetPort() int {
	if info.Port == nil {
		// default port
		return 1238
	}
	return *info.Port
}

// BizModel contains necessary metadata info of an ark biz.
// Usually this BizModel is generated automatically from bundle like jarfile.
type BizModel struct {
	// BizName the biz name
	BizName string `json:"bizName,omitempty"`

	// BizVersion is the version of biz module.
	BizVersion string `json:"bizVersion,omitempty"`

	// BizUrl is the location of source code.
	BizUrl fileutil.FileUrl `json:"bizUrl,omitempty"`
}

// InstallBizRequest is the request for installing biz module to ark container.
type InstallBizRequest struct {
	// BizModel is the metadata a given biz module.
	BizModel BizModel `json:"bizModel"`

	// TargetContainer is the target ark container we want to install a biz module to.
	TargetContainer ArkContainerRuntimeInfo `json:"targetContainer"`

	// InstallType is the type of install.
	// If the InstallType is "filesystem", then it will install the local biz module.
	// If the InstallType is "http", then it will install the remote biz module.
	InstallType string `json:"installType"`

	// BizHomeDir is the location of all biz module.
	// If not given, we will use {tmp}/arkBiz/ dir instead.
	// This will only be used when install biz module from local filesystem.
	BizHomeDir *string `json:"bizHomeDir"`
}

// InstallBizResponse is the response for installing biz module to ark container.
type InstallBizResponse struct {
	ArkResponseBase
}

// UnInstallBizRequest is the request for installing biz module to ark container.
type UnInstallBizRequest struct {
	// BizModel is the metadata a given biz module.
	BizModel BizModel `json:"bizModel"`

	// TargetContainer is the target ark container we want to install a biz module to.
	TargetContainer ArkContainerRuntimeInfo `json:"targetContainer"`
}

// UnInstallBizResponse is the response for installing biz module to ark container.
type UnInstallBizResponse struct {
	ArkResponseBase
}

// QueryAllArkBizRequest is the request for querying all biz module in a given ark container.
type QueryAllArkBizRequest struct {
	// HostName is where the ark container is running
	HostName string

	// Port is where the ark container is serving
	Port int
}

// HealthRequest is the request for health status of base runtime.
type HealthRequest struct {
	// HostName is where the ark container is running
	HostName string

	// Port is where the ark container is serving
	Port int
}

// ArkBizStateRecord is the response for biz module state record
type ArkBizStateRecord struct {
	ChangeTime int64  `json:"changeTime"`
	State      string `json:"state"`
	Reason     string `json:"reason"`
	Message    string `json:"message"`
}

// ArkBizInfo is the response for querying all biz module in a given ark container.
type ArkBizInfo struct {
	BizName         string              `json:"bizName"`
	BizState        string              `json:"bizState"`
	BizVersion      string              `json:"bizVersion"`
	MainClass       string              `json:"mainClass"`
	WebContextPath  string              `json:"webContextPath"`
	BizStateRecords []ArkBizStateRecord `json:"bizStateRecords"`
}

// QueryAllArkBizResponse is the response for querying all biz module in a given ark container.
type QueryAllArkBizResponse struct {
	GenericArkResponseBase[[]ArkBizInfo]
}

type JvmInfo struct {
	MaxNonHeapMemoryM       float64 `json:"max non heap memory(M)"`
	JavaVersion             string  `json:"java version"`
	MaxMemoryM              float64 `json:"max memory(M)"`
	MaxHeapMemoryM          float64 `json:"max heap memory(M)"`
	UsedHeapMemoryM         float64 `json:"used heap memory(M)"`
	UsedNonHeapMemoryM      float64 `json:"used non heap memory(M)"`
	LoadedClassCount        int     `json:"loaded class count"`
	InitNonHeapMemoryM      float64 `json:"init non heap memory(M)"`
	TotalMemoryM            float64 `json:"total memory(M)"`
	FreeMemoryM             float64 `json:"free memory(M)"`
	UnloadClassCount        int     `json:"unload class count"`
	TotalClassCount         int     `json:"total class count"`
	CommittedHeapMemoryM    float64 `json:"committed heap memory(M)"`
	JavaHome                string  `json:"java home"`
	InitHeapMemoryM         float64 `json:"init heap memory(M)"`
	CommittedNonHeapMemoryM float64 `json:"committed non heap memory(M)"`
	RunTimeS                float64 `json:"run time(s)"`
	JavaUsedMetaspace       int64   `json:"java used metaspace"`
	JavaCommittedMetaspace  int64   `json:"java committed metaspace"`
	JavaMaxMetaspace        int64   `json:"java max metaspace"`
}

type CpuInfo struct {
	Count      int     `json:"count"`
	TotalUsed  float64 `json:"total used (%)"`
	Type       string  `json:"type"`
	UserUsed   float64 `json:"user used (%)"`
	Free       float64 `json:"free (%)"`
	SystemUsed float64 `json:"system used (%)"`
}

type MasterBizInfo struct {
	BizName        string `json:"bizName"`
	BizState       string `json:"bizState"`
	BizVersion     string `json:"bizVersion"`
	WebContextPath string `json:"webContextPath"`
}

// HealthData is the response data of querying health status of base runtime.
type HealthData struct {
	Jvm           JvmInfo       `json:"jvm"`
	Cpu           CpuInfo       `json:"cpu"`
	MasterBizInfo MasterBizInfo `json:"masterBizInfo"`
}

// HealthInfo is the response for querying health status of base runtime.
type HealthInfo struct {
	HealthData HealthData `json:"healthData"`
}

// HealthResponse is the response of health check of base runtime.
type HealthResponse struct {
	GenericArkResponseBase[HealthInfo]
}
