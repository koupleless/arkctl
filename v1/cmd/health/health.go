/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package health

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koupleless/arkctl/common/cmdutil"
	"github.com/koupleless/arkctl/common/runtime"
	"github.com/koupleless/arkctl/common/style"
	"github.com/koupleless/arkctl/v1/cmd/root"
	"github.com/koupleless/arkctl/v1/service/ark"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	portFlag int    = 1238
	hostFlag string = "127.0.0.1"

	podFlag      string = ""
	podNamespace string = ""
	podName      string = ""
)

var HealthCommand = &cobra.Command{
	Use:   "health",
	Short: "get master service health status",
	Long: `
The arkctl health subcommand can help you quickly get master service health status, including JVM, CPU, Master Service Info and Biz Status.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if podFlag != "" && strings.Contains(podFlag, "/") {
			podNamespace, podName = strings.Split(podFlag, "/")[0], strings.Split(podFlag, "/")[1]
		} else {
			podNamespace, podName = "default", podFlag
		}

		return execHealth(context.Background())
	},
}

func execHealthLocal(ctx context.Context) error {
	arkService := ark.BuildService(ctx)
	healthStatus, err := arkService.Health(ctx, ark.HealthRequest{
		HostName: hostFlag,
		Port:     portFlag,
	})
	if err != nil {
		return err
	}
	style.InfoPrefix("QueryHealth").Println(healthStatus.Code)
	style.InfoPrefix("JVM").Println(string(runtime.MustReturnResult(json.MarshalIndent(healthStatus.Data.HealthData.Jvm, "", "    "))))
	style.InfoPrefix("CPU").Println(string(runtime.MustReturnResult(json.MarshalIndent(healthStatus.Data.HealthData.Cpu, "", "    "))))
	style.InfoPrefix("MasterBiz").Println(string(runtime.MustReturnResult(json.MarshalIndent(healthStatus.Data.HealthData.MasterBizInfo, "", "    "))))
	return nil
}

func execHealthKubePod(ctx context.Context) error {
	kubeQueryCmd := cmdutil.BuildCommand(
		ctx,
		"kubectl",
		"-n", podNamespace,
		"exec", podName, "--",
		"curl",
		"-X",
		"POST",
		fmt.Sprintf("http://127.0.0.1:%v/health", portFlag),
	)

	if err := kubeQueryCmd.Exec(); err != nil {
		pterm.Error.PrintOnError(err)
		return err
	}

	if stderroutput := <-kubeQueryCmd.Wait(); stderroutput != nil {
		stderrlines := stderroutput.Error()
		stdoutlines := &strings.Builder{}
		for line := range kubeQueryCmd.Output() {
			stdoutlines.WriteString(line)
		}

		if !strings.Contains(stdoutlines.String(), "SUCCESS") {
			pterm.Println(stderrlines)
			pterm.Println(stdoutlines)
			pterm.Error.Println("health status query failed")
			return fmt.Errorf("health status query failed")
		}
		style.InfoPrefix("QueryHealth").Println(stdoutlines)
	}
	return nil
}

func execHealth(ctx context.Context) error {
	switch {
	case podFlag != "":
		return execHealthKubePod(ctx)
	default:
		return execHealthLocal(ctx)
	}
}

func init() {
	root.RootCmd.AddCommand(HealthCommand)
	HealthCommand.Flags().IntVar(&portFlag, "port", portFlag, "ark container's port")
	HealthCommand.Flags().StringVar(&hostFlag, "host", hostFlag, "ark container's host")
	HealthCommand.Flags().StringVar(&podFlag, "pod", podFlag, "ark container's running pod")
}
