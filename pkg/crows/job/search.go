package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/desertfox/gograylog"
)

type Search struct {
	Type     string   `yaml:"type"`
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
}

func (s *Search) GoGraylogQuery(frequency int) gograylog.Query {
	return gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Frequency:   frequency,
		Limit:       10000,
	}
}

func (s *Search) Run(g gograylog.ClientInterface, frequency int) (*Result, error) {
	q := gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Frequency:   frequency,
		Limit:       10000,
	}

	b, err := g.Search(q)
	if err != nil {
		return &Result{}, err
	}

	count := bytes.Count(b, []byte("\n"))

	if count == 0 && len(b) > 0 {
		var j map[string]interface{}
		if err := json.Unmarshal(b, &j); err != nil {
			fmt.Printf("Error parsing json. data:%s\nerr:%s\n", b, err)
		}

		if val, ok := j["total_results"]; ok {
			return &Result{
				Count: int(val.(float64)),
				When:  time.Now(),
			}, nil
		}
	}

	//Remove csv headers
	if count > 2 {
		count -= 1
	}

	r := &Result{
		Count: count,
		When:  time.Now(),
	}

	return r, nil
}

func (s Search) BuildURL(host string, from, to time.Time) string {
	q := gograylog.Query{
		QueryString: s.Query,
		StreamID:    s.Streamid,
		Fields:      s.Fields,
		Limit:       10000,
	}

	return q.Url(host, from, to)
}

func (r Result) From(f int) time.Time {
	return r.When.Add(time.Duration(int64(-1) * int64(f) * int64(time.Minute)))
}

func (r Result) To() time.Time {
	return r.When
}
