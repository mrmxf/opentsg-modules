// package middleware runs input middleware this is not currently not used
package middleware

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Functions is a skelton for the actions that may be made
type Functions struct {
	AnyOf  []string                   `json:"anyOf"`
	Action map[string]json.RawMessage `json:"action"`
}

// check checks if any middleware conditions have been completed.
func Check(functions []Functions, frame int) (bool, map[string]json.RawMessage, string) {

	frameCountReg := regexp.MustCompile(`^framecount\([\d]{1,4},[\d]{1,4}\)`)

	for _, funcs := range functions {
		var pass bool
		// check all the functions that are declared
		for _, funcToRun := range funcs.AnyOf {
			switch {
			case frameCountReg.MatchString(funcToRun):
				min, max := 0, 0
				// scan the digits from the text
				fmt.Sscanf(funcToRun, "framecount(%04d,%04d)", &min, &max)
				// run the function
				pass = framecount(min, max, frame)
			default:
				pass = false
			}
		}
		// if they all pass then return the action
		if pass {

			return pass, funcs.Action, "middleware functions :" + strings.Join(funcs.AnyOf, ",") // a map of the raw json to replace the one being used
		}
	}
	
	return false, nil, ""
}

func framecount(min, max, pos int) bool {
	if pos >= min && pos < max {

		return true
	}

	return false

}
