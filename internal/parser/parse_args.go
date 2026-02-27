package parser

import (
	"fmt"

	"github.com/colinta/ged/internal/rule"
)

// ParseArgs parses a list of CLI arguments into rules, handling { } blocks
// for conditional rules. Returns a flat list of LineRule and DocumentRule values.
func ParseArgs(args []string) ([]any, error) {
	rules, remaining, err := parseArgs(args)
	if err != nil {
		return nil, err
	}
	if len(remaining) > 0 {
		return nil, fmt.Errorf("unexpected '%s'", remaining[0])
	}
	return rules, nil
}

// parseArgs is the recursive workhorse. It consumes args until it hits "}" or
// runs out of input. Returns the parsed rules and any unconsumed args.
//
// When it encounters a condition (from "if/pattern/"), it expects "{" next,
// then recurses to collect inner rules, then expects "}".
func parseArgs(args []string) ([]any, []string, error) {
	var results []any

	for len(args) > 0 {
		if args[0] == "}" {
			// End of block — return so the caller can consume "}"
			return results, args, nil
		}
		if args[0] == "{" {
			return nil, nil, fmt.Errorf("unexpected '{'")
		}

		parsed, err := ParseRule(args[0])
		if err != nil {
			return nil, nil, err
		}
		args = args[1:]

		if cond, ok := parsed.(*condition); ok {
			innerParsed, remaining, err := collectBlock(args, "if condition")
			if err != nil {
				return nil, nil, err
			}
			args = remaining

			if hasDocRule(innerParsed) {
				docRules := buildDocRules(innerParsed)
				results = append(results, rule.NewConditionalDocRule(cond.pattern, cond.inverted, docRules))
			} else {
				var lineRules []rule.LineRule
				for _, p := range innerParsed {
					lineRules = append(lineRules, p.(rule.LineRule))
				}
				results = append(results, rule.NewConditionalLineRule(cond.pattern, cond.inverted, lineRules))
			}
		} else if cond, ok := parsed.(*betweenCondition); ok {
			innerParsed, remaining, err := collectBlock(args, "between condition")
			if err != nil {
				return nil, nil, err
			}
			args = remaining

			if hasDocRule(innerParsed) {
				docRules := buildDocRules(innerParsed)
				results = append(results, rule.NewBetweenDocRule(cond.startPattern, cond.endPattern, cond.inverted, docRules))
			} else {
				var lineRules []rule.LineRule
				for _, p := range innerParsed {
					lineRules = append(lineRules, p.(rule.LineRule))
				}
				results = append(results, rule.NewBetweenLineRule(cond.startPattern, cond.endPattern, cond.inverted, lineRules))
			}
		} else {
			results = append(results, parsed)
		}
	}

	return results, args, nil
}

// collectBlock consumes "{", inner rules, and "}" from args.
// Returns the inner rules and the remaining args after "}".
func collectBlock(args []string, context string) ([]any, []string, error) {
	if len(args) == 0 || args[0] != "{" {
		return nil, nil, fmt.Errorf("expected '{' after %s", context)
	}
	args = args[1:] // consume "{"

	innerParsed, remaining, err := parseArgs(args)
	if err != nil {
		return nil, nil, err
	}
	if len(remaining) == 0 || remaining[0] != "}" {
		return nil, nil, fmt.Errorf("expected '}'")
	}
	return innerParsed, remaining[1:], nil
}

// hasDocRule reports whether any element in parsed is a DocumentRule.
func hasDocRule(parsed []any) bool {
	for _, p := range parsed {
		if _, ok := p.(rule.DocumentRule); ok {
			return true
		}
	}
	return false
}

// buildDocRules converts a mixed list of LineRule/DocumentRule into a
// []DocumentRule by wrapping consecutive LineRules in ApplyAllRule.
// This is the same logic used in main.go's run().
func buildDocRules(parsed []any) []rule.DocumentRule {
	var docRules []rule.DocumentRule
	var pendingLineRules []rule.LineRule

	for _, p := range parsed {
		switch r := p.(type) {
		case rule.LineRule:
			pendingLineRules = append(pendingLineRules, r)
		case rule.DocumentRule:
			if len(pendingLineRules) > 0 {
				docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
				pendingLineRules = nil
			}
			docRules = append(docRules, r)
		}
	}

	if len(pendingLineRules) > 0 {
		docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
	}

	return docRules
}
