package server

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-jsonnet/formatter"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

func (s *server) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	doc, err := s.cache.get(params.TextDocument.URI)
	if err != nil {
		err = fmt.Errorf("Formatting: %s: %w", errorRetrievingDocument, err)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	// TODO(#14): Formatting options should be user configurable.
	formatted, err := formatter.Format(params.TextDocument.URI.SpanURI().Filename(), doc.item.Text, formatter.DefaultOptions())
	if err != nil {
		err = fmt.Errorf("Formatting: unable to format document: %w", err)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	return getTextEdits(doc.item.Text, formatted), nil
}

func getTextEdits(before, after string) []protocol.TextEdit {
	edits := myers.ComputeEdits(span.URI("any"), before, after)

	var result []protocol.TextEdit
	for _, edit := range edits {
		result = append(result, protocol.TextEdit{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(edit.Span.Start().Line()) - 1, Character: uint32(edit.Span.Start().Column()) - 1},
				End:   protocol.Position{Line: uint32(edit.Span.End().Line()) - 1, Character: uint32(edit.Span.End().Column()) - 1},
			},
			NewText: edit.NewText,
		})
	}

	return result
}