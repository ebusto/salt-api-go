package event

import (
	"cmp"
	"encoding/json"
	"slices"
	"strings"
)

type HighStateResult struct {
	Changes  map[string]any `json:"changes"`
	Comment  string         `json:"comment"`
	Duration Duration       `json:"duration"`
	Function string         `json:"function"`
	ID       string         `json:"__id__"`
	Name     string         `json:"name"`
	Order    int            `json:"__run_num__"`
	Result   bool           `json:"result"`
	SLS      string         `json:"__sls__"`
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
