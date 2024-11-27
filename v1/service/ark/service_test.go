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

import (
	"context"
	"encoding/json"
	"github.com/koupleless/arkctl/common/fileutil"
	"github.com/koupleless/arkctl/common/osutil"
	"net"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func mockHttpServer(
	path string,
	handler func(w http.ResponseWriter, r *http.Request),
) (int, func()) {
	// Create a listener on a random port.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// Retrieve the port.
	port := listener.Addr().(*net.TCPAddr).Port
	mux.Handle(path, http.HandlerFunc(handler))

	server := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			logrus.Warn(err)
		}
	}()

	return port, func() {
		listener.Close()
	}
}

func TestInstallBiz_Success(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port, cancel := mockHttpServer("/installBiz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "SUCCESS",
			"message": "install biz success!",
		})
	})
	defer func() {
		cancel()
	}()

	err := client.InstallBiz(ctx, InstallBizRequest{
		BizModel: BizModel{
			BizName:    "biz",
			BizVersion: "0.0.1-SNAPSHOT",
			BizUrl:     "",
		},
		TargetContainer: ArkContainerRuntimeInfo{
			RunType: ArkContainerRunTypeLocal,
			Port:    &port,
		},
	})
	assert.Nil(t, err)

}

func TestInstallBiz_Failed(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port, cancel := mockHttpServer("/installBiz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":            "FAILED",
			"message":         "install biz failed!",
			"errorStackTrace": "this is the error stack trace!",
		})
	})
	defer func() {
		cancel()
	}()

	err := client.InstallBiz(ctx, InstallBizRequest{
		BizModel: BizModel{
			BizName:    "biz",
			BizVersion: "0.0.1-SNAPSHOT",
			BizUrl:     "",
		},
		TargetContainer: ArkContainerRuntimeInfo{
			RunType: ArkContainerRunTypeLocal,
			Port:    &port,
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, "install biz failed: install biz failed! \n Caused by: this is the error stack trace!", err.Error())
}

func TestInstallBiz_NoServer(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port := 8888

	err := client.InstallBiz(ctx, InstallBizRequest{
		BizModel: BizModel{
			BizName:    "biz",
			BizVersion: "0.0.1-SNAPSHOT",
			BizUrl:     fileutil.FileUrl(osutil.GetLocalFileProtocol() + "/foobar"),
		},
		TargetContainer: ArkContainerRuntimeInfo{
			RunType: ArkContainerRunTypeLocal,
			Port:    &port,
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, "Post \"http://127.0.0.1:8888/installBiz\": dial tcp 127.0.0.1:8888: connect: connection refused", err.Error())
}

func TestUnInstallBiz_Success(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port, cancel := mockHttpServer("/uninstallBiz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "SUCCESS",
			"message": "uninstall biz success!",
		})
	})
	defer func() {
		cancel()
	}()

	err := client.UnInstallBiz(ctx, UnInstallBizRequest{
		BizModel: BizModel{
			BizName:    "biz",
			BizVersion: "0.0.1-SNAPSHOT",
			BizUrl:     "",
		},
		TargetContainer: ArkContainerRuntimeInfo{
			RunType: ArkContainerRunTypeLocal,
			Port:    &port,
		},
	})
	assert.Nil(t, err)

}

func TestUnInstallBiz_NotInstalled(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port, cancel := mockHttpServer("/uninstallBiz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "FAILED",
			"message": "uninstall biz failed!",
			"data": map[string]interface{}{
				"code": "NOT_FOUND_BIZ",
			},
		})
	})
	defer func() {
		cancel()
	}()

	err := client.UnInstallBiz(ctx, UnInstallBizRequest{
		BizModel: BizModel{
			BizName:    "biz",
			BizVersion: "0.0.1-SNAPSHOT",
			BizUrl:     "",
		},
		TargetContainer: ArkContainerRuntimeInfo{
			RunType: ArkContainerRunTypeLocal,
			Port:    &port,
		},
	})
	assert.Nil(t, err)

}

