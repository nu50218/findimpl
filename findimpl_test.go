package findimpl

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	target = "error"
	analysistest.Run(t, testdata, Analyzer, "a")
}
