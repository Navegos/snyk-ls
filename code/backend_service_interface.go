package code

import (
	sglsp "github.com/sourcegraph/go-lsp"

	"github.com/snyk/snyk-ls/lsp"
)

type BackendService interface {
	CreateBundle(files map[sglsp.DocumentURI]File) (string, []sglsp.DocumentURI, error)
	ExtendBundle(bundleHash string, files map[sglsp.DocumentURI]File, removedFiles []sglsp.DocumentURI) (string, []sglsp.DocumentURI, error)
	RetrieveDiagnostics(bundleHash string, limitToFiles []sglsp.DocumentURI, severity int) (map[sglsp.DocumentURI][]lsp.Diagnostic, map[sglsp.DocumentURI][]sglsp.CodeLens, string, error)
}
