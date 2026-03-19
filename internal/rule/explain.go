package rule

import "fmt"

// Explainer is an optional interface for rules that can describe
// their behavior in plain English. Used by --explain mode.
type Explainer interface {
	Explain() string
}

// --- Line rules ---

func (r *SubstitutionRule) Explain() string {
	if r.global {
		return fmt.Sprintf("replace all matches of /%s/ with %q", r.patternStr, r.replace)
	}
	return fmt.Sprintf("replace first match of /%s/ with %q", r.patternStr, r.replace)
}

func (r *SubLineNumRule) Explain() string {
	return fmt.Sprintf("replace matching line numbers with %q", r.replacement)
}

func (r *PrintLineRule) Explain() string {
	desc := fmt.Sprintf("keep lines matching /%s/", r.patternStr)
	if r.before > 0 || r.after > 0 {
		desc += fmt.Sprintf(" (with %d before, %d after context)", r.before, r.after)
	}
	return desc
}

func (r *PrintLineNumRule) Explain() string {
	return "keep lines by line number"
}

func (r *DeleteLineRule) Explain() string {
	desc := fmt.Sprintf("delete lines matching /%s/", r.patternStr)
	if r.before > 0 || r.after > 0 {
		desc += fmt.Sprintf(" (with %d before, %d after context)", r.before, r.after)
	}
	return desc
}

func (r *DeleteLineNumRule) Explain() string {
	return "delete lines by line number"
}

func (r *TakeRule) Explain() string {
	if r.global {
		return fmt.Sprintf("extract all matches of /%s/", r.patternStr)
	}
	return fmt.Sprintf("extract first match of /%s/", r.patternStr)
}

func (r *RemoveRule) Explain() string {
	return fmt.Sprintf("remove text matching /%s/", r.patternStr)
}

func (r *GroupRule) Explain() string {
	return fmt.Sprintf("extract capture group %d from /%s/", r.group, r.patternStr)
}

func (r *TrimRule) Explain() string {
	switch r.mode {
	case "left":
		return "trim leading whitespace"
	case "right":
		return "trim trailing whitespace"
	default:
		return "trim whitespace from both ends"
	}
}

func (r *UpperRule) Explain() string { return "convert to uppercase" }
func (r *LowerRule) Explain() string { return "convert to lowercase" }

func (r *PrependRule) Explain() string {
	return fmt.Sprintf("prepend %q to each line", r.text)
}

func (r *AppendRule) Explain() string {
	return fmt.Sprintf("append %q to each line", r.text)
}

func (r *SurroundRule) Explain() string {
	return fmt.Sprintf("surround each line with %q and %q", r.before, r.after)
}

func (r *ColumnsRule) Explain() string {
	return fmt.Sprintf("select/reorder columns (joined by %q)", r.joiner)
}

func (r *SplitRule) Explain() string {
	return "split each line on pattern into multiple lines"
}

func (r *InsertRule) Explain() string {
	return fmt.Sprintf("insert %q after lines matching /%s/", r.text, r.pattern)
}

func (r *OnRule) Explain() string    { return "start printing at first match" }
func (r *OffRule) Explain() string   { return "stop printing at first match" }
func (r *AfterRule) Explain() string { return "start printing after first match" }
func (r *ToggleRule) Explain() string {
	return "toggle printing on/off at each match"
}

func (r *XargsRule) Explain() string {
	return fmt.Sprintf("run %q with each line as argument", r.command)
}

// --- Document rules ---

func (r *SortRule) Explain() string    { return "sort lines alphabetically" }
func (r *ReverseRule) Explain() string { return "reverse line order" }
func (r *JoinRule) Explain() string {
	return fmt.Sprintf("join all lines with %q", r.separator)
}
func (r *LinesRule) Explain() string  { return "prepend line numbers" }
func (r *CountRule) Explain() string  { return "output line count" }
func (r *UniqRule) Explain() string {
	if r.pattern == nil {
		return "remove consecutive duplicates"
	}
	desc := fmt.Sprintf("remove consecutive duplicates matching /%s/", r.pattern.String())
	if r.groupNum > 0 {
		desc += fmt.Sprintf(" (group %d)", r.groupNum)
	}
	return desc
}
func (r *BeginRule) Explain() string  { return fmt.Sprintf("prepend %q to document", r.text) }
func (r *EndRule) Explain() string    { return fmt.Sprintf("append %q to document", r.text) }
func (r *BorderRule) Explain() string { return fmt.Sprintf("add %q border to document", r.text) }
func (r *ExecRule) Explain() string {
	return fmt.Sprintf("pipe document through %q", r.command)
}

// --- Conditional/composite rules ---

func (r *ConditionalLineRule) Explain() string {
	prefix := "if"
	if r.inverted {
		prefix = "if NOT"
	}
	desc := fmt.Sprintf("%s matching /%s/:", prefix, r.condition)
	for _, inner := range r.rules {
		if e, ok := inner.(Explainer); ok {
			desc += "\n    " + e.Explain()
		}
	}
	if len(r.elseRules) > 0 {
		desc += "\n  else:"
		for _, inner := range r.elseRules {
			if e, ok := inner.(Explainer); ok {
				desc += "\n    " + e.Explain()
			}
		}
	}
	return desc
}

func (r *ConditionalDocRule) Explain() string {
	prefix := "if"
	if r.inverted {
		prefix = "if NOT"
	}
	desc := fmt.Sprintf("%s matching /%s/:", prefix, r.condition)
	for _, inner := range r.rules {
		if e, ok := inner.(Explainer); ok {
			desc += "\n    " + e.Explain()
		}
	}
	return desc
}

func (r *BetweenLineRule) Explain() string {
	prefix := "between"
	if r.inverted {
		prefix = "outside"
	}
	desc := fmt.Sprintf("%s /%s/ and /%s/:", prefix, r.startPattern, r.endPattern)
	for _, inner := range r.rules {
		if e, ok := inner.(Explainer); ok {
			desc += "\n    " + e.Explain()
		}
	}
	return desc
}

func (r *BetweenDocRule) Explain() string {
	prefix := "between"
	if r.inverted {
		prefix = "outside"
	}
	desc := fmt.Sprintf("%s /%s/ and /%s/:", prefix, r.startPattern, r.endPattern)
	for _, inner := range r.rules {
		if e, ok := inner.(Explainer); ok {
			desc += "\n    " + e.Explain()
		}
	}
	return desc
}

func (r *IfAnyDocRule) Explain() string {
	prefix := "if any line matches"
	if r.inverted {
		prefix = "if no line matches"
	}
	desc := fmt.Sprintf("%s /%s/:", prefix, r.condition)
	for _, inner := range r.rules {
		if e, ok := inner.(Explainer); ok {
			desc += "\n    " + e.Explain()
		}
	}
	return desc
}

func (r *IfNoneDocRule) Explain() string {
	prefix := "if no line matches"
	if r.inverted {
		prefix = "if any line matches"
	}
	desc := fmt.Sprintf("%s /%s/:", prefix, r.condition)
	for _, inner := range r.rules {
		if e, ok := inner.(Explainer); ok {
			desc += "\n    " + e.Explain()
		}
	}
	return desc
}

func (r *ApplyAllRule) Explain() string {
	desc := "pipeline:"
	for _, inner := range r.rules {
		if e, ok := inner.(Explainer); ok {
			desc += "\n  " + e.Explain()
		}
	}
	return desc
}
