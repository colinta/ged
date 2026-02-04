# Go Migration Plan

This plan breaks down the ged project into incremental phases. Each phase introduces new Go concepts while building working, tested functionality.

You are a professional go developer and are teaching me the basics of Go by writing the 'ged' tool together. Before writing code, you should teach me about the library and concepts that we need for that section. Make sure I understand before we add more code to the project.

## Project Structure

```
ged/
├── cmd/
│   └── ged/
│       └── main.go          # CLI entry point
├── internal/
│   ├── rule/
│   │   ├── rule.go          # Rule interface and base types
│   │   ├── line_rules.go    # Line-based rules
│   │   ├── doc_rules.go     # Document-based rules
│   │   ├── conditional.go   # Conditional wrappers
│   │   └── *_test.go        # Tests for each
│   ├── parser/
│   │   ├── parser.go        # Rule parsing
│   │   ├── delimiter.go     # Delimiter handling
│   │   ├── linerange.go     # Line number parsing
│   │   └── *_test.go
│   ├── engine/
│   │   ├── engine.go        # Processing pipeline
│   │   ├── stream.go        # Line streaming
│   │   └── *_test.go
│   └── cli/
│       ├── cli.go           # Argument parsing
│       └── cli_test.go
├── go.mod
├── go.sum
└── Makefile
```

---

## Phase 1: Hello Go - Basic Substitution

**Goal**: Get a working `ged 's/foo/bar'` that reads stdin and writes stdout.

**Go Concepts Introduced**:
- Package structure and `go mod init`
- Basic types: strings, errors
- `fmt` and `os` packages
- `bufio.Scanner` for line reading
- `regexp` package
- Writing and running tests with `go test`

### Steps

1. **Initialize the project**
   ```bash
   mkdir ged && cd ged
   go mod init github.com/colinta/ged
   ```

2. **Create the Rule interface** (`internal/rule/rule.go`)
   - Define `Rule` interface with `Apply(line string) ([]string, error)`
   - Define `LineRule` interface (embeds `Rule`)
   - Create `SubstitutionRule` struct with `pattern *regexp.Regexp` and `replacement string`

3. **Write tests first** (`internal/rule/line_rules_test.go`)
   ```go
   func TestSubstitutionRule(t *testing.T) {
       rule := NewSubstitutionRule("world", "earth")
       result, _ := rule.Apply("hello world")
       if result[0] != "hello earth" {
           t.Errorf("got %q, want %q", result[0], "hello earth")
       }
   }
   ```

4. **Implement SubstitutionRule**
   - Constructor: `NewSubstitutionRule(pattern, replace string) (*SubstitutionRule, error)`
   - Method: `Apply(line string) ([]string, error)`

5. **Create minimal CLI** (`cmd/ged/main.go`)
   - Parse first argument as `s/pattern/replace/`
   - Read stdin line by line
   - Apply rule, print result

6. **Add GlobalSubstitutionRule** (replaces all matches)

### Tests to Write
- [ ] SubstitutionRule replaces first match only
- [ ] SubstitutionRule handles no match (returns original)
- [ ] SubstitutionRule handles regex patterns
- [ ] GlobalSubstitutionRule replaces all matches
- [ ] Case-insensitive flag works

### Deliverable
```bash
echo "hello world" | ged 's/world/earth'
# Output: hello earth
```

---

## Phase 2: Filtering Rules

**Goal**: Implement `p/pattern/` (print matching) and `d/pattern/` (delete matching).

**Go Concepts Introduced**:
- Multiple return values
- Empty slice vs nil semantics
- Table-driven tests
- Error handling patterns

### Steps

1. **Extend the Rule interface**
   - Returning empty slice `[]string{}` means "delete this line"
   - Returning `nil` means "keep original unchanged"

2. **Implement PrintLineRule** (`p/pattern/`)
   - Match: return `[]string{line}`
   - No match: return `[]string{}`

3. **Implement DeleteLineRule** (`d/pattern/`)
   - Match: return `[]string{}`
   - No match: return `[]string{line}`

4. **Refactor parsing**
   - Extract delimiter parsing to `internal/parser/delimiter.go`
   - Support `/`, `|`, `=` delimiters

### Tests to Write
- [ ] PrintLineRule keeps matching lines
- [ ] PrintLineRule removes non-matching lines
- [ ] DeleteLineRule removes matching lines
- [ ] DeleteLineRule keeps non-matching lines
- [ ] Regex patterns work in both rules
- [ ] Different delimiters parse correctly

### Deliverable
```bash
echo -e "foo\nbar\nfoo" | ged 'p/foo/'
# Output: foo\nfoo

echo -e "foo\nbar\nfoo" | ged 'd/foo/'
# Output: bar
```

