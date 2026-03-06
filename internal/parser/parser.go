// Package parser handles parsing of ged command syntax.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/colinta/ged/internal/rule"
	"github.com/dlclark/regexp2"
)

// ParseRule parses a rule string and returns the appropriate Rule.
// It handles delimiter detection and dispatches to command-specific parsers.
// Returns either a rule.LineRule or rule.DocumentRule (as any).
func ParseRule(input string) (any, error) {
	// Word commands must be checked first — "sort" starts with 's',
	// which would otherwise match the substitution command.
	if input == "sort" {
		return rule.NewSortRule(), nil
	}
	if input == "reverse" {
		return rule.NewReverseRule(), nil
	}
	if input == "lines" {
		return rule.NewLinesRule(), nil
	}
	if input == "count" {
		return rule.NewCountRule(), nil
	}
	if input == "uniq" {
		return rule.NewUniqRule(), nil
	}
	if strings.HasPrefix(input, "begin") {
		return parseDocTextCommand(input, "begin")
	}
	if strings.HasPrefix(input, "end") {
		return parseDocTextCommand(input, "end")
	}
	if strings.HasPrefix(input, "border") {
		return parseDocTextCommand(input, "border")
	}
	if strings.HasPrefix(input, "join") {
		return parseJoin(input)
	}
	if strings.HasPrefix(input, "!between") || strings.HasPrefix(input, "between") {
		return parseBetween(input)
	}
	if strings.HasPrefix(input, "!ifany") || strings.HasPrefix(input, "ifany") {
		return parseIfAny(input)
	}
	if strings.HasPrefix(input, "!ifnone") || strings.HasPrefix(input, "ifnone") {
		return parseIfNone(input)
	}
	if strings.HasPrefix(input, "!if") || strings.HasPrefix(input, "if") {
		return parseIf(input)
	}
	if strings.HasPrefix(input, "on") {
		return parseControl(input, "on")
	}
	if strings.HasPrefix(input, "off") {
		return parseControl(input, "off")
	}
	if strings.HasPrefix(input, "after") {
		return parseControl(input, "after")
	}
	if strings.HasPrefix(input, "toggle") {
		return parseControl(input, "toggle")
	}

	// Text modification word commands
	if input == "trim" {
		return rule.NewTrimRule(), nil
	}
	if input == "triml" {
		return rule.NewTrimLeftRule(), nil
	}
	if input == "trimr" {
		return rule.NewTrimRightRule(), nil
	}
	if input == "upper" {
		return rule.NewUpperRule(), nil
	}
	if input == "lower" {
		return rule.NewLowerRule(), nil
	}
	if strings.HasPrefix(input, "prepend") {
		return parseTextCommand(input, "prepend", 1)
	}
	if strings.HasPrefix(input, "append") {
		return parseTextCommand(input, "append", 1)
	}
	if strings.HasPrefix(input, "surround") {
		return parseTextCommand(input, "surround", 2)
	}
	if strings.HasPrefix(input, "cols") {
		return parseCols(input)
	}
	if strings.HasPrefix(input, "split") {
		return parseSplit(input)
	}
	if strings.HasPrefix(input, "insert") {
		return parseInsert(input)
	}
	if strings.HasPrefix(input, "xargs") {
		return parseXargs(input)
	}
	if strings.HasPrefix(input, "exec") {
		return parseExec(input)
	}

	if len(input) < 2 {
		return nil, fmt.Errorf("invalid rule: too short")
	}

	command := input[0]
	delimiter := input[1]
	rest := input[2:]

	// Split by delimiter, respecting backslash escapes
	parts, err := splitByDelimiter(rest, delimiter)
	if err != nil {
		return nil, err
	}

	// Quote delimiters mean literal matching — escape regex metacharacters
	if (delimiter == '`' || delimiter == '\'' || delimiter == '"') && len(parts) > 0 {
		parts[0] = regexp2.Escape(parts[0])
	}

	if command == 'p' && delimiter == ':' {
		return parsePrintLineNum(parts)
	} else if command == 'p' {
		return parsePrint(parts)
	} else if command == 'd' && delimiter == ':' {
		return parseDeleteLineNum(parts)
	} else if command == 'd' {
		return parseDelete(parts)
	} else if command == 's' && delimiter == ':' {
		return parseSubstitutionLineNum(parts)
	} else if command == 's' {
		return parseSubstitution(parts)
	} else if command == 't' {
		return parseTake(parts)
	} else if command == 'r' {
		return parseRemove(parts)
	} else if command >= '1' && command <= '9' {
		return parseGroup(parts, int(command-'0'))
	} else {
		return nil, fmt.Errorf("unknown command: %c", command)
	}
}

