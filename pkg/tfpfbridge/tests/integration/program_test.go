// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfbridgeintegrationtests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func TestBasicProgram(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	bin := filepath.Join(wd, "..", "bin")

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Env:         []string{fmt.Sprintf("PATH=%s", bin)},
		Dir:         filepath.Join("..", "testdata", "basicprogram"),
		SkipRefresh: true,

		PrepareProject: func(*engine.Projinfo) error {
			return ensureTestBridgeProviderCompiled(wd)
		},

		ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
			requiredInputStringCopy, ok := stack.Outputs["requiredInputStringCopy"]
			assert.True(t, ok)
			assert.Equal(t, "input1", requiredInputStringCopy)
		},
	})
}

func ensureTestBridgeProviderCompiled(wd string) error {
	exe := "pulumi-resource-testbridge"
	cmd := exec.Command("go", "build", "-o", filepath.Join("..", "..", "..", "bin", exe))
	cmd.Dir = filepath.Join(wd, "..", "internal", "cmd", exe)
	return cmd.Run()
}