func TestUnInstallBiz_Failed(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port, cancel := mockHttpServer("/uninstallBiz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "FAILED",
			"message": "uninstall biz failed!",
			"data": map[string]interface{}{
				"code": "FOO",
			},
		})
	})
	defer func() {
		cancel()
	}()

	err := client.UnInstallBiz(ctx, UnInstallBizRequest{
		BizModel: BizModel{
			BizName:    "biz",
			BizVersion: "0.0.1-SNAPSHOT",
			BizUrl:     "",
		},
		TargetContainer: ArkContainerRuntimeInfo{
			RunType: ArkContainerRunTypeLocal,
			Port:    &port,
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, "uninstall biz failed: {{FAILED {FOO  0 []} uninstall biz failed! }}", err.Error())

}

func TestUnInstallBiz_NoServer(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port := 8888

	err := client.UnInstallBiz(ctx, UnInstallBizRequest{
		BizModel: BizModel{
			BizName:    "biz",
			BizVersion: "0.0.1-SNAPSHOT",
			BizUrl:     "",
		},
		TargetContainer: ArkContainerRuntimeInfo{
			RunType: ArkContainerRunTypeLocal,
			Port:    &port,
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, "Post \"http://127.0.0.1:8888/uninstallBiz\": dial tcp 127.0.0.1:8888: connect: connection refused", err.Error())

}

func TestQueryAllBiz_HappyPath(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port, cancel := mockHttpServer("/queryAllBiz", func(w http.ResponseWriter, r *http.Request) {
		// {"code":"SUCCESS","data":[{"bizName":"biz1","bizState":"ACTIVATED","bizVersion":"0.0.1-SNAPSHOT","mainClass":"com.alipay.sofa.web.biz1.Biz1Application","webContextPath":"biz1"}]}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code": "SUCCESS",
			"data": []map[string]interface{}{
				{
					"bizName":        "biz1",
					"bizState":       "ACTIVATED",
					"bizVersion":     "0.0.1-SNAPSHOT",
					"mainClass":      "com.alipay.sofa.web.biz1.Biz1Application",
					"webContextPath": "biz1",
					"bizStateRecords": []map[string]interface{}{
						{
							"changeTime": 12345,
							"state":      "ACTIVATED",
						},
					},
				},
			},
		})
	})

	defer func() {
		cancel()
	}()

	info, err := client.QueryAllBiz(ctx, QueryAllArkBizRequest{
		HostName: "127.0.0.1",
		Port:     port,
	})

	assert.Nil(t, err)
	assert.Equal(t, &QueryAllArkBizResponse{
		GenericArkResponseBase: GenericArkResponseBase[[]ArkBizInfo]{
			Code: "SUCCESS",
			Data: []ArkBizInfo{
				{
					BizName:        "biz1",
					BizState:       "ACTIVATED",
					BizVersion:     "0.0.1-SNAPSHOT",
					MainClass:      "com.alipay.sofa.web.biz1.Biz1Application",
					WebContextPath: "biz1",
					BizStateRecords: []ArkBizStateRecord{
						{
							ChangeTime: 12345,
							State:      "ACTIVATED",
						},
					},
				},
			},
		},
	}, info)
}

func TestIsSuccessResponse(t *testing.T) {
	assert.Nil(t, IsSuccessResponse(&GenericArkResponseBase[int]{
		Code: "SUCCESS",
	}))

	assert.Errorf(t, IsSuccessResponse(&GenericArkResponseBase[int]{
		Code:    "FAILED",
		Message: "failed",
	}), "sofa-ark failed response: %s", "failed")
}

func TestUnImplemented(t *testing.T) {

}

func TestHealth_HappyPath(t *testing.T) {
	ctx := context.Background()
	client := BuildService(ctx)
	port, cancel := mockHttpServer("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code": "SUCCESS",
			"data": map[string]interface{}{
				"healthData": map[string]interface{}{
					"cpu": map[string]interface{}{
						"count":           12,
						"free (%)":        82.90837318159039,
						"system used (%)": 6.207871242418149,
						"total used (%)":  139957397,
						"type":            "Apple M3 Pro",
						"user used (%)":   10.883755575991456,
					},
					"jvm": map[string]interface{}{
						"committed heap memory(M)":     286.5,
						"committed non heap memory(M)": 140.3359375,
						"free memory(M)":               78.14934539794922,
						"init heap memory(M)":          288,
						"init non heap memory(M)":      2.4375,
						"java home":                    "/Library/Java/JavaVirtualMachines/jdk1.8.0_291.jdk/Contents/Home/jre",
						"java version":                 "1.8.0_291",
						"loaded class count":           13155,
						"max heap memory(M)":           4096,
						"max memory(M)":                4096,
						"max non heap memory(M)":       -9.5367431640625e-7,
						"run time(s)":                  497862.907,
						"total class count":            13224,
						"total memory(M)":              286.5,
						"unload class count":           69,
						"used heap memory(M)":          208.35065460205078,
						"used non heap memory(M)":      129.16268920898438,
					},
					"masterBizInfo": map[string]interface{}{
						"bizName":        "base",
						"bizState":       "ACTIVATED",
						"bizVersion":     "1.0.0",
						"webContextPath": "/",
					},
				},
			},
		})
	})

	defer func() {
		cancel()
	}()

	info, err := client.Health(ctx, HealthRequest{
		HostName: "127.0.0.1",
		Port:     port,
	})

	assert.Nil(t, err)
	assert.Equal(t, &HealthResponse{
		GenericArkResponseBase: GenericArkResponseBase[HealthInfo]{
			Code: "SUCCESS",
			Data: HealthInfo{
				HealthData{
					Cpu: CpuInfo{
						Count:      12,
						Free:       82.90837318159039,
						SystemUsed: 6.207871242418149,
						TotalUsed:  139957397,
						Type:       "Apple M3 Pro",
						UserUsed:   10.883755575991456,
					},
					Jvm: JvmInfo{
						CommittedHeapMemoryM:    286.5,
						CommittedNonHeapMemoryM: 140.3359375,
						FreeMemoryM:             78.14934539794922,
						InitHeapMemoryM:         288,
						InitNonHeapMemoryM:      2.4375,
						JavaHome:                "/Library/Java/JavaVirtualMachines/jdk1.8.0_291.jdk/Contents/Home/jre",
						JavaVersion:             "1.8.0_291",
						LoadedClassCount:        13155,
						MaxHeapMemoryM:          4096,
						MaxMemoryM:              4096,
						MaxNonHeapMemoryM:       -9.5367431640625e-7,
						RunTimeS:                497862.907,
						TotalClassCount:         13224,
						TotalMemoryM:            286.5,
						UnloadClassCount:        69,
						UsedHeapMemoryM:         208.35065460205078,
						UsedNonHeapMemoryM:      129.16268920898438,
					},
					MasterBizInfo: MasterBizInfo{
						BizName:        "base",
						BizState:       "ACTIVATED",
						BizVersion:     "1.0.0",
						WebContextPath: "/",
					},
				},
			},
		},
	}, info)
}
