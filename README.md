# ged

A streaming text editor for pipelines. Like `sed`, but with modern regex, intuitive syntax, and composable rules.

## Install

```bash
go install github.com/colinta/ged/cmd/ged@latest
```

## Usage

```
ged — a streaming text editor for pipelines

Usage: ged [flags] <rule> [rule...]

Rules:
  s/pattern/replacement/[flags]  Substitute (replace) matches
  s:linerange:replacement        Substitute by line number
  p/pattern/[flags][/options]    Print (keep) matching lines
  p:linerange                    Print by line number
  d/pattern/[flags][/options]    Delete matching lines
  d:linerange                    Delete by line number
  t/pattern/[flags]              Take (extract) matching text
  r/pattern/[flags]              Remove matching text from lines
  N/pattern/[flags]              Extract capture group N (1-9)

Text rules:
  trim, triml, trimr             Trim whitespace
  upper, lower                   Change case
  prepend/text/                  Add text before each line
  append/text/                   Add text after each line
  surround/before/after/         Wrap each line
  cols/pattern/spec[/joiner]     Select/reorder columns
  split/pattern/                 Split lines on pattern
  insert/pattern/text/           Insert text after matching lines

Document rules:
  sort                           Sort lines alphabetically
  reverse                        Reverse line order
  join/separator/                Join lines into one
  lines                          Prepend line numbers
  count                          Output line count
  uniq                           Remove consecutive duplicates
  begin/text/                    Prepend text to document
  end/text/                      Append text to document
  border/text/                   Add text to both ends
  exec/command/                  Pipe document through command
  xargs/command/                 Run command per line

Conditionals:
  if/pattern/ { rules }          Apply rules to matching lines
  !if/pattern/ { rules }         Apply rules to non-matching lines
  between/start/end/ { rules }   Apply rules between patterns
  ifany/pattern/ { rules }       Apply to all if any line matches
  ifnone/pattern/ { rules }      Apply to all if no line matches
  ... else { rules }             Else clause for any conditional

Control flow:
  on/pattern/                    Start printing at match
  off/pattern/                   Stop printing at match
  after/pattern/                 Start printing after match
  toggle/pattern/                Toggle printing at each match

Flags (in trailing delimiter position):
  g   Global (all matches, not just first)
  i   Case-insensitive matching

Options (for p/ and d/ rules):
  context=N                      Include N lines before and after
  before=N                       Include N lines before match
  after=N                        Include N lines after match

Delimiters:
  Any punctuation character works: s/foo/bar/ s|foo|bar| s#foo#bar#
  Quote delimiters (backtick, ', ") enable literal matching:
    s`foo.bar`baz`    matches "foo.bar" literally

CLI flags:
  --input=FILE    Read from file (repeatable)
  --write         Write back to input file(s)
  --write-to=FILE Write output to specific file
  --diff          Show diff instead of output
  --color         Force colored output
  --no-color      Force plain output
  --explain       Describe what rules do
  --help, -h      Show this help

Examples:
  echo "hello world" | ged 's/world/earth/'
  ged --input=log.txt 'p/ERROR/context=2'
  ged 's/old/new/g' sort uniq --input=data.txt --write
  ged 'if/TODO/' '{' upper '}' --input=notes.txt --diff
```

## Examples

```bash
# Simple substitution
echo "hello world" | ged 's/world/earth/'

# Global replace
echo "aaa" | ged 's/a/b/g'

# Filter lines (like grep)
cat log.txt | ged 'p/ERROR/'

# Filter with context (like grep -C)
cat log.txt | ged 'p/ERROR/context=2'

# Delete matching lines
cat data.txt | ged 'd/^#/'

# Chain multiple rules
echo -e "c\na\nb" | ged 's/$/!/' sort

# Conditional rules
echo -e "hello\nworld" | ged 'if/hello/' '{' upper '}'

# File editing
ged 's/old/new/g' --input=file.txt --write

# Show diff before editing
ged 's/old/new/g' --input=file.txt --diff --color

# Explain what rules do
ged --explain 'p/error/i/context=2' 's/ERROR/WARNING/' upper
```

## License

MIT
