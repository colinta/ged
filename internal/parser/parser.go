// Package parser handles parsing of ged command syntax.
package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/colinta/ged/internal/rule"
	"github.com/dlclark/regexp2"
)

// globalOpts holds options applied to all regex-based rules (e.g. --insensitive).
// Set by ParseArgs before parsing begins; reset after parsing completes.
var globalOpts []rule.RuleOption

// commands is the global command registry, populated at init time.
var commands = &registry{}

func init() {
	r := commands

	// -- Exact word commands (no delimiter) --------------------------------
	r.Add(Exact("sort"), Returns(rule.NewSortRule))
	r.Add(Exact("reverse", "rev"), Returns(rule.NewReverseRule))
	r.Add(Exact("lines", "line"), Returns(rule.NewLinesRule))
	r.Add(Exact("count"), Returns(rule.NewCountRule))
	r.Add(Exact("trim"), Returns(rule.NewTrimRule))
	r.Add(Exact("triml"), Returns(rule.NewTrimLeftRule))
	r.Add(Exact("trimr"), Returns(rule.NewTrimRightRule))
	r.Add(Exact("quote"), handleBareQuote)
	r.Add(Exact("unquote"), handleBareUnquote)
	r.Add(Exact("upper"), Returns(rule.NewUpperRule))
	r.Add(Exact("lower"), Returns(rule.NewLowerRule))
	r.Add(Exact("join"), handleBareJoin)

	// -- Prefix commands (name + delimiter + parts) -----------------------

	// uniq with pattern must come before generic prefix "u" if we ever add one
	r.Add(Prefix("uniq", "unique"), handleUniq)

	// Document text commands
	r.Add(Prefix("begin"), handleDocText)
	r.Add(Prefix("end"), handleDocText)
	r.Add(Prefix("border"), handleDocText)

	// Join with separator
	r.Add(Prefix("join"), handleJoin)

	// Conditionals (support ! prefix for negation)
	r.Add(NegPrefix("between"), handleBetween)
	r.Add(NegPrefix("ifany"), handleIfAny)
	r.Add(NegPrefix("ifnone"), handleIfNone)
	r.Add(NegPrefix("if"), handleIf)

	// Control flow
	r.Add(Prefix("on"), handleControl)
	r.Add(Prefix("off"), handleControl)
	r.Add(Prefix("after"), handleControl)
	r.Add(Prefix("toggle"), handleControl)

	// Text modification with delimiter args
	r.Add(Prefix("quote"), handleQuote)
	r.Add(Prefix("unquote"), handleUnquote)
	r.Add(Prefix("prepend"), handlePrepend)
	r.Add(Prefix("append"), handleAppend)
	r.Add(Prefix("surround"), handleSurround)
	r.Add(Prefix("cols"), handleCols)
	r.Add(Prefix("split"), handleSplit)
	r.Add(Prefix("insert"), handleInsert)
	r.Add(Prefix("xargs"), handleXargs)
	r.Add(Prefix("exec"), handleExec)

	// -- Single-character delimiter commands (with word aliases) -----------
	r.Add(Prefix("sub", "substitute"), handleSubstitution)
	r.Add(Prefix("s"), handleSubstitution)
	r.Add(Prefix("print"), handlePrint)
	r.Add(Prefix("p"), handlePrint)
	r.Add(Prefix("d", "del", "delete", "!p", "!print"), handleDelete)
	r.Add(Prefix("takeprint", "tp"), handleTakePrint)
	r.Add(Prefix("printtake", "pt"), handleTakePrint)
	r.Add(Prefix("removeprint", "rp"), handleRemovePrint)
	r.Add(Prefix("printremove", "pr"), handleRemovePrint)
	r.Add(Prefix("take"), handleTake)
	r.Add(Prefix("t"), handleTake)
	r.Add(Prefix("remove"), handleRemove)
	r.Add(Prefix("r"), handleRemove)

	// Digit group extraction: 1/pat/ through 9/pat/
	r.Add(CharRange('1', '9'), handleGroup)
}