// parseFlags reads a flags string and returns the corresponding RuleOptions.
// Supported flags:
//
//	g — global replacement (SubstitutionRule only)
//	i — case-insensitive matching
func parseFlags(flags string) []rule.RuleOption {
	var opts []rule.RuleOption
	if strings.Contains(flags, "g") {
		opts = append(opts, rule.WithGlobal())
	}
	if strings.Contains(flags, "i") {
		opts = append(opts, rule.WithIgnoreCase())
	}
	return opts
}

// flagsFromParts extracts flags from the trailing element of a parts slice.
// For commands like p/pat/ and d/pat/, flags are in parts[1].
// For substitution s/pat/repl/flags, flags are in parts[2].
// Returns the options parsed from the given index, or nil if index is out of range.
func flagsFromParts(parts []string, flagIndex int) []rule.RuleOption {
	if flagIndex < len(parts) {
		return parseFlags(parts[flagIndex])
	}
	return nil
}

// contextOptions holds parsed before/after context values.
type contextOptions struct {
	before int
	after  int
	has    bool // true if any context option was found
}

// parseContextOptions scans parts starting at startIndex for context options.
// Recognized: context=N, before=N, after=N. Parts without = are skipped (flags).
// Returns the parsed options and any error.
func parseContextOptions(parts []string, startIndex int) (contextOptions, error) {
	var ctx contextOptions
	for i := startIndex; i < len(parts); i++ {
		part := parts[i]
		if !strings.Contains(part, "=") {
			continue // skip flag strings
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, val := kv[0], kv[1]
		n, err := strconv.Atoi(val)
		if err != nil {
			return ctx, fmt.Errorf("invalid value for %s: %q", key, val)
		}
		if n < 0 {
			return ctx, fmt.Errorf("%s must be non-negative, got %d", key, n)
		}
		switch key {
		case "context":
			ctx.before = n
			ctx.after = n
			ctx.has = true
		case "before":
			ctx.before = n
			ctx.has = true
		case "after":
			ctx.after = n
			ctx.has = true
		default:
			return ctx, fmt.Errorf("unknown option: %s", key)
		}
	}
	return ctx, nil
}

// flagsAndOptionsFromParts extracts both flags and context options from trailing parts.
// Flags are parts without =, context options are parts with =.
func flagsAndOptionsFromParts(parts []string, startIndex int) ([]rule.RuleOption, contextOptions, error) {
	var opts []rule.RuleOption
	for i := startIndex; i < len(parts); i++ {
		if !strings.Contains(parts[i], "=") {
			opts = append(opts, parseFlags(parts[i])...)
		}
	}
	ctxOpts, err := parseContextOptions(parts, startIndex)
	return opts, ctxOpts, err
}

// parseTextCommand handles word commands that take text arguments via delimiters.
// name is the command name, requiredParts is how many delimiter-separated parts are needed.
func parseTextCommand(input string, name string, requiredParts int) (rule.LineRule, error) {
	rest := input[len(name):]
	if len(rest) == 0 {
		return nil, fmt.Errorf("%s requires a delimiter and argument(s)", name)
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < requiredParts || parts[0] == "" {
		return nil, fmt.Errorf("%s requires %d argument(s)", name, requiredParts)
	}

	switch name {
	case "prepend":
		return rule.NewPrependRule(parts[0]), nil
	case "append":
		return rule.NewAppendRule(parts[0]), nil
	case "surround":
		if len(parts) < 2 || parts[1] == "" {
			return nil, fmt.Errorf("surround requires before and after text")
		}
		return rule.NewSurroundRule(parts[0], parts[1]), nil
	default:
		return nil, fmt.Errorf("unknown text command: %s", name)
	}
}

// parseCols handles "cols/pattern/spec" and "cols/pattern/spec/joiner" syntax.
// An empty pattern defaults to \s+ (whitespace splitting).
func parseCols(input string) (rule.LineRule, error) {
	rest := input[4:] // skip "cols"
	if len(rest) == 0 {
		return nil, fmt.Errorf("cols requires a delimiter, pattern, and column spec")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 2 {
		return nil, fmt.Errorf("cols requires a pattern and column spec")
	}

	// Pattern: empty means \s+
	patternStr := parts[0]
	if patternStr == "" {
		patternStr = `\s+`
	} else if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		patternStr = regexp2.Escape(patternStr)
	}

	pattern, err := rule.CompilePattern(patternStr)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in cols: %w", err)
	}

	// Column spec
	spec, err := rule.ParseColumnSpec(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid column spec in cols: %w", err)
	}

	// Joiner: default to " "
	joiner := " "
	if len(parts) >= 3 {
		joiner = parts[2]
	}

	return rule.NewColumnsRule(pattern, spec, joiner), nil
}

// parseJoin handles "join" (bare) and "join/sep/" syntax.
func parseJoin(input string) (*rule.JoinRule, error) {
	if input == "join" {
		return rule.NewJoinRule(""), nil
	}

	if len(input) < 5 {
		return nil, fmt.Errorf("invalid join syntax: %q", input)
	}

	delimiter := input[4]
	rest := input[5:]

	parts, err := splitByDelimiter(rest, delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 {
		return nil, fmt.Errorf("join requires a separator")
	}

	return rule.NewJoinRule(parts[0]), nil
}

// splitByDelimiter splits a string by delimiter, respecting backslash escapes.
// The delimiter at the end is required (trailing part can be empty for flags).
// Returns the parts with escape sequences processed.
func splitByDelimiter(input string, delimiter byte) ([]string, error) {
	var parts []string
	var current strings.Builder

	i := 0
	for i < len(input) {
		ch := input[i]

		if ch == '\\' && i+1 < len(input) {
			next := input[i+1]
			if next == delimiter {
				// Escaped delimiter - write the delimiter itself
				current.WriteByte(delimiter)
				i += 2
				continue
			} else if next == '\\' {
				// Escaped backslash - write single backslash
				current.WriteByte('\\')
				i += 2
				continue
			} else if next == 'n' {
				current.WriteByte('\n')
				i += 2
				continue
			} else if next == 't' {
				current.WriteByte('\t')
				i += 2
				continue
			}
			// Not an escape sequence we handle - write backslash and continue
			current.WriteByte(ch)
			i++
			continue
		}

		if ch == delimiter {
			parts = append(parts, current.String())
			current.Reset()
			i++
			continue
		}

		current.WriteByte(ch)
		i++
	}

	// The last part (after final delimiter) contains flags
	// It's okay if there's content - that's the flags portion
	parts = append(parts, current.String())

	return parts, nil
}

// parseSubstitution creates a SubstitutionRule from parsed parts.
// Expected parts: [pattern, replacement, flags]
// The trailing delimiter is required, so we need at least 3 parts.
func parseSubstitution(parts []string) (*rule.SubstitutionRule, error) {
	if len(parts) < 2 {
		return nil, fmt.Errorf("substitution requires pattern and replacement with trailing delimiter")
	}

	pattern := parts[0]
	replace := parts[1]
	opts := flagsFromParts(parts, 2)

	return rule.NewSubstitutionRule(pattern, replace, opts...)
}

// parseSubstitutionLineNum creates a SubLineNumRule for line number replacement.
// Expected parts: [lineRange, replacement]
func parseSubstitutionLineNum(parts []string) (rule.LineRule, error) {
	if len(parts) < 2 {
		return nil, fmt.Errorf("substitution requires a line range and replacement")
	}

	lineRange, err := rule.ParseLineRange(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid line range: %w", err)
	}
	return rule.NewSubLineNumRule(lineRange, parts[1]), nil
}

// parsePrint creates a PrintLineRule for pattern matching.
func parsePrint(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("print requires a pattern")
	}

	opts, ctxOpts, err := flagsAndOptionsFromParts(parts, 1)
	if err != nil {
		return nil, fmt.Errorf("print: %w", err)
	}
	if ctxOpts.has {
		return rule.NewPrintContextRule(parts[0], ctxOpts.before, ctxOpts.after, opts...)
	}
	return rule.NewPrintLineRule(parts[0], opts...)
}

// parsePrintLineNum creates a PrintLineNumRule for line number filtering.
func parsePrintLineNum(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("print requires a line range")
	}

	lineRange, err := rule.ParseLineRange(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid line range: %w", err)
	}
	return rule.NewPrintLineNumRule(lineRange), nil
}

