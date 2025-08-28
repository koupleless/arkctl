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
package deploy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseMavenCommand(t *testing.T) {
	t.Parallel()

	defaultArgs := append([]string(nil), defaultMavenArgs...)

	tests := []struct {
		name     string
		mvnCmd   string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "empty command",
			mvnCmd:   "",
			wantCmd:  defaultMavenCmd,
			wantArgs: defaultArgs,
		},
		{
			name:     "custom command without args",
			mvnCmd:   "mvnw",
			wantCmd:  "mvnw",
			wantArgs: []string{},
		},
		{
			name:     "custom command with args",
			mvnCmd:   "mvn clean install",
			wantCmd:  "mvn",
			wantArgs: []string{"clean", "install"},
		},
		{
			name:     "complex command with args",
			mvnCmd:   "mvnw clean package -DskipTests",
			wantCmd:  "mvnw",
			wantArgs: []string{"clean", "package", "-DskipTests"},
		},
		{
			name:     "command with multiple spaces",
			mvnCmd:   "mvn    clean    package",
			wantCmd:  "mvn",
			wantArgs: []string{"clean", "package"},
		},
		{
			name:     "whitespace-only command",
			mvnCmd:   " \t\n ",
			wantCmd:  defaultMavenCmd,
			wantArgs: defaultArgs,
		},
		{
			name:     "leading and trailing spaces",
			mvnCmd:   "   mvn clean package   ",
			wantCmd:  "mvn",
			wantArgs: []string{"clean", "package"},
		},
		{
			name:     "tabs between tokens",
			mvnCmd:   "mvn\tclean\tinstall",
			wantCmd:  "mvn",
			wantArgs: []string{"clean", "install"},
		},
		{
			name:     "windows wrapper with flags",
			mvnCmd:   "mvnw.cmd -q -T1C",
			wantCmd:  "mvnw.cmd",
			wantArgs: []string{"-q", "-T1C"},
		},
		{
			name:     "absolute path mvn",
			mvnCmd:   "/usr/local/bin/mvn -q",
			wantCmd:  "/usr/local/bin/mvn",
			wantArgs: []string{"-q"},
		},
		{
			name:     "relative path wrapper",
			mvnCmd:   "./mvnw -q",
			wantCmd:  "./mvnw",
			wantArgs: []string{"-q"},
		},
		{
			name:     "args with equals and profile",
			mvnCmd:   "mvn -DskipTests=true -Pprod",
			wantCmd:  "mvn",
			wantArgs: []string{"-DskipTests=true", "-Pprod"},
		},
		{
			name:     "thread-count with separate value",
			mvnCmd:   "mvn -T 1C",
			wantCmd:  "mvn",
			wantArgs: []string{"-T", "1C"},
		},
		{
			name:     "multiple profiles in one flag",
			mvnCmd:   "mvn -Pprod,dev",
			wantCmd:  "mvn",
			wantArgs: []string{"-Pprod,dev"},
		},
		{
			name:     "reactor and project list",
			mvnCmd:   "mvn -am -pl :module",
			wantCmd:  "mvn",
			wantArgs: []string{"-am", "-pl", ":module"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotCmd, gotArgs := parseMavenCommand(tt.mvnCmd)
			if gotCmd != tt.wantCmd {
				t.Errorf("parseMavenCommand() gotCmd = %q, want %q", gotCmd, tt.wantCmd)
			}
			if diff := cmp.Diff(tt.wantArgs, gotArgs, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("parseMavenCommand() args mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
