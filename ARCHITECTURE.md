# ged Architecture

**ged** is a modern text transformation tool that processes input through a pipeline of composable rules. Unlike traditional sed, it uses standard regex syntax and provides a rich set of operations for line-by-line and document-level transformations.

## Core Concepts

### Rules

Rules are the fundamental unit of operation. Each rule takes input text and produces transformed output. Rules are classified into two categories based on how they process input:

**Line Rules** process input one line at a time:
- Can stream input without buffering the entire document
- Enable real-time processing of large files or infinite streams
- Examples: substitution, filtering, control flow (on/off/after/toggle)

**Document Rules** require the entire document:
- Buffer all input before processing begins
- Enable operations that span multiple lines
- Examples: sorting, reversing, joining lines

### Rule Pipeline

Rules are processed in a pipeline where the output of one rule feeds into the next:

```
Input → Rule 1 → Rule 2 → Rule 3 → Output
```

The pipeline automatically handles the boundary between line rules and document rules:
1. Consecutive line rules are grouped into an `ApplyAllRule` (a document rule wrapper)
2. When only line rules exist, input streams line-by-line without buffering
3. When document rules are present, all input is buffered first
4. Each document rule processes the full buffer in sequence

### Delimiters

Rules use delimiters to separate their arguments. The choice of delimiter affects matching behavior:

| Delimiter | Matching Mode | Example |
|-----------|---------------|---------|
| `/`, `\|`, `=`, `#` | Regex pattern | `s/foo.*/bar/` |
| `` ` ``, `'`, `"` | Literal string | `` s`foo`bar `` |
| `:` | Line numbers | `s:1:replacement` |

### Line Number Syntax

Line-based operations support flexible line specification:

- `1` - Single line
- `1-5` - Inclusive range
- `5-` - From line 5 to end
- `-5` - From start to line 5
- `1,3,5-7` - Multiple ranges/lines (composite)

### Conditional Blocks

Rules can be applied conditionally using `if` wrappers:

```
if/pattern/ { rules }       # Apply to matching lines
!if/pattern/ { rules }      # Apply to non-matching lines
```

Conditions can be nested: `if/foo/ { if/bar/ { rules } }`

When all inner rules are line rules, a `ConditionalLineRule` streams line-by-line. When inner rules include document rules, a `ConditionalDocRule` collects matching lines, applies inner rules, and weaves results back into position.

## Implemented Rules

### Substitution Rules
- **s/pattern/replace/** - Replace first match per line
- **s/pattern/replace/g** - Replace all matches per line (global flag)
- **s:linerange:replacement** - Replace entire line content by line number

### Filtering Rules
- **p/pattern/** - Print only matching lines (grep)
- **d/pattern/** - Delete matching lines (inverse grep)
- **p:linerange** - Print lines by number
- **d:linerange** - Delete lines by number

### Control Flow Rules
- **on/pattern/** - Start printing at match (match line included)
- **off/pattern/** - Stop printing at match (match line excluded)
- **after/pattern/** - Start printing after match (match line excluded)
- **toggle/pattern/** - Toggle printing state at each match

### Document Rules
- **sort** - Alphabetic sort
- **reverse** - Reverse line order
- **join/separator/** - Join lines with separator
- **join** - Join lines with empty separator

### Conditional Rules
- **if/pattern/ { rules }** - Apply rules to matching lines
- **!if/pattern/ { rules }** - Apply rules to non-matching lines

## Processing Pipeline

```
┌─────────────────────────────────────────────────────────────┐
│                        Input Stage                          │
├─────────────────────────────────────────────────────────────┤
│  Source: stdin (files planned for Phase 9)                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Rule Processing                         │
├─────────────────────────────────────────────────────────────┤
│  If only LineRules:                                         │
│    Stream line-by-line through Pipeline                     │
│    Check ctx.Printing after each line for output decision   │
│                                                             │
│  If DocumentRules present:                                  │
│    Buffer all input, apply DocumentRules in sequence        │
│    ApplyAllRule wraps consecutive LineRules                  │
│    ApplyAllRule checks ctx.Printing per line                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        Output Stage                         │
├─────────────────────────────────────────────────────────────┤
│  Destination: stdout                                        │
└─────────────────────────────────────────────────────────────┘
```

## State Management

### LineContext

All per-line state flows through `LineContext`, passed to every `LineRule.Apply` call:

```go
type LineContext struct {
    LineNum  int        // 1-indexed line number
    Printing PrintState // controls output inclusion
}
```

### PrintState

An enum with three values:
- `PrintDefault` (0) — no control rule active, print everything
- `PrintOn` — printing is enabled
- `PrintOff` — printing is suppressed

The zero value is `PrintDefault`, so a fresh `LineContext{}` prints everything by default.

### SetupRule Interface

Control flow rules (on/off/after/toggle) implement the optional `SetupRule` interface:

```go
type SetupRule interface {
    Setup(ctx *LineContext)
}
```

`Setup` is called once before processing begins to set the initial `PrintState`. Each rule only sets the initial state if `Printing` is still `PrintDefault` — so the first control rule in the pipeline determines the starting state.

### Rule-Local State

Some rules maintain internal mutable state:
- **AfterRule**: `matched bool` tracks whether the pattern has been seen, delaying the print-on by one line
- **ConditionalDocRule**: Tracks which lines match via `[]bool` to weave results back

State is tied to the rule instance. If the same rule processes multiple documents, state carries over (relevant for future multi-file support).

## Design Principles

1. **Composability**: Rules combine naturally through piping
2. **Predictability**: Standard regex syntax, no surprises
3. **Streaming**: Line rules enable processing of infinite streams
4. **Shared context over rule-local state**: Mutable state that multiple rules need (like print on/off) lives on `LineContext`, not on individual rules
5. **Optional interfaces**: Rules only implement what they need (`SetupRule` is optional)
6. **Small interfaces**: `LineRule` and `DocumentRule` are minimal; extensions use separate interfaces

## Planned Features

See CLAUDE.md for the full phase-by-phase roadmap. Key upcoming features:
- Between conditions (`between/start/end/ { rules }`)
- File I/O (`--input`, `--write`)
- Text modification rules (`trim`, `prepend`, `append`)
- Column operations (`cols`)
- Extraction rules (`t/pattern/`, `r/pattern/`)
- External commands (`xargs`, `exec`)
- Diff output and colors