// ParseRule parses a rule string and returns the appropriate Rule.
// It handles delimiter detection and dispatches to command-specific handlers.
// Returns either a rule.LineRule, rule.DocumentRule, or a condition type (as any).
func ParseRule(input string) (any, error) {
	return commands.Parse(input)
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func handleBareQuote(_ Command) (any, error) {
	return rule.NewQuoteRule(""), nil
}

func handleQuote(cmd Command) (any, error) {
	if err := cmd.RequireParts(1); err != nil {
		return nil, err
	}
	return rule.NewQuoteRule(cmd.Parts[0]), nil
}

func handleBareUnquote(_ Command) (any, error) {
	return rule.NewUnquoteRule(""), nil
}

func handleUnquote(cmd Command) (any, error) {
	if err := cmd.RequireParts(1); err != nil {
		return nil, err
	}
	return rule.NewUnquoteRule(cmd.Parts[0]), nil
}

func handleBareJoin(_ Command) (any, error) {
	return rule.NewJoinRule(" "), nil
}

func handleJoin(cmd Command) (any, error) {
	if err := cmd.RequireParts(1); err != nil {
		return nil, err
	}
	return rule.NewJoinRule(cmd.Parts[0]), nil
}

func handleDocText(cmd Command) (any, error) {
	if len(cmd.Parts) < 1 || cmd.Parts[0] == "" {
		return nil, fmt.Errorf("%s requires text", cmd.Name)
	}

	switch cmd.Name {
	case "begin":
		return rule.NewBeginRule(cmd.Parts[0]), nil
	case "end":
		return rule.NewEndRule(cmd.Parts[0]), nil
	case "border":
		return rule.NewBorderRule(cmd.Parts[0]), nil
	default:
		return nil, fmt.Errorf("unknown document text command: %s", cmd.Name)
	}
}

func handlePrepend(cmd Command) (any, error) {
	if len(cmd.Parts) < 1 || cmd.Parts[0] == "" {
		return nil, fmt.Errorf("prepend requires 1 argument(s)")
	}
	return rule.NewPrependRule(cmd.Parts[0]), nil
}

func handleAppend(cmd Command) (any, error) {
	if len(cmd.Parts) < 1 || cmd.Parts[0] == "" {
		return nil, fmt.Errorf("append requires 1 argument(s)")
	}
	return rule.NewAppendRule(cmd.Parts[0]), nil
}

func handleSurround(cmd Command) (any, error) {
	if len(cmd.Parts) < 2 || cmd.Parts[0] == "" || cmd.Parts[1] == "" {
		return nil, fmt.Errorf("surround requires before and after text")
	}
	return rule.NewSurroundRule(cmd.Parts[0], cmd.Parts[1]), nil
}

func handleCols(cmd Command) (any, error) {
	if len(cmd.Parts) < 2 {
		return nil, fmt.Errorf("cols requires a pattern and column spec")
	}

	patternStr := cmd.Parts[0]
	if patternStr == "" {
		patternStr = `\s+`
	} else {
		patternStr = cmd.EscapePattern(patternStr)
	}

	pattern, err := rule.CompilePattern(patternStr, emptyFlags()...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in cols: %w", err)
	}

	spec, err := rule.ParseColumnSpec(cmd.Parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid column spec in cols: %w", err)
	}

	joiner := " "
	if len(cmd.Parts) >= 3 {
		joiner = cmd.Parts[2]
	}

	return rule.NewColumnsRule(pattern, spec, joiner), nil
}

func handleSplit(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, fmt.Errorf("split requires a pattern")
	}
	opts := cmd.Flags(1)
	return rule.NewSplitRule(cmd.Pattern(), opts...)
}

func handleInsert(cmd Command) (any, error) {
	if len(cmd.Parts) < 2 || cmd.Parts[0] == "" {
		return nil, fmt.Errorf("insert requires a pattern and text")
	}
	opts := cmd.Flags(2)
	return rule.NewInsertRule(cmd.Pattern(), cmd.Parts[1], opts...)
}

func handleXargs(cmd Command) (any, error) {
	if len(cmd.Parts) < 1 || cmd.Parts[0] == "" {
		return nil, fmt.Errorf("xargs requires a command")
	}
	return rule.NewXargsRule(cmd.Parts[0]), nil
}

func handleExec(cmd Command) (any, error) {
	if len(cmd.Parts) < 1 || cmd.Parts[0] == "" {
		return nil, fmt.Errorf("exec requires a command")
	}
	return rule.NewExecRule(cmd.Parts[0]), nil
}

func handleSubstitution(cmd Command) (any, error) {
	if cmd.Delimiter == ':' {
		if len(cmd.Parts) < 2 {
			return nil, fmt.Errorf("substitution requires a line range and replacement")
		}
		lineRange, err := rule.ParseLineRange(cmd.Parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid line range: %w", err)
		}
		return rule.NewSubLineNumRule(lineRange, cmd.Parts[1]), nil
	}

	if len(cmd.Parts) < 2 {
		return nil, fmt.Errorf("substitution requires pattern and replacement with trailing delimiter")
	}
	opts := cmd.Flags(2)
	return rule.NewSubstitutionRule(cmd.Pattern(), cmd.Parts[1], opts...)
}

