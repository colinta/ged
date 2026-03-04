# Go Migration Plan

This plan breaks down the ged project into incremental phases. Each phase introduces new Go concepts while building working, tested functionality.

You are a professional go developer and are teaching me the basics of Go by writing the 'ged' tool together. Before writing code, you should teach me about the library and concepts that we need for that section. Make sure I understand before we add more code to the project.

---

## Current Progress

| Phase | Status | Description |
|-------|--------|-------------|
| 1 | ✅ Complete | Basic substitution (`s/foo/bar`) |
| 2 | ✅ Complete | Filtering rules (`p/pattern/`, `d/pattern/`) |
| 3 | ✅ Complete | Rule chaining (multiple rules) |
| 4 | ✅ Complete | Line numbers (`p:1-5`, `d:2-4`, `s:1-3:replacement`) |
| 5 | ✅ Complete | Literal string matching (backtick/quote delimiters) |
| 6 | ✅ Complete | Document rules (`sort`, `reverse`, `join`) |
| 7 | ✅ Complete | Conditional rules (`if/pattern/ { rules }`) |
| 7b | ✅ Complete | LineContext refactor + control flow rules (`on/off/after/toggle`) |
| 8 | ✅ Complete | Between condition (`between/start/end/ { rules }`) |
| 9 | ✅ Complete | File I/O (`--input`, `--write`, `--write-to`, `--`) |
| 10 | ✅ Complete | Text modification (`trim`, `upper`, `lower`, `prepend`, `append`, `surround`) |
| 11 | ✅ Complete | Column operations (`cols`) |
| 12 | ✅ Complete | Extraction rules |
| 13 | ✅ Moved to 7b | Control flow rules (done early, needed LineContext) |
| 14 | ✅ Complete | External commands (`xargs`, `exec`) |
| 15-20 | 🔲 Pending | Diff/colors, more rules, polish |

**To continue**: Run `go test ./...` to verify everything works, then start Phase 15.

## Lesson Files

Each phase introduces Go concepts. After completing a phase, create a lesson file in `lessons/` for each new concept covered. Lesson files are numbered sequentially across all phases (not per-phase). Each file should be 1–3 paragraphs with example code showing the concept in context. Follow the existing format in `lessons/`.

**Current lessons** (from Phases 1–9):
| File | Concept |
|------|---------|
| `01-packages-and-modules.md` | Package structure, `go mod init`, `internal/` |
| `02-basic-types-and-errors.md` | Types, error interface, `%w` wrapping |
| `03-interfaces.md` | Implicit conformance, small interfaces |
| `04-io-reader-writer.md` | `io.Reader`/`io.Writer` for testable I/O |
| `05-regexp.md` | `regexp` package, `QuoteMeta` |
| `06-functional-options.md` | Variadic option functions for constructors |
| `07-table-driven-tests.md` | `testing` package, `t.Run`, sub-tests |
| `08-slices.md` | Dynamic arrays, `append`, nil vs empty |
| `09-custom-types-with-methods.md` | Named types, method receivers |
| `10-strconv.md` | `strconv.Atoi`, string↔int conversion |
| `11-strings-builder.md` | `strings.Builder`, `Join`, `Split`, `TrimSpace` |
| `12-sort-and-slices.md` | In-place sort/reverse, copy-first pattern |
| `13-type-switches-and-any.md` | `any` return type, type switch dispatch |
| `14-recursive-parsing.md` | Recursive descent, token threading |
| `15-iota-enums.md` | `iota` constants, zero-value conventions |
| `16-optional-interfaces.md` | Separate interfaces, type assertion checks |
| `17-defer.md` | Deferred cleanup, LIFO order |
| `18-file-io.md` | `os.Open`/`os.Create`, atomic write-back |
| `19-shared-mutable-context.md` | Pointer-shared state across pipeline |
| `20-string-transforms.md` | `TrimSpace`, `ToUpper`, `ToLower`, simple concatenation |
| `21-regex-split.md` | Splitting strings by regex pattern with match iteration |
| `22-index-resolution.md` | 1-based/negative index resolution, `(int, bool)` pattern |
| `23-regex-match-groups.md` | `FindStringMatch`, `Groups()`, iterating matches |

| `24-os-exec.md` | `os/exec`, `exec.Command`, shell quoting, stdin/stdout wiring |

**Next lesson number**: 25

## Project Structure

```
ged/
├── cmd/
│   └── ged/
│       ├── main.go              # CLI entry point
│       └── main_test.go         # CLI integration tests
├── internal/
│   ├── rule/
│   │   ├── rule.go              # LineRule, DocumentRule, LineContext, SetupRule, PrintState
│   │   ├── sub_line_rule.go     # SubstitutionRule (pattern-based)
│   │   ├── sub_linenum_rule.go  # SubLineNumRule (line number-based)
│   │   ├── print_line_rule.go   # PrintLineRule (pattern-based)
│   │   ├── delete_line_rule.go  # DeleteLineRule (pattern-based)
│   │   ├── print_linenum_rule.go # PrintLineNumRule (line number-based)
│   │   ├── delete_linenum_rule.go # DeleteLineNumRule (line number-based)
│   │   ├── linerange.go         # LineRange types for line number parsing
│   │   ├── sort_rule.go         # SortRule (document rule)
│   │   ├── reverse_rule.go      # ReverseRule (document rule)
│   │   ├── join_rule.go         # JoinRule (document rule)
│   │   ├── apply_all_rule.go    # ApplyAllRule (wraps LineRules into DocumentRule)
│   │   ├── conditional_rule.go  # ConditionalLineRule and ConditionalDocRule
│   │   ├── on_rule.go           # OnRule (control flow)
│   │   ├── off_rule.go          # OffRule (control flow)
│   │   ├── after_rule.go        # AfterRule (control flow)
│   │   ├── toggle_rule.go       # ToggleRule (control flow)
│   │   ├── trim_rule.go         # TrimRule (both/left/right)
│   │   ├── case_rule.go         # UpperRule, LowerRule
│   │   ├── text_rule.go         # PrependRule, AppendRule, SurroundRule
│   │   ├── columns_rule.go     # ColumnsRule, ColumnSpec, regexSplit
│   │   ├── take_rule.go         # TakeRule (extract matching portion)
│   │   ├── remove_rule.go       # RemoveRule (remove matching portion)
│   │   ├── group_rule.go        # GroupRule (extract capture group)
│   │   ├── xargs_rule.go       # XargsRule (execute command per line)
│   │   ├── exec_rule.go        # ExecRule (pipe document through command)
│   │   └── *_test.go            # Tests for each
│   ├── parser/
│   │   ├── parser.go            # Rule parsing (single rules + control rules)
│   │   ├── parse_args.go        # Multi-arg parsing with { } blocks
│   │   └── *_test.go
│   └── engine/
│       ├── pipeline.go          # Processing pipeline
│       └── pipeline_test.go
├── go.mod
└── go.sum
```