// parseDelete creates a DeleteLineRule for pattern matching.
func parseDelete(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("delete requires a pattern")
	}

	opts, ctxOpts, err := flagsAndOptionsFromParts(parts, 1)
	if err != nil {
		return nil, fmt.Errorf("delete: %w", err)
	}
	if ctxOpts.has {
		return rule.NewDeleteContextRule(parts[0], ctxOpts.before, ctxOpts.after, opts...)
	}
	return rule.NewDeleteLineRule(parts[0], opts...)
}

// parseDeleteLineNum creates a DeleteLineNumRule for line number filtering.
func parseDeleteLineNum(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("delete requires a line range")
	}

	lineRange, err := rule.ParseLineRange(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid line range: %w", err)
	}
	return rule.NewDeleteLineNumRule(lineRange), nil
}

// parseTake creates a TakeRule for extracting matching text.
// Syntax: t/pattern/, t/pattern/flags, t/pattern/flags/joiner
// The joiner is used in global mode to join multiple matches (default: space).
func parseTake(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("take requires a pattern")
	}
	opts := flagsFromParts(parts, 1)
	joiner := " " // default
	if len(parts) > 2 {
		joiner = parts[2]
	}
	return rule.NewTakeRule(parts[0], joiner, opts...)
}

// parseRemove creates a RemoveRule for removing matching text.
func parseRemove(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("remove requires a pattern")
	}
	opts := flagsFromParts(parts, 1)
	return rule.NewRemoveRule(parts[0], opts...)
}

