package event

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type HighStateResult struct {
	Changes  map[string]any `json:"changes"`
	Comment  string         `json:"comment"`
	Duration Duration       `json:"duration"`
	Function string         `json:"function"`
	ID       string         `json:"__id__"`
	Name     Name           `json:"name"`
	Order    int            `json:"__run_num__"`
	Result   bool           `json:"result"`
	SLS      string         `json:"__sls__"`
}

type Name string

func (n *Name) UnmarshalJSON(b []byte) error {
	// Try string
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*n = Name(s)
		return nil
	}

	// Try []string
	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil {
		*n = Name(strings.Join(arr, "!!"))
		return nil
	}

	// Neither string nor []string
	return fmt.Errorf("name: expected string or []string, got %s", string(b))
}

// HighState parses the job return as a highstate return.
func (e *JobReturn) HighState() ([]HighStateResult, error) {
	var results []HighStateResult

	for key, value := range e.Return.Result().Map() {
		var r HighStateResult

		if err := json.Unmarshal([]byte(value.Raw), &r); err != nil {
			return nil, err
		}

		// Determine the function from the key.
		//   "file_|-/etc/promtail/promtail.yaml_|-/etc/promtail/promtail.yaml_|-managed"
		if p := strings.Split(key, "_|-"); len(p) == 4 {
			r.Function = p[0] + "." + p[3]
		}

		results = append(results, r)
	}

	slices.SortFunc(results, func(a, b HighStateResult) int {
		return cmp.Compare(a.Order, b.Order)
	})

	return results, nil
}
