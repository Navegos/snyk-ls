/*
 * © 2024 Snyk Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package command

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/snyk/snyk-ls/domain/observability/performance"
	"github.com/snyk/snyk-ls/domain/snyk"
	"github.com/snyk/snyk-ls/infrastructure/code"
	"github.com/snyk/snyk-ls/infrastructure/snyk_api"
	"github.com/snyk/snyk-ls/internal/notification"
)

func Test_codeFixDiffs_Command(t *testing.T) {

}

type mockIssueProvider struct {
}

func (m mockIssueProvider) IssuesFor(_ string, _ snyk.Range) []snyk.Issue {
	panic("this should not be called")
}
func (m mockIssueProvider) Issue(id string) snyk.Issue {
	return snyk.Issue{ID: id}
}

func Test_codeFixDiffs_Execute(t *testing.T) {
	instrumentor := performance.NewInstrumentor()
	snykCodeClient := &code.FakeSnykCodeClient{
		UnifiedDiffSuggestions: []code.AutofixUnifiedDiffSuggestion{
			{
				FixId:               uuid.NewString(),
				UnifiedDiffsPerFile: nil,
			},
		},
	}
	snykApiClient := &snyk_api.FakeApiClient{CodeEnabled: true}
	codeScanner := &code.Scanner{
		BundleUploader: code.NewBundler(snykCodeClient, instrumentor),
		SnykApiClient:  snykApiClient,
	}
	cut := codeFixDiffs{
		notifier:    notification.NewMockNotifier(),
		codeScanner: codeScanner,
	}

	t.Run("no context", func(t *testing.T) {
		cut.issueProvider = mockIssueProvider{}
		_, err := cut.Execute(nil)
		require.Error(t, err)
	})

	t.Run("happy path", func(t *testing.T) {
		cut.issueProvider = mockIssueProvider{}
		codeScanner.BundleHashes = map[string]string{"folderPath": "bundleHash"}
		cut.command = snyk.CommandData{
			Arguments: []any{"folderPath", "issuePath", "issueId"},
		}

		suggestions, err := cut.Execute(context.Background())

		require.NotEmptyf(t, suggestions, "suggestions should not be empty")
		require.NoError(t, err)
	})
}
