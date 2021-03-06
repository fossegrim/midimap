package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fossegrim/midimap/lang/matcher"
	"gitlab.com/gomidi/midi"
)

// matcherMatchesMessage reports whether m matches message.
func matcherMatchesMessage(m matcher.Matcher, msg midi.Message) bool {
	switch m := m.(type) {
	case matcher.MatcherWithoutLogicalOperator:
		var data int64
		// msg.Raw()[0] is status
		// msg.Raw()[1] is data1
		// msg.Raw()[2] is data2
		switch m.LeftOperand {
		case matcher.Data1:
			data = int64(msg.Raw()[1])
		case matcher.Data2:
			data = int64(msg.Raw()[2])
		default:
			panic("unreachable")
		}

		switch m.Operator {
		case matcher.LessThanOperator:
			return data < m.RightOperand
		case matcher.LessThanOrEqualToOperator:
			return data <= m.RightOperand
		case matcher.EqualToOperator:
			return data == m.RightOperand
		case matcher.UnequalToOperator:
			return data != m.RightOperand
		case matcher.GreaterThanOrEqualToOperator:
			return data >= m.RightOperand
		case matcher.GreaterThanOperator:
			return data >= m.RightOperand
		default:
			panic("unreachable")
		}
	case matcher.MatcherWithLogicalOperator:
		switch m.Operator {
		case matcher.LogicalAndOperator:
			return matcherMatchesMessage(m.LeftMatcher, msg) &&
				matcherMatchesMessage(m.RightMatcher, msg)
		case matcher.LogicalOrOperator:
			return matcherMatchesMessage(m.LeftMatcher, msg) ||
				matcherMatchesMessage(m.RightMatcher, msg)
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}

// getInByPortNumber retrieves the midi.In by number(should not be
// confused with index) portNumber from ins.
func getInByPortNumber(ins []midi.In, number uint64) (in midi.In, err error) {
	for _, innerIn := range ins {
		if uint64(innerIn.Number()) == number {
			in = innerIn
			return
		}

	}
	err = fmt.Errorf("no MIDI port by number %d", number)
	return
}

// parsePortNumber parses a port number.
func parsePortNumber(s string) (portNumber uint64, err error) {
	portNumber, err = strconv.ParseUint(s, 10, 0)
	if err != nil {
		err = errors.New("portnumber must be a valid unsigned integer")
	}
	return
}
