package parameters

import "net/url"

// URLValuesToMap will generate a string map from URL values
func URLValuesToMap(values url.Values) map[string]string {
	r := make(map[string]string)
	for k, v := range values {
		if v != nil && len(v) != 0 {
			r[k] = v[0]
		}
	}
	return r
}
