// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package executors

import (
	pb "beam.apache.org/playground/backend/internal/api/v1"
	"beam.apache.org/playground/backend/internal/environment"
	"beam.apache.org/playground/backend/internal/validators"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

var (
	executorConfig = environment.NewExecutorConfig("javac", "java", []string{"-d", "bin", "-classpath"}, []string{"-cp", "bin:"})
	env            = environment.NewEnvironment(environment.NetworkEnvs{}, *environment.NewBeamEnvs(pb.Sdk_SDK_JAVA, executorConfig), environment.ApplicationEnvs{})
)

// BaseExecutorBuilder fills up an executor with base parameters
func BaseExecutorBuilder(envs environment.BeamEnvs, workingDir string, filePath string, validatorsFuncs *[]validators.Validator) *ExecutorBuilder {
	if validatorsFuncs == nil {
		v := make([]validators.Validator, 0)
		validatorsFuncs = &v
	}
	builder := NewExecutorBuilder().
		WithCompiler().
		WithCommand(envs.ExecutorConfig.CompileCmd).
		WithArgs(envs.ExecutorConfig.CompileArgs).
		WithFileName(filePath).
		WithWorkingDir(workingDir).
		WithRunner().
		WithCommand(envs.ExecutorConfig.RunCmd).
		WithArgs(envs.ExecutorConfig.RunArgs).
		WithClassName("HelloWorld").
		WithWorkingDir(workingDir).
		WithValidator().
		WithSdkValidators(validatorsFuncs).ExecutorBuilder
	return &builder
}

func TestExecutor_Compile(t *testing.T) {
	type fields struct {
		compileArgs CmdConfiguration
		runArgs     CmdConfiguration
		validators  []validators.Validator
	}
	tests := []struct {
		name   string
		fields fields
		want   *exec.Cmd
	}{
		{
			name: "TestCompile",
			fields: fields{
				compileArgs: CmdConfiguration{
					fileName:    "filePath",
					workingDir:  "./",
					commandName: "javac",
					commandArgs: []string{"-d", "bin", "-classpath", "/opt/apache/beam/jars/beam-sdks-java-harness.jar"},
				},
			},
			want: &exec.Cmd{
				Path:         "/usr/bin/javac",
				Args:         []string{"javac", "-d", "bin", "-classpath", "/opt/apache/beam/jars/beam-sdks-java-harness.jar", "filePath"},
				Env:          nil,
				Dir:          "",
				Stdin:        nil,
				Stdout:       nil,
				Stderr:       nil,
				ExtraFiles:   nil,
				SysProcAttr:  nil,
				Process:      nil,
				ProcessState: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Executor{
				compileArgs: tt.fields.compileArgs,
				runArgs:     tt.fields.runArgs,
				validators:  tt.fields.validators,
			}
			if got := ex.Compile(); !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("WithCompiler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutor_Run(t *testing.T) {
	type fields struct {
		compileArgs CmdConfiguration
		runArgs     CmdConfiguration
		validators  []validators.Validator
	}
	tests := []struct {
		name   string
		fields fields
		want   *exec.Cmd
	}{
		{
			name: "TestRun",
			fields: fields{
				runArgs: CmdConfiguration{
					fileName:    "HelloWorld",
					workingDir:  "./",
					commandName: "java",
					commandArgs: []string{"-cp", "bin:/opt/apache/beam/jars/beam-sdks-java-harness.jar:" +
						"/opt/apache/beam/jars/beam-runners-direct.jar:/opt/apache/beam/jars/slf4j-jdk14.jar"},
				},
			},
			want: &exec.Cmd{
				Path: "/usr/bin/java",
				Args: []string{"java", "-cp", "bin:/opt/apache/beam/jars/beam-sdks-java-harness.jar:" +
					"/opt/apache/beam/jars/beam-runners-direct.jar:/opt/apache/beam/jars/slf4j-jdk14.jar", "HelloWorld"},
				Env:          nil,
				Dir:          "",
				Stdin:        nil,
				Stdout:       nil,
				Stderr:       nil,
				ExtraFiles:   nil,
				SysProcAttr:  nil,
				Process:      nil,
				ProcessState: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Executor{
				compileArgs: tt.fields.compileArgs,
				runArgs:     tt.fields.runArgs,
				validators:  tt.fields.validators,
			}
			if got := ex.Run(); !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("WithRunner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseExecutorBuilder(t *testing.T) {
	validatorsFuncs := validators.GetJavaValidators("filePath")

	type args struct {
		envs            environment.BeamEnvs
		workingDir      string
		filePath        string
		validatorsFuncs *[]validators.Validator
	}
	tests := []struct {
		name string
		args args
		want Executor
	}{
		{
			name: "NewCmdProvider",
			args: args{
				envs:            env.BeamSdkEnvs,
				workingDir:      "./",
				filePath:        "filePath",
				validatorsFuncs: validatorsFuncs,
			},
			want: Executor{
				compileArgs: CmdConfiguration{
					fileName:    "filePath",
					workingDir:  "./",
					commandName: "javac",
					commandArgs: []string{"-d", "bin", "-classpath"},
				},
				runArgs: CmdConfiguration{
					fileName:    "HelloWorld",
					workingDir:  "./",
					commandName: "java",
					commandArgs: []string{"-cp", "bin:"},
				},
				validators: *validatorsFuncs,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseExecutorBuilder(tt.args.envs, tt.args.workingDir, tt.args.filePath, tt.args.validatorsFuncs).Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseExecutorBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
	err := os.RemoveAll("configs")
	if err != nil {
		return
	}
}