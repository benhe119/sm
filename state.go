package sm

import (
	"strconv"
	"strings"
)

const (
	_ State = iota
	STATE_CREATED
	STATE_MITIGATED
	STATE_FIXED
	STATE_CLOSED
)

// State ...
type State int

func ParseState(s string) State {
	switch strings.ToLower(s) {
	case "created":
		return State(STATE_CREATED)
	case "mitigated":
		return State(STATE_MITIGATED)
	case "fixed":
		return State(STATE_FIXED)
	case "closed":
		return State(STATE_CLOSED)
	default:
		i, err := strconv.Atoi(s)
		if err != nil {
			return State(0)
		}
		return State(i)
	}
}

func (s State) String() string {
	switch s {
	case STATE_CREATED:
		return "CREATED"
	case STATE_MITIGATED:
		return "MITIGATED"
	case STATE_FIXED:
		return "FIXED"
	case STATE_CLOSED:
		return "CLOSED"
	default:
		return "???"
	}
}