func handlePrint(cmd Command) (any, error) {
	if cmd.Delimiter == ':' {
		if len(cmd.Parts) < 1 {
			return nil, fmt.Errorf("print requires a line range")
		}
		lineRange, err := rule.ParseLineRange(cmd.Parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid line range: %w", err)
		}
		return rule.NewPrintLineNumRule(lineRange), nil
	}

	if len(cmd.Parts) < 1 {
		return nil, fmt.Errorf("print requires a pattern")
	}
	opts, ctxOpts, err := cmd.FlagsAndContext(1)
	if err != nil {
		return nil, fmt.Errorf("print: %w", err)
	}
	pattern := cmd.Pattern()
	if ctxOpts.has {
		return rule.NewPrintContextRule(pattern, ctxOpts.before, ctxOpts.after, opts...)
	}
	return rule.NewPrintLineRule(pattern, opts...)
}

func handleDelete(cmd Command) (any, error) {
	if cmd.Delimiter == ':' {
		if len(cmd.Parts) < 1 {
			return nil, fmt.Errorf("delete requires a line range")
		}
		lineRange, err := rule.ParseLineRange(cmd.Parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid line range: %w", err)
		}
		return rule.NewDeleteLineNumRule(lineRange), nil
	}

	if len(cmd.Parts) < 1 {
		return nil, fmt.Errorf("delete requires a pattern")
	}
	opts, ctxOpts, err := cmd.FlagsAndContext(1)
	if err != nil {
		return nil, fmt.Errorf("delete: %w", err)
	}
	pattern := cmd.Pattern()
	if ctxOpts.has {
		return rule.NewDeleteContextRule(pattern, ctxOpts.before, ctxOpts.after, opts...)
	}
	return rule.NewDeleteLineRule(pattern, opts...)
}

func handleTake(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, err
	}
	opts := cmd.Flags(1)
	joiner := " "
	if len(cmd.Parts) > 2 {
		joiner = cmd.Parts[2]
	}
	return rule.NewTakeRule(cmd.Parts[0], joiner, opts...)
}

func handleRemove(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, err
	}
	opts := cmd.Flags(1)
	return rule.NewRemoveRule(cmd.Parts[0], opts...)
}

func handleTakePrint(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, err
	}
	opts := cmd.Flags(1)
	joiner := " "
	if len(cmd.Parts) > 2 {
		joiner = cmd.Parts[2]
	}
	return rule.NewTakePrintRule(cmd.Parts[0], joiner, opts...)
}

func handleRemovePrint(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, err
	}
	opts := cmd.Flags(1)
	return rule.NewRemovePrintRule(cmd.Parts[0], opts...)
}

func handleGroup(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, err
	}
	groupNum := int(cmd.Name[0] - '0')
	opts := cmd.Flags(1)
	return rule.NewGroupRule(cmd.Parts[0], groupNum, opts...)
}

func handleUniq(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, fmt.Errorf("uniq requires a pattern")
	}

	pattern := cmd.Pattern()
	opts := emptyFlags()
	groupNum := 0

	for i := 1; i < len(cmd.Parts); i++ {
		part := cmd.Parts[i]
		if part == "" {
			continue
		}
		if len(part) == 1 && part[0] >= '1' && part[0] <= '9' {
			groupNum = int(part[0] - '0')
		} else {
			if strings.Contains(part, "g") {
				opts = append(opts, rule.WithGlobal())
			}
			if strings.Contains(part, "i") {
				opts = append(opts, rule.WithIgnoreCase())
			}
		}
	}

	return rule.NewUniqPatternRule(pattern, groupNum, opts...)
}

func handleControl(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, err
	}
	pattern := cmd.Pattern()
	opts := cmd.Flags(1)

	switch cmd.Name {
	case "on":
		return rule.NewOnRule(pattern, opts...)
	case "off":
		return rule.NewOffRule(pattern, opts...)
	case "after":
		return rule.NewAfterRule(pattern, opts...)
	case "toggle":
		return rule.NewToggleRule(pattern, opts...)
	default:
		return nil, fmt.Errorf("unknown control command: %s", cmd.Name)
	}
}

func handleIf(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, fmt.Errorf("missing pattern in if condition")
	}
	opts := cmd.Flags(1)
	compiled, err := rule.CompilePattern(cmd.Pattern(), opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in if condition: %w", err)
	}
	return &condition{pattern: compiled, inverted: cmd.Inverted}, nil
}

func handleIfAny(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, fmt.Errorf("missing pattern in ifany condition")
	}
	opts := cmd.Flags(1)
	compiled, err := rule.CompilePattern(cmd.Pattern(), opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in ifany condition: %w", err)
	}
	return &ifAnyCondition{pattern: compiled, inverted: cmd.Inverted}, nil
}

