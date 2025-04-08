// Copyright 2023 OpenSSF Scorecard Authors
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

//nolint:stylecheck
package hasDangerousWorkflowUntrustedCheckout

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/ossf/scorecard/v5/checker"
	"github.com/ossf/scorecard/v5/finding"
	"github.com/ossf/scorecard/v5/probes/internal/utils/test"
)

func Test_Run(t *testing.T) {
	t.Parallel()
	//nolint:govet
	tests := []struct {
		name     string
		raw      *checker.RawResults
		outcomes []finding.Outcome
		err      error
	}{
		{
			name: "Three workflows none of which do untrusted checkout.",
			raw: &checker.RawResults{
				DangerousWorkflowResults: checker.DangerousWorkflowData{
					NumWorkflows: 3,
					Workflows: []checker.DangerousWorkflow{
						{
							Type: checker.DangerousWorkflowScriptInjection,
						},
					},
				},
			},
			outcomes: []finding.Outcome{
				finding.OutcomeFalse,
			},
		},
		{
			name: "Three workflows one of which has possibility of untrusted checkout.",
			raw: &checker.RawResults{
				DangerousWorkflowResults: checker.DangerousWorkflowData{
					NumWorkflows: 3,
					Workflows: []checker.DangerousWorkflow{
						{
							Type: checker.DangerousWorkflowUntrustedCheckout,
						},
					},
				},
			},
			outcomes: []finding.Outcome{
				finding.OutcomeTrue,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			findings, s, err := Run(tt.raw)
			if !cmp.Equal(tt.err, err, cmpopts.EquateErrors()) {
				t.Errorf("mismatch (-want +got):\n%s", cmp.Diff(tt.err, err, cmpopts.EquateErrors()))
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(Probe, s); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			test.AssertOutcomes(t, findings, tt.outcomes)
		})
	}
}
