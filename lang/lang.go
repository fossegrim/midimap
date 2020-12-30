// Functions for parsing midimap-lang code.
//
// Before reading this code it is advised that you read through and comprehend the midimap-lang specification.
package lang

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ComparisonOperator int

const (
	LessThanOperator ComparisonOperator = iota
	LessThanOrEqualToOperator
	EqualToOperator
	UnequalToOperator
	GreaterThanOrEqualToOperator
	GreaterThanOperator
)

type Data1OrData2 int

const (
	Data1 Data1OrData2 = iota
	Data2
)

type Comparison struct {
	LeftOperand  Data1OrData2
	Operator     ComparisonOperator
	RightOperand int64
}

// parseCOMPARISON parses a COMPARISON as specified in Section 1.2.1.1 COMPARISONS of the midimap-lang specification.
// If s is a valid COMPARISON as described by the specification, parseCOMPARISON returns comparison, true.
// Otherwise parseCOMPARISON returns comparison, false.
func parseCOMPARISON(s string) (Comparison, bool) {
	s = strings.ReplaceAll(s, " ", "") // Remove spaces

	var comparison Comparison
	switch {
	case strings.HasPrefix(s, "data1"):
		comparison.LeftOperand = Data1
	case strings.HasPrefix(s, "data2"):
		comparison.LeftOperand = Data2
	default:
		return comparison, false
	}
	s = s[len("datax"):] // Discard parsed leftOperand

	var operatorLength int
	switch {
	case strings.HasPrefix(s, "=="):
		comparison.Operator = EqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "!="):
		comparison.Operator = UnequalToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "<="):
		comparison.Operator = LessThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, ">="):
		comparison.Operator = GreaterThanOrEqualToOperator
		operatorLength = 2
	case strings.HasPrefix(s, "<"):
		comparison.Operator = LessThanOperator
		operatorLength = 1
	case strings.HasPrefix(s, ">"):
		comparison.Operator = GreaterThanOperator
		operatorLength = 1
	default:
		return comparison, false
	}
	s = s[operatorLength:] // Discard parsed operator

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return comparison, false
	}
	comparison.RightOperand = int64(n)

	return comparison, true
}

type LogicalOperator int

const (
	LogicalAndOperator LogicalOperator = iota
	LogicalOrOperator
	NoLogicalOperator LogicalOperator = -1
)

type Matcher struct {
	LeftComparison  Comparison
	Operator        LogicalOperator
	RightComparison Comparison
}

// beforeAndAfter splits a string into two substrings, one before the leftmost match of a regexp and the other after it.
// If r matches s, beforeAndAfter returns (left, right, true) where left is the characters prior to the leftmost match and right is the characters to the right of the leftmost match.
// Otherwise, that is if r does not match s, ("", "", false) is returned.
func beforeAndAfter(r *regexp.Regexp, s string) (string, string, bool) {
	loc := r.FindStringIndex(s)
	if loc == nil {
		return "", "", false
	}

	before, after := s[:loc[0]], s[loc[1]:]
	return before, after, true

}

// beforeAndAfterLogicalOperator splits a string into two substrings, one before a logical operator and the other after it.
// If s contains "&&" beforeAndAfterLogicalOperator returns before, after, LogicalAndOperator, where before is the substring of s which appear before "&&" and after is the substring of s which appear after "&&".
// Otherwise if s contains "||" beforeAndAfterLogicalOperator returns before, after, LogicalOrOperator, where before is the substring of s which appear before "||" and after is the substring of s which appear after it.
// If s contains neither "||" nor "&&", beforeAndAfterLogicalOperator returns "", "", NoLogicalOperator
func beforeAndAfterLogicalOperator(s string) (string, string, LogicalOperator) {
	logicalAndRegexp := regexp.MustCompilePOSIX("&&")
	logicalOrRegexp := regexp.MustCompilePOSIX(`\|\|`)

	before, after, ok := beforeAndAfter(logicalAndRegexp, s)
	if ok {
		return before, after, LogicalAndOperator
	}

	before, after, ok = beforeAndAfter(logicalOrRegexp, s)
	if ok {
		return before, after, LogicalOrOperator
	}

	return "", "", NoLogicalOperator
}

// parseMATCHER parses a MATCHER as specified in Section 1.2.1 MATCHERS of the midimap-lang specification.
// If s is a valid MATCHER as described by the specification, parseMATCHER returns matcher, true.
// Otherwise parseMATCHER returns matcher, false.
func parseMATCHER(s string) (Matcher, bool) {
	var matcher Matcher

	var left, right string
	left, right, matcher.Operator = beforeAndAfterLogicalOperator(s)
	if matcher.Operator == NoLogicalOperator {
		return matcher, false
	}

	var ok bool
	matcher.LeftComparison, ok = parseCOMPARISON(left)
	if !ok {
		return matcher, false
	}
	matcher.RightComparison, ok = parseCOMPARISON(right)
	if !ok {
		return matcher, false
	}

	return matcher, true
}

// parseKeyCode parses a KeyCode of the following form.
// integer
// Spaces may be be intersped anywhere without changing the result.
// If s is of the specified form, it returns keyCode, true, otherwise it returns keyCode, false.
func parseKeyCode(s string) (int, bool) {
	s = strings.ReplaceAll(s, " ", "") // Remove spaces
	keyCode, err := strconv.Atoi(s)
	if err != nil {
		return keyCode, false
	}
	return keyCode, true
}

type Mapping struct {
	Matcher Matcher
	KeyCode int
}

// parseMAPPING parses a MAPPING as specified in Section 1.2 MAPPINGS of the midimap-lang specification.
// If s is a valid MAPPING as described by the specification, parseMAPPING returns mapping, true.
// Otherwise parseMAPPING returns mapping, false.
func parseMAPPING(s string) (Mapping, bool) {
	var mapping Mapping
	r := regexp.MustCompilePOSIX("- *>")
	before, after, ok := beforeAndAfter(r, s)
	if !ok {
		return mapping, false
	}
	mapping.Matcher, ok = parseMATCHER(before)
	if !ok {
		return mapping, false
	}
	mapping.KeyCode, ok = parseKeyCode(after)
	if !ok {
		return mapping, false
	}
	return mapping, true
}

// NextMAPPING attemps to parse the next MAPPING, as specified in Section 1.2 MAPPINGS of the midimap-lang specification, from r by parsing lines until a MAPPING is reached or an io error occurs.
// If an io error occured NextMAPPING returns mapping, ioError.
// If an invalid mapping is reached NextMAPPING returns mapping, err, where err is an error describing how the mapping is invalid.
// Otherwise NextMAPPING returns mapping, nil.
func NextMAPPING(r *bufio.Reader) (Mapping, error) {
	var mapping Mapping
	var line string
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			return mapping, err
		}
		line = s[:len(s)-1]
		// skip comments
		if !strings.HasPrefix(line, "#") {
			break
		}
	}

	// I wonder if there is an idiom to return all the return values of a called function.
	// It sounds a bit sugarish, so probably not.
	mapping, ok := parseMAPPING(line)
	var err error = nil
	if !ok {
		err = errors.New(fmt.Sprintf("invalid mapping %q", line))
	}

	return mapping, err
}
