// References:
// Go maps to store entries: https://go.dev/blog/maps
// Check for value in map: https://golang.cafe/blog/how-to-check-if-a-map-contains-a-key-in-go#:~:text=The%20most%20common%20way%20to,if%20it%20does%20not%20exist.
// Initializing maps in Go: https://bitfieldconsulting.com/posts/map-declaring-initializing

package problem2

import (
	"strconv"
	"strings"
)

func ProcessLog(log []string, maxDuration int) int {
	// Input: entry logs (slice of strings), maxDuration (int)
	//        Log entry format: "visitor_num timestamp action"
	// Returns: total number of visitors who stayed less than or equal to the maximum duration time.

	timeLog := make(map[int][]int) // Store 'in' and 'out' time for each visitor
	ruleFollowers := 0

	// read over entries provided
	for _, entry := range log {
		// remove whitespaces  and split entry by comma
		entry = strings.Trim(entry, " ")
		parts := strings.Split(entry, ",")
		if len(parts) != 3 {
			continue
		}

		visitorNum, visErr := strconv.Atoi(strings.TrimSpace(parts[0]))
		timestamp, timeErr := strconv.Atoi(strings.TrimSpace(parts[1]))
		action := strings.TrimSpace(parts[2])

		// ensure visitor_num and timestamp are digits only
		if visErr != nil || timeErr != nil {
			continue
		}

		// Mapping: initialize map if visitor not yet tracked
		if action == "IN" {
			_, ok := timeLog[visitorNum]
			if !ok {
				timeLog[visitorNum] = []int{-1, -1}
			}
			timeLog[visitorNum][0] = timestamp

		} else if action == "OUT" {
			_, ok := timeLog[visitorNum]
			if !ok {
				timeLog[visitorNum] = []int{-1, -1}
			}
			timeLog[visitorNum][1] = timestamp
		}
	}

	// calculate duration of stay
	for _, timeLog := range timeLog {
		in := timeLog[0]
		out := timeLog[1]

		// calculate duration and record valid visitors
		if in != -1 && out != -1 {
			duration := out - in
			if duration <= maxDuration {
				ruleFollowers++
			}
		}
	}

	return ruleFollowers
}