// parseGroup creates a GroupRule for extracting a numbered capture group.
func parseGroup(parts []string, groupNum int) (rule.LineRule, error) {
	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("group requires a pattern")
	}
	opts := flagsFromParts(parts, 1)
	return rule.NewGroupRule(parts[0], groupNum, opts...)
}

// parseDocTextCommand handles document-level commands that take text arguments: begin, end, border.
func parseDocTextCommand(input string, name string) (rule.DocumentRule, error) {
	rest := input[len(name):]
	if len(rest) == 0 {
		return nil, fmt.Errorf("%s requires a delimiter and text", name)
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("%s requires text", name)
	}

	switch name {
	case "begin":
		return rule.NewBeginRule(parts[0]), nil
	case "end":
		return rule.NewEndRule(parts[0]), nil
	case "border":
		return rule.NewBorderRule(parts[0]), nil
	default:
		return nil, fmt.Errorf("unknown document text command: %s", name)
	}
}

// parseXargs parses "xargs/command/" and returns a LineRule that runs the command for each line.
func parseXargs(input string) (rule.LineRule, error) {
	rest := input[5:] // skip "xargs"
	if len(rest) == 0 {
		return nil, fmt.Errorf("xargs requires a command")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("xargs requires a command")
	}

	return rule.NewXargsRule(parts[0]), nil
}

// parseExec parses "exec/command/" and returns a DocumentRule that pipes document through the command.
func parseExec(input string) (rule.DocumentRule, error) {
	rest := input[4:] // skip "exec"
	if len(rest) == 0 {
		return nil, fmt.Errorf("exec requires a command")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("exec requires a command")
	}

	return rule.NewExecRule(parts[0]), nil
}

