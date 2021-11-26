package bundle

import "github.com/sourcegraph/go-lsp"

var (
	DummyUri = lsp.DocumentURI("/Dummy.java")
)

type FakeBackendService struct {
	BundleHash string
}

func (f *FakeBackendService) createBundle(files map[lsp.DocumentURI]File) (string, []lsp.DocumentURI) {
	return f.BundleHash, nil
}
func (f *FakeBackendService) extendBundle(files map[lsp.DocumentURI]File, removedFiles []lsp.DocumentURI) []lsp.DocumentURI {
	return nil
}
func (f *FakeBackendService) retrieveDiagnostics() map[lsp.DocumentURI][]lsp.Diagnostic {
	diagnostic := lsp.Diagnostic{
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      2,
				Character: 5,
			},
			End: lsp.Position{
				Line:      2,
				Character: 7,
			},
		},
		Severity: lsp.Error,
		Code:     "123",
		Source:   "snyk code",
		Message:  "This is a dummy with severity Error",
	}

	fakes := map[lsp.DocumentURI][]lsp.Diagnostic{}
	var diagnostics []lsp.Diagnostic
	diagnostics = append(diagnostics, diagnostic)
	fakes[DummyUri] = diagnostics
	return fakes
}
