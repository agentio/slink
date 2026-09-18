package froda

import (
	"fmt"
	"net/url"
)

// makeParams converts a map of string keys and any values into a URL-encoded string.
// If a value is a slice of strings, it will be joined with commas.
// Generally the values will be strings, numbers, booleans, or slices of strings
func makeParams(p map[string]any) string {
	params := url.Values{}
	for k, v := range p {
		if s, ok := v.([]string); ok {
			for _, v := range s {
				if v != "" {
					// omit parameters with empty values
					params.Add(k, v)
				}
			}
		} else {
			if v != "" {
				// omit parameters with empty values
				params.Add(k, fmt.Sprint(v))
			}
		}
	}
	return params.Encode()
}
