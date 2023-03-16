package notification_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	lsp2 "github.com/snyk/snyk-ls/application/server/lsp"
	notification2 "github.com/snyk/snyk-ls/application/server/notification"
	"github.com/snyk/snyk-ls/domain/snyk"
	"github.com/snyk/snyk-ls/infrastructure/code"
	"github.com/snyk/snyk-ls/infrastructure/iac"
	"github.com/snyk/snyk-ls/internal/notification"
	"github.com/snyk/snyk-ls/internal/product"
	"github.com/snyk/snyk-ls/internal/testutil"
)

type sendMessageTestCase struct {
	name           string
	act            func(scanNotifier snyk.ScanNotifier)
	expectedStatus lsp2.ScanStatus
}

func Test_SendMessage(t *testing.T) {
	testutil.UnitTest(t)

	const folderPath = "/test/folderPath"

	tests := []sendMessageTestCase{
		{
			name: "SendInProgressMessage",
			act: func(scanNotifier snyk.ScanNotifier) {
				scanNotifier.SendInProgress(folderPath)
			},
			expectedStatus: lsp2.InProgress,
		},
		{
			name: "SendSuccessMessage",
			act: func(scanNotifier snyk.ScanNotifier) {
				scanNotifier.SendSuccess(folderPath, []snyk.Issue{})
			},
			expectedStatus: lsp2.Success,
		},
		{
			name: "SendErrorMessage",
			act: func(scanNotifier snyk.ScanNotifier) {
				scanNotifier.SendError(product.ProductCode, folderPath)
			},
			expectedStatus: lsp2.ErrorStatus,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedProduct := "code"
			mockNotifier := notification.NewMockNotifier()
			scanNotifier, _ := notification2.NewScanNotifier(mockNotifier)

			// Act - run the test
			test.act(scanNotifier)

			// Assert - search through all the messages for the expected message
			for _, msg := range mockNotifier.SentMessages() {
				if containsMatchingMessage(t, msg, test, expectedProduct, folderPath) {
					return
				}
			}
			assert.Fail(t, "Scan message was not sent")
		})
	}
}