---

## Phase 1: Hello Go - Basic Substitution ✅ COMPLETE

**Goal**: Get a working `ged 's/foo/bar'` that reads stdin and writes stdout.

**Go Concepts Learned**:
- Package structure and `go mod init`
- Basic types: strings, errors
- `fmt` and `os` packages
- `bufio.Scanner` for line reading
- `regexp` package
- Writing and running tests with `go test`
- **Functional options pattern** for configurable constructors
- **Table-driven tests** for comprehensive test coverage
- **Implicit interface conformance** (no explicit `implements`)
- **Multiple return values** for error handling
- **`strings.Builder`** for efficient string building
- **`io.Reader` and `io.Writer`** interfaces for testable I/O
- **Error wrapping** with `fmt.Errorf("...: %w", err)`

### Implementation Notes

**Functional Options**: We use the idiomatic Go pattern for optional parameters:
```go
rule, _ := NewSubstitutionRule("foo", "bar")              // defaults
rule, _ := NewSubstitutionRule("foo", "bar", WithGlobal()) // with options
```

**Parser Design**: Two-layer parsing:
- `ParseRule()` - handles delimiter detection, escape sequences, dispatches to command parsers
- `parseSubstitution()` - validates and creates SubstitutionRule

**Flexible Syntax**: Trailing delimiter is optional unless flags are needed:
- `s/foo/bar` ✓
- `s/foo/bar/` ✓
- `s/foo/bar/g` ✓ (need delimiter before flags)
- `s/foo/` ✓ (empty replacement)

**Escape Handling**: `splitByDelimiter()` handles `\/` and `\\` escape sequences.

**Testable CLI**: `main()` is a thin wrapper that calls `run(args, stdin, stdout, stderr)`. The `run()` function accepts `io.Reader`/`io.Writer` interfaces, allowing tests to use `strings.NewReader` and `bytes.Buffer` instead of real I/O.