func handleIfNone(cmd Command) (any, error) {
	if err := cmd.RequirePattern(); err != nil {
		return nil, fmt.Errorf("missing pattern in ifnone condition")
	}
	opts := cmd.Flags(1)
	compiled, err := rule.CompilePattern(cmd.Pattern(), opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in ifnone condition: %w", err)
	}
	return &ifNoneCondition{pattern: compiled, inverted: cmd.Inverted}, nil
}

func handleBetween(cmd Command) (any, error) {
	if len(cmd.Parts) < 2 || cmd.Parts[0] == "" || cmd.Parts[1] == "" {
		return nil, fmt.Errorf("between requires start and end patterns")
	}

	opts := cmd.Flags(2)
	startCompiled, err := rule.CompilePattern(cmd.EscapePattern(cmd.Parts[0]), opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid start pattern in between: %w", err)
	}
	endCompiled, err := rule.CompilePattern(cmd.EscapePattern(cmd.Parts[1]), opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid end pattern in between: %w", err)
	}

	return &betweenCondition{
		startPattern: startCompiled,
		endPattern:   endCompiled,
		inverted:     cmd.Inverted,
	}, nil
}

// ---------------------------------------------------------------------------
// Shared utilities (kept from original)
// ---------------------------------------------------------------------------

func isValidDelimiter(delimiter byte) bool {
	if delimiter == ':' || delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		return true
	}
	r := rune(delimiter)
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

// parseFlags reads a flags string and returns the corresponding RuleOptions,
// merged with any globalOpts (e.g. from --insensitive).
func parseFlags(flags string) []rule.RuleOption {
	var opts []rule.RuleOption
	opts = append(opts, globalOpts...)
	if strings.Contains(flags, "g") {
		opts = append(opts, rule.WithGlobal())
	}
	if strings.Contains(flags, "i") {
		opts = append(opts, rule.WithIgnoreCase())
	}
	return opts
}

// emptyFlags returns globalOpts when no explicit flags are present.
func emptyFlags() []rule.RuleOption {
	if len(globalOpts) == 0 {
		return nil
	}
	opts := make([]rule.RuleOption, len(globalOpts))
	copy(opts, globalOpts)
	return opts
}

// flagsFromParts extracts flags from the trailing element of a parts slice.
func flagsFromParts(parts []string, flagIndex int) []rule.RuleOption {
	if flagIndex < len(parts) && parts[flagIndex] != "" {
		return parseFlags(parts[flagIndex])
	}
	return emptyFlags()
}

// contextOptions holds parsed before/after context values.
type contextOptions struct {
	before int
	after  int
	has    bool
}

// parseContextOptions scans parts starting at startIndex for context options.
func parseContextOptions(parts []string, startIndex int) (contextOptions, error) {
	var ctx contextOptions
	for i := startIndex; i < len(parts); i++ {
		part := parts[i]
		if !strings.Contains(part, "=") {
			continue
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
func flagsAndOptionsFromParts(parts []string, startIndex int) ([]rule.RuleOption, contextOptions, error) {
	opts := emptyFlags()
	for i := startIndex; i < len(parts); i++ {
		if !strings.Contains(parts[i], "=") {
			flagStr := parts[i]
			if strings.Contains(flagStr, "g") {
				opts = append(opts, rule.WithGlobal())
			}
			if strings.Contains(flagStr, "i") {
				opts = append(opts, rule.WithIgnoreCase())
			}
		}
	}
	ctxOpts, err := parseContextOptions(parts, startIndex)
	return opts, ctxOpts, err
}

// splitByDelimiter splits a string by delimiter, respecting backslash escapes.
func splitByDelimiter(input string, delimiter byte) ([]string, error) {
	var parts []string
	var current strings.Builder

	i := 0
	for i < len(input) {
		ch := input[i]

		if ch == '\\' && i+1 < len(input) {
			next := input[i+1]
			if next == delimiter {
				current.WriteByte(delimiter)
				i += 2
				continue
			} else if next == '\\' {
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

	parts = append(parts, current.String())
	return parts, nil
}

// condition is a parser-internal type representing a parsed if/!if condition.
type condition struct {
	pattern  *regexp2.Regexp
	inverted bool
}

// ifAnyCondition is a parser-internal type for ifany/!ifany conditions.
type ifAnyCondition struct {
	pattern  *regexp2.Regexp
	inverted bool
}

// ifNoneCondition is a parser-internal type for ifnone/!ifnone conditions.
type ifNoneCondition struct {
	pattern  *regexp2.Regexp
	inverted bool
}

// betweenCondition is a parser-internal type for between/!between conditions.
type betweenCondition struct {
	startPattern *regexp2.Regexp
	endPattern   *regexp2.Regexp
	inverted     bool
}