// parseSplit parses "split/pattern/" and returns a SplitRule.
func parseSplit(input string) (rule.LineRule, error) {
	rest := input[5:] // skip "split"
	if len(rest) == 0 {
		return nil, fmt.Errorf("split requires a pattern")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("split requires a pattern")
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	opts := flagsFromParts(parts, 1)
	return rule.NewSplitRule(pattern, opts...)
}

// parseInsert parses "insert/pattern/text/" and returns an InsertRule.
func parseInsert(input string) (rule.LineRule, error) {
	rest := input[6:] // skip "insert"
	if len(rest) == 0 {
		return nil, fmt.Errorf("insert requires a pattern and text")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 2 || parts[0] == "" {
		return nil, fmt.Errorf("insert requires a pattern and text")
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	text := parts[1]
	opts := flagsFromParts(parts, 2)
	return rule.NewInsertRule(pattern, text, opts...)
}

// parseControl parses "name/pattern/" for control rules (on, off, after, toggle).
func parseControl(input string, name string) (rule.LineRule, error) {
	rest := input[len(name):]
	if len(rest) == 0 {
		return nil, fmt.Errorf("%s requires a pattern", name)
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("%s requires a pattern", name)
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	opts := flagsFromParts(parts, 1)

	switch name {
	case "on":
		return rule.NewOnRule(pattern, opts...)
	case "off":
		return rule.NewOffRule(pattern, opts...)
	case "after":
		return rule.NewAfterRule(pattern, opts...)
	case "toggle":
		return rule.NewToggleRule(pattern, opts...)
	default:
		return nil, fmt.Errorf("unknown control command: %s", name)
	}
}

// condition is a parser-internal type representing a parsed if/!if condition.
// It's not a rule — it gets converted into a ConditionalRule once the inner
// rules are collected from the { } block.
type condition struct {
	pattern  *regexp2.Regexp
	inverted bool
}

// parseIf parses "if/pattern/" or "!if/pattern/" and returns a condition.
func parseIf(input string) (*condition, error) {
	inverted := false
	rest := input

	if strings.HasPrefix(rest, "!if") {
		inverted = true
		rest = rest[3:]
	} else {
		rest = rest[2:]
	}

	if len(rest) == 0 {
		return nil, fmt.Errorf("missing pattern in if condition")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("missing pattern in if condition")
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	opts := flagsFromParts(parts, 1)
	compiled, err := rule.CompilePattern(pattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in if condition: %w", err)
	}

	return &condition{
		pattern:  compiled,
		inverted: inverted,
	}, nil
}

// ifAnyCondition is a parser-internal type for ifany/!ifany conditions.
type ifAnyCondition struct {
	pattern  *regexp2.Regexp
	inverted bool
}

// parseIfAny parses "ifany/pattern/" or "!ifany/pattern/" and returns an ifAnyCondition.
func parseIfAny(input string) (*ifAnyCondition, error) {
	inverted := false
	rest := input

	if strings.HasPrefix(rest, "!ifany") {
		inverted = true
		rest = rest[6:]
	} else {
		rest = rest[5:]
	}

	if len(rest) == 0 {
		return nil, fmt.Errorf("missing pattern in ifany condition")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("missing pattern in ifany condition")
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	opts := flagsFromParts(parts, 1)
	compiled, err := rule.CompilePattern(pattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in ifany condition: %w", err)
	}

	return &ifAnyCondition{
		pattern:  compiled,
		inverted: inverted,
	}, nil
}

// ifNoneCondition is a parser-internal type for ifnone/!ifnone conditions.
type ifNoneCondition struct {
	pattern  *regexp2.Regexp
	inverted bool
}

// parseIfNone parses "ifnone/pattern/" or "!ifnone/pattern/" and returns an ifNoneCondition.
func parseIfNone(input string) (*ifNoneCondition, error) {
	inverted := false
	rest := input

	if strings.HasPrefix(rest, "!ifnone") {
		inverted = true
		rest = rest[7:]
	} else {
		rest = rest[6:]
	}

	if len(rest) == 0 {
		return nil, fmt.Errorf("missing pattern in ifnone condition")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("missing pattern in ifnone condition")
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	opts := flagsFromParts(parts, 1)
	compiled, err := rule.CompilePattern(pattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in ifnone condition: %w", err)
	}

	return &ifNoneCondition{
		pattern:  compiled,
		inverted: inverted,
	}, nil
}

// betweenCondition is a parser-internal type representing a parsed between condition.
// Like condition, it gets assembled with inner rules from { } blocks in parseArgs.
type betweenCondition struct {
	startPattern *regexp2.Regexp
	endPattern   *regexp2.Regexp
	inverted     bool
}

// parseBetween parses "between/start/end/" or "!between/start/end/" and returns a betweenCondition.
func parseBetween(input string) (*betweenCondition, error) {
	inverted := false
	rest := input

	if strings.HasPrefix(rest, "!between") {
		inverted = true
		rest = rest[8:]
	} else {
		rest = rest[7:]
	}

	if len(rest) == 0 {
		return nil, fmt.Errorf("between requires start and end patterns")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("between requires start and end patterns")
	}

	startPattern := parts[0]
	endPattern := parts[1]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		startPattern = regexp2.Escape(startPattern)
		endPattern = regexp2.Escape(endPattern)
	}

	opts := flagsFromParts(parts, 2)
	startCompiled, err := rule.CompilePattern(startPattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid start pattern in between: %w", err)
	}

	endCompiled, err := rule.CompilePattern(endPattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid end pattern in between: %w", err)
	}

	return &betweenCondition{
		startPattern: startCompiled,
		endPattern:   endCompiled,
		inverted:     inverted,
	}, nil
}
