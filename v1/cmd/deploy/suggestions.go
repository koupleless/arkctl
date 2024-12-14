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

package deploy

import (
	"strings"

	"github.com/koupleless/arkctl/common/style"
)

const (
	faq_url = "https://koupleless.io/en/docs/faq/faq/"
)

var suggestionFuncs = []func(errorOutputLines []string) bool{
	suggestionBaseNotStart,
	suggestionMavenExecutableNotFound,
	suggestionMavenVersionTooLow,
	suggestWebContextPathConflict,
	suggestApplicationProperties,
	suggestImportSpringBootAutoConfiguration,
	suggestJvmInitializingFailed,
}

func printSuggestion(err error) {
	printSuggestionWithMore(err, nil)
}

func printSuggestionWithMore(err error, subprocessOutput []string) {
	var errorOutputLines []string
	if err != nil {
		if lines := strings.Split(err.Error(), "\n"); len(lines) > 0 {
			errorOutputLines = append(errorOutputLines, lines...)
		}
	}
	if len(subprocessOutput) > 0 {
		errorOutputLines = append(errorOutputLines, subprocessOutput...)
	}

	for _, suggestionFunc := range suggestionFuncs {
		if suggestionFunc(errorOutputLines) {
			break
		}
	}

	doPrintSuggestion("you can go to faq for more help at " + faq_url)
}

func suggestionBaseNotStart(errorOutputLines []string) bool {
	hasBaseNotStart := false
	for _, line := range errorOutputLines {
		if strings.HasSuffix(line, "connect: connection refused") {
			if strings.Contains(line, "installBiz") {
				hasBaseNotStart = true
				break
			}
		}
	}
	if hasBaseNotStart {
		doPrintSuggestion("ensure target base is running")
		return true
	}
	return false
}

func suggestionMavenExecutableNotFound(errorOutputLines []string) bool {
	hasMavenExecutableNotFound := false
	for _, line := range errorOutputLines {
		if strings.Contains(line, "exec: \"mvn\": executable file not found") {
			hasMavenExecutableNotFound = true
			break
		}
	}
	if hasMavenExecutableNotFound {
		doPrintSuggestion("install latest maven or just put mvn executable path into your $PATH")
		return true
	}
	return false
}

func suggestionMavenVersionTooLow(errorOutputLines []string) bool {
	hasMavenVersionTooLow := false
	var featureSubStrings = []string{
		"Unable to parse configuration of mojo com.alipay.sofa:sofa-ark-maven-plugin",
		"com.google.inject.ProvisionException: Unable to provision",
		"Error injecting: private org.eclipse.aether.spi.log.Logger",
		"Can not set org.eclipse.aether.spi.log.Logger field",
	}
	for _, line := range errorOutputLines {
		for _, featureSubString := range featureSubStrings {
			if strings.Contains(line, featureSubString) {
				hasMavenVersionTooLow = true
				break
			}
		}
	}
	if hasMavenVersionTooLow {
		doPrintSuggestion("your maven is outdated, update it to 3.6.1 or higher version")
		return true
	}
	return false
}

func suggestWebContextPathConflict(errorOutputLines []string) bool {
	hasStartWebServer := false
	hasChildNameNotUnique := false
	for _, line := range errorOutputLines {
		if strings.Contains(line, "Unable to start web server") {
			hasStartWebServer = true
		}
		if hasStartWebServer {
			if strings.Contains(line, "Child name") && strings.Contains(line, "is not unique") {
				hasChildNameNotUnique = true
				break
			}
		}
	}
	if hasChildNameNotUnique {
		doPrintSuggestion("another installed biz module has the same webContextPath as yours")
		doPrintSuggestion("change your <webContextPath> in pom.xml or uninstall another biz module")
		return true
	}
	return false
}

func suggestApplicationProperties(errorOutputLines []string) bool {
	for _, line := range errorOutputLines {
		if strings.Contains(line, "spring.application.name must be configured") {
			doPrintSuggestion("add \"spring.application.name\" config into your application.properties")
			return true
		}
	}
	return false
}

func suggestImportSpringBootAutoConfiguration(errorOutputLines []string) bool {
	for _, line := range errorOutputLines {
		if strings.Contains(line, "The following classes could not be excluded because they are not auto-configuration classes") &&
			strings.Contains(line, "org.springframework.boot.actuate.autoconfigure.startup.StartupEndpointAutoConfiguration") {
			doPrintSuggestion("import sprign-boot-actuator-autoconfiguration artifact in your pom.xml file")
			return true
		}
	}
	return false
}

func suggestJvmInitializingFailed(errorOutputLines []string) bool {
	for _, line := range errorOutputLines {
		if strings.Contains(line, "Error occurred during initialization of VM") {
			doPrintSuggestion("check your jvm staring paramaters")
			return true
		}
	}
	return false
}

func doPrintSuggestion(line string) {
	style.InfoPrefix("Suggestion").Printfln("%s", line)
}
