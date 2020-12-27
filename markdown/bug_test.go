package markdown

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

type testRenderer struct {
	t *testing.T
}

func (t testRenderer) AddOptions(...renderer.Option) { return }
func (t testRenderer) Render(_ io.Writer, source []byte, n ast.Node) error {
	fencedCodeBlock := n.FirstChild().FirstChild().FirstChild().NextSibling().(*ast.FencedCodeBlock)

	line := fencedCodeBlock.Lines().At(0)
	if val := line.Value(source); !bytes.Equal([]byte("include .bingo/Variables.mk\n"), val) {
		t.t.Errorf("not what we expected, got %q", string(val))
	}
	line = fencedCodeBlock.Lines().At(1)
	if val := line.Value(source); !bytes.Equal([]byte("\n"), val) {
		t.t.Errorf("not what we expected, got %q", string(val)) // BUG1: bug_test.go:28: not what we expected, got "  \n"
	}
	line = fencedCodeBlock.Lines().At(2)
	if val := line.Value(source); !bytes.Equal([]byte("run:\n"), val) {
		t.t.Errorf("not what we expected, got %q", string(val))
	}
	line = fencedCodeBlock.Lines().At(3)
	if val := line.Value(source); !bytes.Equal([]byte("\t@$(GOIMPORTS) <args>\n"), val) {
		t.t.Errorf("not what we expected, got %q", string(val)) // BUG 2: bug_test.go:36: not what we expected, got "  @$(GOIMPORTS) <args>\n"	}
	}
	return nil
}

func TestGoldmarkCodeBlockWhitespaces(t *testing.T) {
	var codeBlock = "```"
	mdContent := []byte(fmt.Sprintf(`* Some item with nested code with strict whitespaces.
  %sMakefile
  include .bingo/Variables.mk
  
  run:
  	@$(GOIMPORTS) <args>
  %s`, codeBlock, codeBlock))

	var buf bytes.Buffer
	if err := goldmark.New(goldmark.WithRenderer(&testRenderer{t: t})).Convert(mdContent, &buf); err != nil {
		t.Fatal(err)
	}
}