### Tests Written (36 total)
- [x] SubstitutionRule replaces first match only
- [x] SubstitutionRule handles no match (returns original)
- [x] SubstitutionRule handles regex patterns
- [x] SubstitutionRule with `WithGlobal()` replaces all matches
- [x] SubstitutionRule handles capture group replacements ($1, $2)
- [x] Invalid regex returns error
- [x] Parser handles various delimiters (/, |, =, #)
- [x] Parser handles escaped delimiters
- [x] Parser handles escaped backslashes
- [x] Parser preserves whitespace
- [x] Parser rejects invalid input
- [x] CLI handles basic substitution end-to-end
- [x] CLI handles multiple lines
- [x] CLI returns errors for invalid input

### Files Created
- `internal/rule/rule.go` - Rule interface
- `internal/rule/line_rules.go` - SubstitutionRule with functional options
- `internal/rule/line_rules_test.go` - Rule tests (7 tests)
- `internal/parser/parser.go` - ParseRule with escape handling
- `internal/parser/parser_test.go` - Parser tests, table-driven (21 tests)
- `cmd/ged/main.go` - CLI entry point with testable `run()` function
- `cmd/ged/main_test.go` - CLI integration tests (8 tests)

### Deliverable ✅
```bash
echo "hello world" | ./ged 's/world/earth'
# Output: hello earth

echo "hello world world" | ./ged 's/world/earth'
# Output: hello earth world  (first match only)

echo "hello world world" | ./ged 's/world/earth/g'
# Output: hello earth earth  (global)

echo "foo 123 bar 456" | ./ged 's/\d+/NUM/g'
# Output: foo NUM bar NUM
```

---

## Phase 2: Filtering Rules ✅ COMPLETE

**Goal**: Implement `p/pattern/` (print matching) and `d/pattern/` (delete matching).

**Go Concepts Learned**:
- **Empty slice vs nil semantics**: `[]string{}` signals "delete line", slice with content keeps line(s)
- Separate files per rule type for better organization

### Implementation Notes

- `PrintLineRule` - keeps lines matching pattern, deletes non-matching
- `DeleteLineRule` - deletes lines matching pattern, keeps non-matching
- Parser extended with `parsePrint()` and `parseDelete()` functions

### Tests Written
- [x] PrintLineRule keeps matching lines
- [x] PrintLineRule removes non-matching lines
- [x] DeleteLineRule removes matching lines
- [x] DeleteLineRule keeps non-matching lines
- [x] Regex patterns work in both rules
- [x] Different delimiters parse correctly

### Deliverable ✅
```bash
echo -e "foo\nbar\nfoo" | ged 'p/foo/'
# Output: foo\nfoo

echo -e "foo\nbar\nfoo" | ged 'd/foo/'
# Output: bar
```

---

## Phase 3: Rule Chaining ✅ COMPLETE

**Goal**: Support multiple rules: `ged 'p/foo/' 's/o/x/'`

**Go Concepts Learned**:
- **Slices**: Dynamic arrays with `append()` - always reassign result
- **Variadic functions**: `func NewPipeline(rules ...Rule)` accepts any number of arguments
- **Spread operator**: `rules...` to pass a slice as variadic arguments

### Implementation Notes

- `Pipeline` type chains multiple rules together
- Each rule's output feeds into the next rule
- Empty output stops the chain (for filtering)
- CLI updated to parse multiple rule arguments

### Tests Written
- [x] Two rules chain correctly
- [x] Filter then substitute works
- [x] Substitute then filter works
- [x] Empty output stops the chain
- [x] Delete rule in chain works

### Deliverable ✅
```bash
echo -e "hello\nworld\nhello" | ged 'p/hello/' 's/o/x/'
# Output: hellx\nhellx
```

---

## Phase 4: Line Numbers ✅ COMPLETE

**Goal**: Support line number operations: `p:1-5`, `d:2-4`

**Go Concepts Learned**:
- **Custom types with methods**: `type SingleLine int` with `Contains(lineNum int) bool`
- **Parsing with `strconv`**: `strconv.Atoi()` for string-to-int conversion
- **Interface polymorphism**: `LineRange` interface with multiple implementations
- **Breaking change management**: Updated `Rule.Apply()` signature to include `lineNum`

### Implementation Notes

**Rule Interface Change**: All rules now receive a context carrying the line number (refactored from bare `lineNum int` in Phase 7b):
```go
type LineRule interface {
    Apply(line string, ctx *LineContext) ([]string, error)
}
```

**LineRange Types** (in `internal/rule/linerange.go`):
- `SingleLine` - matches one line: `5`
- `Range` - matches range: `2-4`
- `OpenRange` - matches open-ended: `5-` or `-5`
- `CompositeRange` - combines with OR: `1,3,5-7`

**Colon Delimiter**: `:` indicates line number rules vs `/` for pattern rules:
- `p:2-4` → PrintLineNumRule (lines 2, 3, 4)
- `d:2-4` → DeleteLineNumRule (remove lines 2, 3, 4)
- `s:2-4:text` → SubLineNumRule (replace lines 2, 3, 4 with "text")
- `p/foo/` → PrintLineRule (pattern match)

**Parser Refactor**: `ParseRule` uses `if/else if` with compound conditions to dispatch based on both command and delimiter. Specific cases (e.g. `command == 'p' && delimiter == ':'`) come before general cases (e.g. `command == 'p'`). Parse functions no longer receive the delimiter parameter.

### Tests Written
- [x] Single line number matches correctly
- [x] Range `2-4` matches lines 2, 3, 4
- [x] Open range `5-` matches 5 and beyond
- [x] Open range `-5` matches 1 through 5
- [x] Comma-separated ranges work
- [x] PrintLineNumRule filters by line number
- [x] DeleteLineNumRule filters by line number
- [x] SubLineNumRule replaces matching lines
- [x] SubLineNumRule keeps non-matching lines
- [x] SubLineNumRule with newline in replacement returns multiple lines

### Deliverable ✅
```bash
echo -e "1\n2\n3\n4\n5" | ged 'p:2-4'
# Output: 2\n3\n4

echo -e "one\ntwo\nthree" | ged 's:2:replaced'
# Output: one\nreplaced\nthree
```

---

## Phase 5: Literal String Matching ✅ COMPLETE

**Goal**: Support quote delimiters for literal matching

**Go Concepts Learned**:
- **`regexp.QuoteMeta`**: Escapes all regex metacharacters in a string
- **Escape sequences in `splitByDelimiter`**: `\n` → newline, `\t` → tab
- **`strings.Split`**: Splitting substitution results on newlines to produce multiple output lines

### Implementation Notes

**Literal Matching**: When the delimiter is a quote character (`` ` ``, `'`, `"`), the pattern is run through `regexp.QuoteMeta` before being compiled as a regex. This happens centrally in `ParseRule` before dispatching to parse functions.

**Escape Sequences**: `splitByDelimiter` now expands `\n` and `\t` in addition to `\\` and escaped delimiters. This works in both patterns and replacements.

**Newline in Replacements**: `SubstitutionRule.Apply` and `SubLineNumRule.Apply` split results on `\n` and return multiple entries, so a replacement containing `\n` produces multiple output lines.

### Tests Written
- [x] Backtick treats `.` as literal dot
- [x] Backtick treats `[` `]` as literal brackets
- [x] Single quote activates literal matching
- [x] Double quote activates literal matching
- [x] Escape sequences expand correctly (`\n`, `\t`)
- [x] Newline in substitution replacement produces multiple output lines
- [x] QuoteMeta'd pattern matches literal but not regex wildcards

### Deliverable ✅
```bash
echo "foo.bar" | ged 's`foo.bar`baz`'
# Output: baz  (literal match, not regex)

echo "foo.bar" | ged "s'foo.bar'baz'"
# Output: baz

echo "hello" | ged 's/hello/line1\nline2/'
# Output: line1
#         line2
```

---

## Phase 6: Document Rules ✅ COMPLETE

**Goal**: Implement `sort`, `reverse`, `join`

**Go Concepts Learned**:
- **`sort.Strings`**: Sorts a string slice in place — always copy first to avoid mutating the caller's data
- **`slices.Reverse`**: Reverses a slice in place (Go 1.21+, `slices` package)
- **`strings.Join`**: Joins slice elements with a separator string
- **Type switches**: `switch r := parsed.(type) { case X: ... }` dispatches on runtime type
- **`any` type**: Alias for `interface{}`, used when a function returns different interface types
- **Circular import avoidance**: Go forbids circular imports; `ApplyAllRule` inlines pipeline logic to avoid `rule` importing `engine`

### Implementation Notes

**Architecture Change**: Renamed `Rule` to `LineRule` (per-line processing) and added `DocumentRule` (whole-document processing). The rename is transparent to existing code because Go uses implicit interface conformance.

**Two Interfaces** (note: `lineNum int` was later refactored to `*LineContext` in Phase 7b):
```go
type LineRule interface {
    Apply(line string, ctx *LineContext) ([]string, error)
}
type DocumentRule interface {
    ApplyDocument(lines []string) ([]string, error)
}
```

**Parser Returns `any`**: `ParseRule` now returns `(any, error)` because it can produce either a `LineRule` or a `DocumentRule`. Word commands (`sort`, `reverse`, `join`) are checked *before* single-character command dispatch, since `sort` starts with `s` and would otherwise match the substitution command.

**ApplyAllRule**: Wraps consecutive `LineRule`s into a `DocumentRule` by inlining the pipeline chaining logic. This avoids a circular import between `rule` and `engine`.

**main.go Rewrite**: `run()` now:
1. Parses all args, building a `[]DocumentRule` list
2. Consecutive `LineRule`s are wrapped in `ApplyAllRule`
3. All stdin is buffered into `[]string`
4. Each `DocumentRule` is applied in sequence
5. Output is written

### Tests Written
- [x] Sort orders alphabetically
- [x] Sort handles empty/single-line input
- [x] Sort does not mutate input slice
- [x] Reverse reverses line order
- [x] Reverse handles empty/single-line input
- [x] Reverse does not mutate input slice
- [x] Join combines lines with comma
- [x] Join combines lines with space
- [x] Join combines lines with empty separator
- [x] Join handles empty/single-line input
- [x] ApplyAllRule applies substitution to all lines
- [x] ApplyAllRule filters lines
- [x] ApplyAllRule chains multiple rules
- [x] ApplyAllRule preserves line numbering
- [x] Parser parses `sort`, `reverse`, `join`, `join/,/`
- [x] `sort` does not match as substitution command
- [x] CLI: sort, reverse, join end-to-end
- [x] CLI: line rules then sort
- [x] CLI: sort then line rules
- [x] CLI: bare join (empty separator)

### Files Created
- `internal/rule/sort_rule.go` - SortRule (DocumentRule)
- `internal/rule/reverse_rule.go` - ReverseRule (DocumentRule)
- `internal/rule/join_rule.go` - JoinRule (DocumentRule)
- `internal/rule/apply_all_rule.go` - ApplyAllRule (wraps LineRules into DocumentRule)
- `internal/rule/sort_rule_test.go` - Tests
- `internal/rule/reverse_rule_test.go` - Tests
- `internal/rule/join_rule_test.go` - Tests
- `internal/rule/apply_all_rule_test.go` - Tests
- `internal/parser/parse_document_test.go` - Tests

### Files Modified
- `internal/rule/rule.go` - Renamed `Rule` to `LineRule`, added `DocumentRule`
- `internal/engine/pipeline.go` - `rule.Rule` → `rule.LineRule`
- `internal/parser/parser.go` - Returns `any`, word command dispatch, helper return types
- `cmd/ged/main.go` - Rewritten for document-rule architecture

### Deliverable ✅
```bash
echo -e "c\na\nb" | ged sort
# Output: a\nb\nc

echo -e "a\nb\nc" | ged reverse
# Output: c\nb\na

echo -e "a\nb\nc" | ged 'join/,/'
# Output: a,b,c

echo -e "c3\na1\nb2" | ged 's/[0-9]//g' sort
# Output: a\nb\nc
```

---

## Phase 7: Conditional Rules ✅ COMPLETE

**Goal**: Implement `if/pattern/ { rules }`

**Go Concepts Learned**:
- **Recursive parsing**: `parseArgs` calls itself to handle nested `{ }` blocks
- **Intermediate types**: `condition` is a parser-internal type that bridges parsing and rule creation
- **Tree structures**: `ConditionalLineRule` contains child rules, forming a tree instead of a flat list
- **`make([]bool, n)`**: Pre-allocated boolean slice for tracking match positions

### Implementation Notes

**Two Conditional Rule Types**:
- `ConditionalLineRule` (implements `LineRule`) — all inner rules are `LineRule`s, can stream line-by-line
- `ConditionalDocRule` (implements `DocumentRule`) — inner rules include `DocumentRule`s, buffers matching lines into a sub-document

The parser decides which to create based on what's inside the block.

**Syntax**: Each token is a separate CLI argument:
```bash
ged 'if/hello/' '{' 's/o/x/' '}'
ged '!if/hello/' '{' 's/o/x/' '}'      # inverted
ged 'if/foo/' '{' 'if/bar/' '{' 's/x/y/' '}' '}'  # nested
ged 'if/item/' '{' 'sort' '}'           # document rule inside block
```

**ParseArgs**: New top-level parser function that handles multi-arg block syntax. Replaces the per-arg loop in `main.go`. Uses recursion to handle nested blocks.

**ConditionalDocRule semantics**: Matching lines are collected into a sub-document, inner rules are applied, then results are woven back into their original positions. Non-matching lines stay in place as fixed anchors.

**buildDocRules helper**: Converts mixed `[]any` into `[]DocumentRule` by wrapping consecutive `LineRule`s in `ApplyAllRule`. Same logic as `main.go`'s rule grouping.

### Tests Written
- [x] If condition applies rules to matches only
- [x] Inverted if applies to non-matches
- [x] Non-matching lines pass through unchanged
- [x] Multiple inner rules chain as pipeline
- [x] Inner delete removes line entirely
- [x] Line numbers passed through to inner rules
- [x] ConditionalDocRule sorts only matching lines
- [x] ConditionalDocRule reverses only matching lines
- [x] ConditionalDocRule joins matching lines
- [x] ConditionalDocRule inverted works
- [x] ConditionalDocRule no matches passes through
- [x] ConditionalDocRule mixed line+doc inner rules
- [x] Parser: if/pattern/ returns condition
- [x] Parser: !if/pattern/ returns inverted condition
- [x] Parser: literal delimiters work with if
- [x] Parser: missing pattern errors
- [x] ParseArgs: simple rules without conditionals
- [x] ParseArgs: conditional block creates LineRule
- [x] ParseArgs: document rule inside block creates DocRule
- [x] ParseArgs: nested conditionals
- [x] ParseArgs: error on missing braces
- [x] CLI: if condition end-to-end
- [x] CLI: inverted if end-to-end
- [x] CLI: if with multiple inner rules
- [x] CLI: if then sort
- [x] CLI: if with document rule inside block
- [x] CLI: nested if (chained conditions)

### Files Created
- `internal/rule/conditional_rule.go` - ConditionalLineRule and ConditionalDocRule
- `internal/rule/conditional_rule_test.go` - Tests
- `internal/parser/parse_args.go` - ParseArgs with recursive block parsing
- `internal/parser/parse_args_test.go` - Tests

### Files Modified
- `internal/parser/parser.go` - Added `parseIf()`, `condition` type, `if`/`!if` dispatch
- `cmd/ged/main.go` - Uses `ParseArgs()` instead of per-arg loop
- `cmd/ged/main_test.go` - CLI integration tests for conditionals

### Deliverable ✅
```bash
echo -e "hello\nworld\nhello" | ged 'if/hello/' '{' 's/o/x/' '}'
# Output: hellx\nworld\nhellx

echo -e "hello\nworld\nhello" | ged '!if/hello/' '{' 's/o/x/' '}'
# Output: hello\nwxrld\nhello

echo -e "b_item\na_item\nc_other\nd_item" | ged 'if/item/' '{' sort '}'
# Output: a_item\nb_item\nc_other\nd_item
```

---

## Phase 7b: LineContext Refactor + Control Flow Rules ✅ COMPLETE

**Goal**: Refactor `Apply` signature to use shared context, then implement `on/off/after/toggle` print-control rules (originally Phase 13, moved up because control flow motivates the context design).

**Go Concepts Learned**:
- **`iota` for enums**: Constants auto-increment from 0 in a `const` block — zero value conventionally means "unset/default"
- **Optional interfaces**: `SetupRule` is a separate interface; the caller uses type assertion `if s, ok := r.(SetupRule); ok { ... }` to call it only on rules that implement it
- **Shared mutable state via context**: Multiple rules read/write `ctx.Printing` instead of maintaining rule-local state
- **Self-initializing Setup**: Each `Setup` method guards with `if ctx.Printing == PrintDefault` so the first control rule in the pipeline determines the starting state

### Implementation Notes

**LineContext Refactor**: Replaced `lineNum int` parameter with `*LineContext` across all `LineRule.Apply` signatures. Used `ged` itself (with backtick literal matching) to perform the mechanical refactor across test files.

**PrintState Enum**:
```go
type PrintState int
const (
    PrintDefault PrintState = iota  // 0 — no control rule, print everything
    PrintOn                         // 1 — printing enabled
    PrintOff                        // 2 — printing suppressed
)
```

**Control Rules Don't Filter**: They set `ctx.Printing` but always return `[]string{line}`. The caller (main.go streaming loop or ApplyAllRule) checks `ctx.Printing` after processing each line and decides whether to include it in output. This means other rules in the pipeline still see every line.

**SetupRule**: Optional interface called once before the processing loop to set initial `PrintState`. Guards with `PrintDefault` check so multiple control rules don't clobber each other — first rule wins.

**AfterRule Local State**: Uses rule-local `matched bool` in addition to shared `ctx.Printing`. Checks `r.matched` before checking the pattern, so the matching line itself stays off and the next line turns on.

### Semantics

| Rule | Initial state | Match line printed? | Lines after match |
|------|--------------|--------------------|--------------------|
| `on/pat/` | off | yes | on |
| `off/pat/` | on | no | off |
| `after/pat/` | off | no | on |
| `toggle/pat/` | off | flips | flipped |

### Tests Written
- [x] OnRule starts at match, includes match line
- [x] OnRule with no match prints nothing
- [x] OffRule stops at match, excludes match line
- [x] OffRule with no match prints everything
- [x] AfterRule starts after match, excludes match line
- [x] AfterRule with no match prints nothing
- [x] ToggleRule flips on each match
- [x] ToggleRule match line follows new state
- [x] On + Off combined (first rule sets initial state)
- [x] Parser: on/pattern/, off/pattern/, after/pattern/, toggle/pattern/
- [x] Parser: literal delimiters work with control rules
- [x] Parser: missing/empty pattern errors
- [x] CLI: on, off, after, toggle end-to-end
- [x] CLI: on with substitution
- [x] CLI: on + off combined

### Files Created
- `internal/rule/on_rule.go` - OnRule (SetupRule + LineRule)
- `internal/rule/off_rule.go` - OffRule (SetupRule + LineRule)
- `internal/rule/after_rule.go` - AfterRule (SetupRule + LineRule)
- `internal/rule/toggle_rule.go` - ToggleRule (SetupRule + LineRule)
- `internal/rule/control_rule_test.go` - Tests
- `internal/parser/parse_control_test.go` - Tests

### Files Modified
- `internal/rule/rule.go` - Added LineContext, PrintState, SetupRule; updated LineRule.Apply signature
- `internal/rule/apply_all_rule.go` - Calls Setup, checks ctx.Printing
- `internal/rule/*.go` - All rule Apply signatures updated (lineNum int → ctx *LineContext)
- `internal/engine/pipeline.go` - Process signature updated
- `internal/parser/parser.go` - Added parseControl, on/off/after/toggle dispatch
- `cmd/ged/main.go` - Calls Setup, checks ctx.Printing in streaming path
- All `*_test.go` files - Updated Apply/Process calls to use &LineContext{}

### Deliverable ✅
```bash
echo -e "a\nstart\nb\nc" | ged 'on/start/'
# Output: start\nb\nc

echo -e "a\nb\nstop\nc" | ged 'off/stop/'
# Output: a\nb

echo -e "a\nmarker\nb\nc" | ged 'after/marker/'
# Output: b\nc

echo -e "off\n---\non1\non2\n---\noff2" | ged 'toggle/---/'
# Output: ---\non1\non2

echo -e "before\nstart\nmiddle\nend\nafter" | ged 'on/start/' 'off/end/'
# Output: start\nmiddle
```

---

## Phase 8: Between Condition

**Goal**: Implement `between/start/end/ { rules }`

**Go Concepts Introduced**:
- Stateful rule processing
- Range tracking
- Edge case handling (inclusive/exclusive)

### Steps

1. **Create BetweenCondition**
   - Track state: before, inside, after
   - Apply rules only when inside

2. **Handle inclusive boundaries**
   - Start line is inside the range
   - End line is inside the range

3. **Support inverted between**
   - `!between` applies rules outside the range

### Tests to Write
- [ ] Rules apply inside range
- [ ] Start line is included
- [ ] End line is included
- [ ] Rules don't apply outside range
- [ ] Inverted between works
- [ ] Multiple ranges in one document work
- [ ] Nested between conditions work

### Deliverable
```bash
echo -e "start\n1\n2\nend\n3" | ged 'between/start/end/ { s/\d/x }'
# Output: start\nx\nx\nend\n3
```

---

## Phase 9: File I/O ✅ COMPLETE

**Goal**: Support `--input=file` and `--write`

**Go Concepts Learned**:
- **`os.Open` / `os.Create`**: File opening returns `*os.File` which implements `io.Reader`/`io.Writer`
- **`defer`**: Schedules cleanup (like `f.Close()`) to run when the function returns
- **`os.CreateTemp` + `os.Rename`**: Safe atomic write pattern for in-place editing
- **`os.Stat` / `os.Chmod`**: Preserving file permissions during write-back
- **`filepath.Dir`**: Extracting directory from a path for temp file placement
- **`strings.HasPrefix`**: Simple flag parsing without the `flag` package

### Implementation Notes

**CLI Options Parsing**: `parseCliOptions` separates flags from rule args:
- `--input=FILE` — read from file (repeatable for multiple files)
- `--write` — overwrite input files in place (requires `--input`)
- `--write-to=FILE` — write output to specific file
- `--` — everything after is a rule argument (not a flag)

**Architecture Refactor**: Extracted `processStream()` as the core processing function that takes `io.Reader`/`io.Writer`. Both stdin and file inputs use this same function. `processFile()` handles file opening, and `writeBack()` handles the temp-file-then-rename pattern.

**Atomic Write**: `writeBack` writes to a temp file in the same directory as the target, then uses `os.Rename` for an atomic swap. This prevents data loss if the program crashes mid-write. File permissions are preserved via `os.Stat`/`os.Chmod`.

**Validation**: Mutually exclusive flags are caught early:
- `--write` requires `--input`
- `--write` and `--write-to` can't be used together
- `--write-to` can't be used with multiple input files

### Tests Written
- [x] Read from file with `--input=`
- [x] File not found returns error
- [x] Write in-place with `--write`
- [x] Write-back preserves file permissions
- [x] `--write-to=` writes to specific file, leaves original unchanged
- [x] `--write-to=` works with stdin (no `--input`)
- [x] `--write` without `--input` errors
- [x] `--write` + `--write-to` errors (mutually exclusive)
- [x] Multiple input files output to stdout
- [x] Multiple input files with `--write` updates each
- [x] Multiple input files with `--write-to` errors
- [x] Bare rules after `--` (not treated as flags)
- [x] File input with document rules (sort)

### Files Modified
- `cmd/ged/main.go` — Added `cliOptions`, `parseCliOptions`, `processStream`, `processFile`, `writeBack`, `readLines`, `applyDocRules`; refactored `run()` to use them
- `cmd/ged/main_test.go` — Added 13 file I/O tests

### Deliverable ✅
```bash
ged 's/foo/bar' --input=test.txt --write
# Transforms test.txt in place

ged 's/foo/bar' --input=test.txt --write-to=output.txt
# Writes to output.txt, leaves test.txt unchanged

ged 's/foo/bar' --input=a.txt --input=b.txt
# Processes both files, outputs to stdout

ged -- 's/--flag/value/'
# Everything after -- is a rule, even if it looks like a flag
```

---

## Phase 10: Text Modification Rules ✅ COMPLETE

**Goal**: Implement `trim`, `triml`, `trimr`, `upper`, `lower`, `prepend`, `append`, `surround`

**Go Concepts Learned**:
- **`strings.TrimSpace`**: Removes all Unicode whitespace from both ends
- **`strings.TrimLeft` / `strings.TrimRight`**: Removes characters from a cutset (set of chars, not substring)
- **`strings.ToUpper` / `strings.ToLower`**: Unicode-aware case conversion
- **Simple concatenation**: `prefix + line` is idiomatic for small string joins

### Implementation Notes

**Word Commands**: `trim`, `triml`, `trimr`, `upper`, `lower` are bare word commands (like `sort`, `reverse`) — no delimiter or arguments needed. The parser matches the exact word.

**Text Commands**: `prepend`, `append`, `surround` take text arguments via delimiters (like `join`). A shared `parseTextCommand` helper handles delimiter parsing and dispatches by name.

**TrimRule modes**: A single `TrimRule` struct with a `mode` field ("both", "left", "right") instead of three separate types. The three constructors `NewTrimRule`, `NewTrimLeftRule`, `NewTrimRightRule` set the mode.

### Tests Written
- [x] Trim removes whitespace from both ends
- [x] Trim handles tabs, mixed whitespace, internal spaces
- [x] TrimLeft removes leading whitespace only
- [x] TrimRight removes trailing whitespace only
- [x] Upper converts to uppercase
- [x] Lower converts to lowercase
- [x] Prepend adds text before each line
- [x] Append adds text after each line
- [x] Surround wraps with before and after text
- [x] Parser: trim, triml, trimr, upper, lower (word commands)
- [x] Parser: prepend, append, surround (with delimiters)
- [x] Parser: missing argument errors
- [x] CLI: trim, triml, trimr end-to-end
- [x] CLI: upper, lower end-to-end
- [x] CLI: prepend, append, surround end-to-end
- [x] CLI: trim then upper (chained)
- [x] CLI: if condition with prepend

### Files Created
- `internal/rule/trim_rule.go` — TrimRule with both/left/right modes
- `internal/rule/case_rule.go` — UpperRule and LowerRule
- `internal/rule/text_rule.go` — PrependRule, AppendRule, SurroundRule
- `internal/rule/text_rule_test.go` — Rule tests
- `internal/parser/parse_text_test.go` — Parser tests

### Files Modified
- `internal/parser/parser.go` — Added word command dispatch + `parseTextCommand` helper
- `cmd/ged/main_test.go` — CLI integration tests

### Deliverable ✅
```bash
echo "  hello  " | ged trim
# Output: hello

echo "hello" | ged upper
# Output: HELLO

echo "hello" | ged 'prepend/>> /'
# Output: >> hello

echo "hello" | ged 'surround/[/]/'
# Output: [hello]

echo "  hello  " | ged trim upper
# Output: HELLO
```

---

## Phase 11: Column Operations ✅ COMPLETE

**Goal**: Implement `cols/pattern/spec` for column selection and reordering

**Go Concepts Learned**:
- **Regex splitting**: Iterating `FindStringMatch`/`FindNextMatch` to split strings by regex pattern
- **Index resolution**: Converting 1-based/negative indices to 0-based, `(int, bool)` validity pattern
- **`strconv.Atoi`**: Parsing column spec strings into integers

### Implementation Notes

**Syntax**: `cols/pattern/spec` or `cols/pattern/spec/joiner`
- Pattern is a regex (empty = `\s+` for whitespace splitting)
- Spec is comma-separated: `1,3-5,-1` (1-based, negatives from end, ranges inclusive)
- Joiner defaults to `" "` (single space)
- Quote delimiters for literal split patterns: `` cols`,`1,2 ``

**ColumnSpec**: Parsed into `[]colEntry` where each entry is a single index or range. `Resolve(numCols)` returns ordered 0-based indices, silently skipping out-of-bounds.

**regexSplit**: Custom function that splits a string by a `regexp2.Regexp` pattern, collecting text between matches. Preserves empty strings at boundaries (like most languages' regex split).

### Tests Written
- [x] ColumnSpec: single, multiple, reorder, negative, range, open range, reverse range, mixed
- [x] ColumnSpec: out-of-bounds skipped, duplicates preserved
- [x] ColumnSpec: error on empty, zero index, invalid text
- [x] ColumnsRule: select single/multiple columns, reorder
- [x] ColumnsRule: negative index, range, open range
- [x] ColumnsRule: comma delimiter, custom joiner, comma-space regex
- [x] ColumnsRule: multiple spaces collapsed, out-of-bounds, empty joiner
- [x] regexSplit: whitespace, comma, no match, leading/trailing sep, multi-char
- [x] Parser: whitespace default, comma pattern, with joiner, regex pattern, literal backtick, negative, range
- [x] Parser: errors for missing spec, missing delimiter, invalid spec, invalid pattern
- [x] CLI: whitespace splitting, comma splitting, custom joiner, negative index
- [x] CLI: cols then substitute (chained), cols inside if condition

### Files Created
- `internal/rule/columns_rule.go` — ColumnsRule, ColumnSpec, regexSplit
- `internal/rule/columns_rule_test.go` — Rule and spec tests
- `internal/parser/parse_cols_test.go` — Parser tests

### Files Modified
- `internal/parser/parser.go` — Added `cols` dispatch and `parseCols` function

### Deliverable ✅
```bash
echo "alice 25 engineer" | ged 'cols//3,1'
# Output: engineer alice

echo "a,b,c,d" | ged 'cols/,/1,3/-'
# Output: a-c

echo "a  b  c  d  e" | ged 'cols//-1,1'
# Output: e a

echo "alice,25,engineer" | ged 'cols/,/1,3/ | '
# Output: alice | engineer
```

---

## Phase 12: Extraction Rules ✅ COMPLETE

**Goal**: Implement `t/pattern/`, `r/pattern/`, group capture (`1/pattern/`)

**Go Concepts Learned**:
- **`FindStringMatch`**: Returns a `*Match` object with full match and capture groups
- **`Groups()`**: Returns `[]Group` where index 0 is the full match, 1+ are captures
- **`FindNextMatch`**: Iterates to the next match for global extraction
- **Group validity check**: `group.Length == 0` means the group didn't participate

### Implementation Notes

**Three Extraction Rules**:
- `TakeRule` (`t/pattern/`) — extracts the matching portion. With capture groups, returns the first group. With `g` flag, collects all matches space-separated.
- `RemoveRule` (`r/pattern/`) — removes the matching portion using `Replace` with empty string. Supports `g` flag for removing all matches.
- `GroupRule` (`1/pattern/` through `9/pattern/`) — extracts a specific numbered capture group. Group numbers are 1-based.

**No-match behavior**: All three rules pass the line through unchanged if the pattern doesn't match. This is the safe default — no data loss.

**Parser dispatch**: Single-digit commands (`1`-`9`) are matched after all letter commands. The `r` command doesn't conflict with `reverse` because word commands are checked first.

### Tests Written
- [x] Take extracts first match
- [x] Take no match passes through
- [x] Take extracts first capture group when present
- [x] Take returns full match without capture groups
- [x] Take global collects all matches
- [x] Take global with capture groups
- [x] Take first match only without global
- [x] Take case insensitive
- [x] Remove removes first match
- [x] Remove no match passes through
- [x] Remove global removes all matches
- [x] Remove regex match (e.g., strip comments)
- [x] Remove case insensitive
- [x] Remove only first without global
- [x] Group extracts group 1, group 2
- [x] Group no match passes through
- [x] Group out of range passes through
- [x] Group optional group that didn't participate passes through
- [x] Group case insensitive
- [x] Invalid group number (0, negative) errors
- [x] Parser: t/pattern/, r/pattern/, 1/pattern/ through 9/pattern/
- [x] Parser: flags (g, i), alternate delimiters, literal backtick
- [x] Parser: empty pattern errors
- [x] CLI: take, take global, remove, remove global, group capture
- [x] CLI: take then substitute (chained), remove with trim

### Files Created
- `internal/rule/take_rule.go` — TakeRule with global support
- `internal/rule/remove_rule.go` — RemoveRule (uses Replace with empty string)
- `internal/rule/group_rule.go` — GroupRule for numbered capture groups
- `internal/rule/take_rule_test.go` — Tests
- `internal/rule/remove_rule_test.go` — Tests
- `internal/rule/group_rule_test.go` — Tests
- `internal/parser/parse_extract_test.go` — Parser tests

### Files Modified
- `internal/parser/parser.go` — Added `t`, `r`, `1`-`9` dispatch + `parseTake`, `parseRemove`, `parseGroup`
- `cmd/ged/main_test.go` — CLI integration tests

### Deliverable ✅
```bash
echo "hello world" | ged '1/(\w+) (\w+)/'
# Output: hello

echo "abc 123 def 456" | ged 't/\d+/g'
# Output: 123 456

echo "code  # comment" | ged 'r/\s*#.*/'
# Output: code

echo "name: alice, age: 30" | ged '2/name: (\w+), age: (\d+)/'
# Output: 30
```

---

## Phase 13: Control Flow Rules ✅ MOVED TO PHASE 7b

Implemented early as Phase 7b because control flow rules motivated the LineContext refactor.
See Phase 7b above for full details.

---

## Phase 14: External Commands ✅ COMPLETE

**Goal**: Implement `xargs/command/` and `exec/command/`

**Go Concepts Learned**:
- **`os/exec`**: `exec.Command` creates a command; `cmd.Run()` blocks until done
- **Shell execution**: `exec.Command("sh", "-c", cmd)` for shell features (pipes, redirects)
- **`cmd.Stdin`/`cmd.Stdout`**: Wire `io.Reader`/`io.Writer` to control I/O
- **`*exec.ExitError`**: Non-zero exit codes produce this error type
- **Shell quoting**: Single-quote wrapping with `'\\''` escaping prevents injection

### Implementation Notes

**Two Rule Types**:
- `XargsRule` (LineRule) — For each input line, runs `sh -c "command 'line'"`. The line is shell-quoted to prevent injection. If the command fails, the line passes through unchanged.
- `ExecRule` (DocumentRule) — Pipes the entire document into the command's stdin, captures stdout as the new document. Errors propagate (unlike xargs).

**Shell Quoting**: `shellQuote()` wraps strings in single quotes and escapes embedded single quotes as `'\''`. This prevents `$`, backticks, and other metacharacters from being interpreted.

**Parser Syntax**: Uses delimiter-based parsing like other commands:
- `xargs/echo hello/` — runs `echo hello 'line'` for each line
- `exec/sort -n/` — pipes document through `sort -n`
- `exec|grep foo | wc -l|` — alternate delimiters for commands containing `/`

### Tests Written
- [x] Xargs echo with argument
- [x] Xargs transforms each line
- [x] Xargs command producing multiple lines
- [x] Xargs command producing empty output (deletes line)
- [x] Xargs failing command passes line through
- [x] Xargs handles special characters in line (no injection)
- [x] Shell quoting handles single quotes, `$()`, backticks
- [x] Exec sort command
- [x] Exec cat passes through
- [x] Exec grep filters lines
- [x] Exec awk transforms with line numbers
- [x] Exec wc -l counts lines (with pipe)
- [x] Exec grep no match returns error
- [x] Exec nonexistent command returns error
- [x] Exec head limits output
- [x] Parser: xargs/cmd/, exec/cmd/, alternate delimiters, errors
- [x] CLI: xargs echo end-to-end
- [x] CLI: xargs after substitution (chained)
- [x] CLI: exec sort end-to-end
- [x] CLI: exec grep with pipe
- [x] CLI: exec after substitution
- [x] CLI: xargs inside if condition

### Files Created
- `internal/rule/xargs_rule.go` — XargsRule (LineRule) with shellQuote
- `internal/rule/exec_rule.go` — ExecRule (DocumentRule)
- `internal/rule/xargs_rule_test.go` — Rule tests
- `internal/rule/exec_rule_test.go` — Rule tests
- `internal/parser/parse_exec_test.go` — Parser tests

### Files Modified
- `internal/parser/parser.go` — Added `xargs`/`exec` dispatch + `parseXargs`, `parseExec`
- `cmd/ged/main_test.go` — CLI integration tests

### Deliverable ✅
```bash
echo -e "hello\nworld" | ged 'xargs/echo hi/'
# Output: hi hello\nhi world

echo -e "c\na\nb" | ged 'exec/sort/'
# Output: a\nb\nc

echo -e "hello\nworld\nhello" | ged 'if/hello/' '{' 'xargs/echo matched:/' '}'
# Output: matched: hello\nworld\nmatched: hello

echo -e "3\n1\n2" | ged 's/$/x/' 'exec/sort/'
# Output: 1x\n2x\n3x
```

---

## Phase 15: Diff Output and Colors

**Goal**: Implement `--diff` mode and colored output

**Go Concepts Introduced**:
- ANSI escape codes
- Terminal detection (`os.IsTerminal`)
- Diff algorithms (or use a library)
- Optional dependencies

### Steps

1. **Implement diff generation**
   - Compare original vs transformed
   - Generate unified diff format

2. **Add color support**
   - Detect if stdout is a TTY
   - `--color` and `--no-color` flags
   - Color additions green, deletions red

### Tests to Write
- [ ] Diff shows changes correctly
- [ ] Diff context lines configurable
- [ ] Colors applied when enabled
- [ ] Colors disabled on non-TTY
- [ ] --no-color flag works

### Deliverable
```bash
ged 's/foo/bar' --input=file.txt --diff
# Shows unified diff output
```

---

## Phase 16: More Document Rules

**Goal**: Implement `lines/`, `begin/`, `end/`, `border/`, `count`, `uniq`

**Go Concepts Introduced**:
- String formatting with `fmt.Sprintf`
- Document manipulation patterns

### Steps

1. **Implement remaining document rules**:
   - `LinesRule` - prepend line numbers
   - `BeginRule` - prepend to document
   - `EndRule` - append to document
   - `BorderRule` - both begin and end
   - `CountRule` - output line count
   - `UniqueRule` - remove consecutive duplicates

### Tests to Write
- [ ] Lines numbers correctly
- [ ] Begin prepends to document
- [ ] End appends to document
- [ ] Border does both
- [ ] Count outputs number
- [ ] Uniq removes consecutive duplicates

---

## Phase 17: Advanced Conditionals

**Goal**: Implement `ifany/`, `ifnone/`, `else`

**Go Concepts Introduced**:
- Two-pass processing
- Document-level conditions
- Else clause handling

### Steps

1. **Implement IfAnyCondition**
   - Scan entire document first
   - Apply rules to all lines if any matches

2. **Implement IfNoneCondition**
   - Apply rules only if no lines match

3. **Implement else clause**
   - `if/pattern/ { rules } else { other }`

### Tests to Write
- [ ] IfAny applies to all when one matches
- [ ] IfAny applies to none when no match
- [ ] IfNone applies when no match
- [ ] Else clause works with if
- [ ] Else clause works with between

---

## Phase 18: Split and Insert

**Goal**: Implement `split/pattern/` and `insert/pattern/text/`

**Go Concepts Introduced**:
- Rules that produce multiple outputs
- Insertion patterns

### Steps

1. **Implement SplitRule**
   - Split line on pattern
   - Return multiple lines

2. **Implement InsertRule**
   - Insert new line after matching lines

### Tests to Write
- [ ] Split produces multiple lines
- [ ] Split handles no match
- [ ] Insert adds line after match
- [ ] Insert doesn't affect non-matches

---

## Phase 19: Error Handling and Help

**Goal**: Comprehensive error messages, `--help`, `--explain`

**Go Concepts Introduced**:
- Custom error types
- `errors.Is` and `errors.As`
- Help text generation
- Explanation mode

### Steps

1. **Create structured error types**
   - `ParseError` with position info
   - `RuleError` with rule context

2. **Implement --explain**
   - Print what each rule does in plain English

3. **Implement --help**
   - Generate comprehensive help text

### Tests to Write
- [ ] Parse errors include position
- [ ] Rule errors include context
- [ ] Help text is accurate
- [ ] Explain describes rules correctly

---

## Phase 20: Polish and Performance

**Goal**: Optimize, benchmark, and finalize

**Go Concepts Introduced**:
- Benchmarking with `go test -bench`
- Profiling with `pprof`
- `sync.Pool` for object reuse
- Build tags and cross-compilation

### Steps

1. **Add benchmarks**
   - Benchmark common operations
   - Compare with original Node.js version

2. **Optimize hot paths**
   - Regex compilation caching
   - String allocation reduction

3. **Cross-platform builds**
   - Linux, macOS, Windows
   - Create release binaries

4. **Documentation**
   - README with examples
   - Man page generation

---

## Learning Checkpoints

After each phase, you should be comfortable with:

| Phase | Key Go Concepts |
|-------|-----------------|
| 1 | Packages, interfaces, regexp, basic tests |
| 2 | Multiple returns, nil vs empty slice |
| 3 | Slices, variadic functions |
| 4 | Custom types, closures, strconv |
| 5 | strings package, escape handling |
| 6 | sort package, type assertions, buffering |
| 7 | Recursive parsing, tree structures |
| 7b | iota enums, optional interfaces, shared mutable context |
| 8 | Stateful processing |
| 9 | io interfaces, defer, file handling |
| 10 | String manipulation |
| 11 | Index manipulation |
| 12 | Regex submatches |
| 14 | os/exec, subprocesses |
| 15 | Terminal I/O, ANSI codes |
| 16-18 | Pattern consolidation |
| 19 | Error handling patterns |
| 20 | Benchmarking, optimization |

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/rule/...

# Run with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```
