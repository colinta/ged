package rule

// FlushRule is an optional interface for rules that buffer output
// and need to emit remaining lines at end-of-document.
// The caller checks for this with a type assertion and calls Flush
// once after all lines have been processed.
type FlushRule interface {
	Flush(ctx *LineContext) ([]string, error)
}
