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
| 7-20 | 🔲 Pending | See details below |

**To continue**: Run `go test ./...` to verify everything works, then start Phase 7.

## Project Structure

```
ged/
├── cmd/
│   └── ged/
│       ├── main.go              # CLI entry point
│       └── main_test.go         # CLI integration tests
├── internal/
│   ├── rule/
│   │   ├── rule.go              # LineRule and DocumentRule interfaces
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
│   │   └── *_test.go            # Tests for each
│   ├── parser/
│   │   ├── parser.go            # Rule parsing
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

**Rule Interface Change**: All rules now receive line number:
```go
type Rule interface {
    Apply(line string, lineNum int) ([]string, error)
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

**Two Interfaces**:
```go
type LineRule interface {
    Apply(line string, lineNum int) ([]string, error)
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

## Phase 7: Conditional Rules

**Goal**: Implement `if/pattern/ { rules }`

**Go Concepts Introduced**:
- Recursive parsing
- Tree structures
- Nested rule execution
- Boolean logic

### Steps

1. **Extend parser for block syntax**
   - Detect `{` and `}` tokens
   - Parse nested rules recursively

2. **Create ConditionalRule wrapper**
   ```go
   type ConditionalRule struct {
       condition *regexp.Regexp
       inverted  bool
       rules     []Rule
   }
   ```

3. **Implement condition types**:
   - `if/pattern/` - apply rules to matching lines
   - `!if/pattern/` - apply rules to non-matching lines

4. **Handle condition chaining**
   - `if/foo/ if/bar/ { ... }` - both must match

### Tests to Write
- [ ] If condition applies rules to matches only
- [ ] Inverted if applies to non-matches
- [ ] Non-matching lines pass through unchanged
- [ ] Nested rules execute in order
- [ ] Chained conditions require all to match

### Deliverable
```bash
echo -e "hello\nworld\nhello" | ged 'if/hello/ { s/o/x }'
# Output: hellx\nworld\nhellx
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

## Phase 9: File I/O

**Goal**: Support `--input=file` and `--write`

**Go Concepts Introduced**:
- `os.Open`, `os.Create`
- `defer` for cleanup
- `io.Reader` and `io.Writer` interfaces
- Error wrapping with `fmt.Errorf`
- File permissions

### Steps

1. **Refactor engine to use io interfaces**
   ```go
   func (e *Engine) Process(r io.Reader, w io.Writer) error
   ```

2. **Implement CLI flags**
   - `--input=file` or positional argument
   - `--write` for in-place editing
   - `--write-to=file` for explicit output

3. **Add backup support**
   - `--write-rename=%.backup` creates backup first

4. **Handle multiple input files**
   - Process each file separately
   - Support `--ls` mode (filenames from stdin)

### Tests to Write
- [ ] Read from file works
- [ ] Write to stdout by default
- [ ] Write in-place with --write
- [ ] Write to different file works
- [ ] Backup before writing works
- [ ] Multiple input files process separately
- [ ] ls mode processes each filename

### Deliverable
```bash
ged 's/foo/bar' --input=test.txt --write
```

---

## Phase 10: Text Modification Rules

**Goal**: Implement `trim`, `prepend`, `append`, `surround`, `quote`, `unquote`

**Go Concepts Introduced**:
- `strings.TrimSpace`, `strings.TrimLeft`, `strings.TrimRight`
- String concatenation vs `strings.Builder`
- Unicode handling

### Steps

1. **Implement rules**:
   - `TrimRule` with left/right/both variants
   - `PrependRule` and `AppendRule`
   - `SurroundRule`
   - `QuoteRule` and `UnquoteRule`

### Tests to Write
- [ ] Trim removes whitespace
- [ ] Trim left/right variants work
- [ ] Prepend adds to start
- [ ] Append adds to end
- [ ] Surround wraps with both
- [ ] Quote handles existing quotes
- [ ] Unquote removes outer quotes only

### Deliverable
```bash
echo "  hello  " | ged 'trim'
# Output: hello
```

---

## Phase 11: Column Operations

**Goal**: Implement `cols//1,3,2` for column selection

**Go Concepts Introduced**:
- `strings.Fields` and `strings.Split`
- Index manipulation
- Negative indexing pattern

### Steps

1. **Parse column specification**
   - Positive indices (1-based)
   - Negative indices (from end)
   - Ranges like `1-3`

2. **Implement ColumnsRule**
   - Split line by delimiter (whitespace default)
   - Select and reorder columns
   - Join with output separator

### Tests to Write
- [ ] Select single column
- [ ] Select multiple columns
- [ ] Reorder columns
- [ ] Negative indices work
- [ ] Custom delimiter works
- [ ] Custom output separator works
- [ ] Out-of-bounds columns handled gracefully

### Deliverable
```bash
echo "a b c d" | ged 'cols//3,1'
# Output: c a
```

---

## Phase 12: Extraction Rules

**Goal**: Implement `t/pattern/`, `r/pattern/`, group capture (`1/pattern/`)

**Go Concepts Introduced**:
- Regex submatches
- `regexp.FindStringSubmatch`
- Slice indexing safety

### Steps

1. **Implement TakeRule**
   - Return only the matching portion
   - Return whole line if no match

2. **Implement RemoveRule**
   - Remove the matching portion
   - Return whole line if no match

3. **Implement GroupMatchRule**
   - Extract specific capture group
   - Handle missing groups

### Tests to Write
- [ ] Take extracts match
- [ ] Take returns line on no match
- [ ] Remove deletes match
- [ ] Group extracts numbered group
- [ ] Invalid group number handled
- [ ] TakePrint and RemovePrint variants work

### Deliverable
```bash
echo "hello world" | ged '1/(\w+) (\w+)/'
# Output: hello
```

---

## Phase 13: Control Flow Rules

**Goal**: Implement `on/`, `off/`, `after/`, `toggle/`

**Go Concepts Introduced**:
- Mutable state across lines
- State machines
- Pointer receivers for methods

### Steps

1. **Add print state to Pipeline**
   - `printOn *bool` - nil means no control active

2. **Implement control rules**:
   - `OnRule` - start printing at match
   - `OffRule` - stop printing at match
   - `AfterRule` - start printing after match
   - `ToggleRule` - flip state at each match

### Tests to Write
- [ ] On starts printing at match
- [ ] Off stops printing at match
- [ ] After starts one line after match
- [ ] Toggle flips state
- [ ] Control rules combine correctly
- [ ] State resets between files

### Deliverable
```bash
echo -e "a\nstart\nb\nc" | ged 'on/start/'
# Output: start\nb\nc
```

---

## Phase 14: External Commands

**Goal**: Implement `xargs/command/` and `exec`

**Go Concepts Introduced**:
- `os/exec` package
- `exec.Command`
- Capturing stdout/stderr
- Process environment

### Steps

1. **Implement XargsExecRule**
   - Execute command with line as argument
   - Capture output as new line(s)

2. **Implement DocumentExecRule**
   - Execute entire document as shell script
   - Return output

3. **Handle errors and timeouts**

### Tests to Write
- [ ] Xargs executes for each line
- [ ] Xargs captures output
- [ ] Exec runs document as script
- [ ] Command failures handled
- [ ] Environment paged correctly

### Deliverable
```bash
echo -e "hello\nworld" | ged 'xargs/echo hi/'
# Output: hi hello\nhi world
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
| 8 | Stateful processing |
| 9 | io interfaces, defer, file handling |
| 10 | String manipulation |
| 11 | Index manipulation |
| 12 | Regex submatches |
| 13 | State machines |
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
