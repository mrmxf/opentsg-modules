// Package get, gets the information from map[string]any
package get

import "fmt"

// Get gets all the keys of a map and the associated vals
// arrayDot true includes array values as part of the dotpath, e,g,
// dot.path.0.val. Else the array is returned as the final value
func Get(extractMap map[string]any, prevKeys []string, arrayDot bool) ([]any, [][]string) {

	// loop through the map and get every pathway and each relative value
	count := 0
	var res []any
	var resKeys [][]string

	for key, extractValue := range extractMap {

		switch valueType := extractValue.(type) {
		case map[string]any:

			mapKeys := appendWithCopy(prevKeys, key) // append the keys on each run

			morev, moreK := Get(valueType, mapKeys, arrayDot)
			res = append(res, morev...)
			resKeys = append(resKeys, moreK...)
		case []any:
			if arrayDot {
				for i, arrVal := range valueType {
					arrKeys := appendWithCopy(prevKeys, key, fmt.Sprintf("%v", i))

					if m, ok := arrVal.(map[string]any); ok {
						morev, moreK := Get(m, arrKeys, arrayDot)
						res = append(res, morev...)
						resKeys = append(resKeys, moreK...)
					} else {
						res = append(res, arrVal)
						resKeys = append(resKeys, arrKeys)
					}
				}
			} else {
				valueKeys := appendWithCopy(prevKeys, key)

				res = append(res, valueType)
				resKeys = append(resKeys, valueKeys)
			}
		default:
			// when a value is found add to the list of values and strings
			valueKeys := appendWithCopy(prevKeys, key)

			res = append(res, valueType)
			resKeys = append(resKeys, valueKeys)
			// return that v type
		}
		count++
	}

	return res, resKeys
}

// appendWithCopy creates a copy of the previous string then appends the additions.
// It is designed to combat {appendAssign: append result not assigned to the same slice} and any
// of the errors that may stem from that.
func appendWithCopy(previous []string, additions ...string) []string {
	valueKeys := make([]string, len(previous))
	copy(valueKeys, previous)

	return append(valueKeys, additions...)
}