---

## Phase 3: Rule Chaining

**Goal**: Support multiple rules: `ged 'p/foo/' 's/o/x/'`

**Go Concepts Introduced**:
- Slices and iteration
- Variadic functions
- Method chaining patterns
- Interface composition

### Steps

1. **Create Pipeline type** (`internal/engine/engine.go`)
   ```go
   type Pipeline struct {
       rules []rule.Rule
   }

   func (p *Pipeline) Process(line string) []string
   ```

2. **Handle rule output propagation**
   - A rule can output 0, 1, or many lines
   - Each output line feeds into the next rule
   - Empty output stops processing for that line

3. **Update CLI to accept multiple rules**

### Tests to Write
- [ ] Two rules chain correctly
- [ ] Filter then substitute works
- [ ] Substitute then filter works
- [ ] Empty output stops the chain
- [ ] Multiple output lines all get proceged

### Deliverable
```bash
echo -e "hello\nworld\nhello" | ged 'p/hello/' 's/o/x/'
# Output: hellx\nhellx
```

---

## Phase 4: Line Numbers

**Goal**: Support line number operations: `p:1-5`, `s:2:replacement`

**Go Concepts Introduced**:
- Custom types and methods
- Parsing with `strconv`
- Closures (for line range test functions)
- State in structs

### Steps

1. **Create LineRange type** (`internal/parser/linerange.go`)
   ```go
   type LineRange interface {
       Contains(lineNum int) bool
   }
   ```

2. **Implement range types**:
   - `SingleLine` - matches one line
   - `Range` - matches start to end
   - `OpenRange` - matches from N or to N
   - `Modulo` - matches every Nth line
   - `CompositeRange` - combines multiple ranges

3. **Add line number tracking to Pipeline**
   - Track current line number
   - Pass to rules that need it

4. **Implement line-number-aware rules**
   - `PrintLineNumRule` - print by line number
   - `SubstituteLineRule` - replace entire line by number

### Tests to Write
- [ ] Single line number matches correctly
- [ ] Range `2-4` matches lines 2, 3, 4
- [ ] Open range `5-` matches 5 and beyond
- [ ] Open range `-5` matches 1 through 5
- [ ] Modulo `%2` matches even lines
- [ ] Modulo with offset `%2-1` works
- [ ] Comma-separated ranges work
- [ ] Line substitution replaces entire line

### Deliverable
```bash
echo -e "1\n2\n3\n4\n5" | ged 'p:2-4'
# Output: 2\n3\n4
```

---

## Phase 5: Literal String Matching

**Goal**: Support backtick/quote delimiters for literal matching

**Go Concepts Introduced**:
- `strings` package functions
- `regexp.QuoteMeta`
- Delimiter type system

### Steps

1. **Extend delimiter parser**
   - Track delimiter type (regex vs literal)
   - Return metadata with parsed result

2. **Modify rule constructors**
   - Accept flag for literal vs regex
   - Use `strings.Replace` for literal, `regexp` for regex

3. **Implement escape sequence handling**
   - `\n` → newline
   - `\t` → tab
   - `\\` → backslash

### Tests to Write
- [ ] Backtick treats `.` as literal dot
- [ ] Single quote treats `*` as literal asterisk
- [ ] Escape sequences expand correctly
- [ ] Mixed literal and regex rules work together

### Deliverable
```bash
echo "foo.bar" | ged 's`foo.bar`baz'
# Output: baz  (literal match, not regex)
```

---

## Phase 6: Document Rules

**Goal**: Implement `sort`, `reverse`, `join`

**Go Concepts Introduced**:
- `sort` package
- Slices manipulation
- `strings.Join`
- Interface type assertions
- Buffering strategies

### Steps

1. **Define DocumentRule interface**
   ```go
   type DocumentRule interface {
       ApplyDocument(lines []string) ([]string, error)
   }
   ```

2. **Modify Pipeline to detect document rules**
   - If any document rule exists, buffer all line output
   - Apply document rules to buffer

3. **Implement document rules**:
   - `SortRule` - alphabetic sort
   - `SortNumericRule` - numeric sort
   - `ReverseRule` - reverse order
   - `JoinRule` - join with separator

### Tests to Write
- [ ] Sort orders alphabetically
- [ ] SortNumeric handles numbers correctly
- [ ] SortNumeric handles non-numeric lines
- [ ] Reverse reverses line order
- [ ] Join combines lines with separator
- [ ] Line rules then document rules work
- [ ] Document rules then line rules work

### Deliverable
```bash
echo -e "c\na\nb" | ged 'sort'
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