func Test_SendSuccess_SendsForAllEnabledProducts(t *testing.T) {
	testutil.UnitTest(t)

	mockNotifier := notification.NewMockNotifier()
	scanNotifier, _ := notification2.NewScanNotifier(mockNotifier)

	const folderPath = "/test/iac/folderPath"

	// expected message uses lsp2.ScanIssue && lsp2.CodeIssueData
	expectedIacIssue := []lsp2.ScanIssue{
		{
			Id:       "098f6bcd4621d373cade4e832627b4f6",
			Title:    "iacTitle",
			Severity: "critical",
			FilePath: "iacAffectedFilePath",
			AdditionalData: lsp2.IacIssueData{
				PublicId:      "iacID",
				Documentation: "iacDocumentation",
				LineNumber:    1,
				Issue:         "iacIssue",
				Impact:        "iacImpact",
				Path:          []string{"iacPath"},
			},
		},
	}

	expectedCodeIssue := []lsp2.ScanIssue{
		{
			Id:       "5a105e8b9d40e1329780d62ea2265d8a",
			Title:    "codeMessage",
			Severity: "low",
			FilePath: "codeAffectedFilePath",
			AdditionalData: lsp2.CodeIssueData{
				Message:            "codeMessage",
				Rule:               "codeRule",
				RuleId:             "codeRuleID",
				RepoDatasetSize:    2,
				ExampleCommitFixes: []lsp2.ExampleCommitFix{},
				CWE:                []string{},
				IsSecurityType:     false,
				Text:               "codeText",
				Cols:               lsp2.Point{1, 1},
				Rows:               lsp2.Point{1, 1},
				Markers:            []lsp2.Marker{},
			},
		},
	}

	scanIssues := []snyk.Issue{
		{ // IaC issue
			ID:        "iacID",
			Severity:  snyk.Critical,
			IssueType: 1,
			Range: snyk.Range{
				Start: snyk.Position{
					Line:      1,
					Character: 1,
				},
				End: snyk.Position{
					Line:      1,
					Character: 2,
				},
			},
			Message:             "iacMessage",
			FormattedMessage:    "iacFormattedMessage",
			AffectedFilePath:    "iacAffectedFilePath",
			Product:             product.ProductInfrastructureAsCode,
			References:          []snyk.Reference{},
			IssueDescriptionURL: &url.URL{},
			CodeActions:         []snyk.CodeAction{},
			Commands:            []snyk.Command{},
			AdditionalData: iac.IssueData{
				Key:           "098f6bcd4621d373cade4e832627b4f6",
				Title:         "iacTitle",
				PublicId:      "iacID",
				Documentation: "iacDocumentation",
				LineNumber:    1,
				Issue:         "iacIssue",
				Impact:        "iacImpact",
				Path:          []string{"iacPath"},
			},
		},
		{ // Code issue
			ID:        "codeID",
			Severity:  snyk.Low,
			IssueType: 1,
			Range: snyk.Range{
				Start: snyk.Position{
					Line:      1,
					Character: 1,
				},
				End: snyk.Position{
					Line:      1,
					Character: 2,
				},
			},
			Message:             "codeMessage",
			FormattedMessage:    "codeFormattedMessage",
			AffectedFilePath:    "codeAffectedFilePath",
			Product:             product.ProductCode,
			References:          []snyk.Reference{},
			IssueDescriptionURL: &url.URL{},
			CodeActions:         []snyk.CodeAction{},
			Commands:            []snyk.Command{},
			AdditionalData: code.IssueData{
				Key:                "5a105e8b9d40e1329780d62ea2265d8a",
				Message:            "codeMessage",
				Rule:               "codeRule",
				RuleId:             "codeRuleID",
				RepoDatasetSize:    2,
				ExampleCommitFixes: []code.ExampleCommitFix{},
				CWE:                []string{},
				IsSecurityType:     false,
				Text:               "codeText",
				Cols:               code.CodePoint{1, 1},
				Rows:               code.CodePoint{1, 1},
				Markers:            []code.Marker{},
			},
		},
	}

	// Act - run the test
	scanNotifier.SendSuccess(folderPath, scanIssues)

	// Assert - check the messages matches the expected message for each product
	for _, msg := range mockNotifier.SentMessages() {
		if msg.(lsp2.SnykScanParams).Product == "code" {
			actualCodeIssue := msg.(lsp2.SnykScanParams).Issues
			assert.Equal(t, expectedCodeIssue, actualCodeIssue)
			return
		}
		if msg.(lsp2.SnykScanParams).Product == "iac" {
			actualIacIssue := msg.(lsp2.SnykScanParams).Issues
			assert.Equal(t, expectedIacIssue, actualIacIssue)
			return
		}
	}
}

func Test_NewScanNotifier_NilNotifier_Errors(t *testing.T) {
	t.Parallel()
	scanNotifier, err := notification2.NewScanNotifier(nil)
	assert.Error(t, err)
	assert.Nil(t, scanNotifier)
}

func Test_SendInProgress_SendsForAllEnabledProducts(t *testing.T) {
	testutil.UnitTest(t)

	// Arrange
	mockNotifier := notification.NewMockNotifier()
	scanNotifier, _ := notification2.NewScanNotifier(mockNotifier)

	// Act
	scanNotifier.SendInProgress("/test/folderPath")

	// Assert
	assert.Equal(t, 2, len(mockNotifier.SentMessages()))
}

func containsMatchingMessage(t *testing.T,
	msg any,
	testCase sendMessageTestCase,
	expectedProduct string,
	folderPath string,
) bool {
	t.Helper()
	scanMessage, ok := msg.(lsp2.SnykScanParams)
	if ok &&
		scanMessage.Status == testCase.expectedStatus &&
		scanMessage.Product == expectedProduct &&
		scanMessage.FolderPath == folderPath {
		return true
	}
	return false
}
