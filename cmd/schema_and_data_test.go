// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
        "flag"
        "testing"
        "github.com/GoogleCloudPlatform/spanner-migration-tool/common/constants"
        "github.com/stretchr/testify/assert"
)

func TestSchemaAndDataSetFlags(t *testing.T) {
        testCases:=[]struct {
                testName       string
                flagArgs      []string
                expectedValues SchemaAndDataCmd
        }{
                {
                        testName: "Default Values",
                        flagArgs: []string{},
                        expectedValues: SchemaAndDataCmd{
                                source:           "",
                                sourceProfile:    "",
                                target:           "Spanner",
                                targetProfile:    "",
                                filePrefix:       "",
                                WriteLimit:       DefaultWritersLimit,
                                dryRun:           false,
                                logLevel:         "DEBUG",
                                SkipForeignKeys:  false,
                                validate:         false,
                                dataflowTemplate: constants.DEFAULT_TEMPLATE_PATH,
                        },
                },
        }

        for _, tc:= range testCases {
                t.Run(tc.testName, func(t *testing.T) {
                        fs:= flag.NewFlagSet("testSetFlags", flag.ContinueOnError)
                        fs.Parse(tc.flagArgs)
                        schemaAndDataCmd:= SchemaAndDataCmd{}
                        schemaAndDataCmd.SetFlags(fs)
                        assert.Equal(t, tc.expectedValues, schemaAndDataCmd, tc.testName)
                })
        }
}
