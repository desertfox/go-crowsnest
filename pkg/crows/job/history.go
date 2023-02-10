package job

import "time"

const maxHistory int = 20

// History container for Results
type History struct {
	results    []*Result
	alertCount *int
}

type Result struct {
	Count int
	When  time.Time
	Alert bool
}

func newHistory() *History {
	return &History{
		results:    make([]*Result, 0, maxHistory),
		alertCount: new(int),
	}
}

// Results accessor
func (h History) Results() []*Result {
	return h.results
}

func (h *History) AlertCount() int {
	return *h.alertCount
}

// Add
func (h *History) Add(r *Result) {
	if r.Alert {
		*h.alertCount++
	} else {
		*h.alertCount = 0
	}

	results := []*Result{r}
	results = append(results, h.results...)

	if len(h.results) >= maxHistory {
		results = results[:maxHistory]
	}

	h.results = results
}

func (h History) Avg() int {
	if len(h.results) == 0 {
		return 0
	}

	var sum int = 0
	for _, v := range h.results {
		sum += v.Count
	}

	return sum / len(h.results)
}
