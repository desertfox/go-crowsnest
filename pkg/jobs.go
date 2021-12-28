package crowsnest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/desertfox/crowsnest/pkg/graylog/search"
	"github.com/desertfox/crowsnest/pkg/graylog/session"
	"github.com/desertfox/crowsnest/pkg/teams/report"
	"gopkg.in/yaml.v2"
)

type job struct {
	Name      string        `yaml:"name"`
	Frequency int           `yaml:"frequency"`
	Threshold int           `yaml:"threshold"`
	TeamsURL  string        `yaml:"teamsurl"`
	Search    searchOptions `yaml:"options"`
}

type searchOptions struct {
	Username string   `yaml:"envusername"`
	Password string   `yaml:"envpassword"`
	Host     string   `yaml:"host"`
	Type     string   `yaml:"type"`
	Streamid string   `yaml:"streamid"`
	Query    string   `yaml:"query"`
	Fields   []string `yaml:"fields"`
}

func (s searchOptions) getUsername() string {
	return os.Getenv(s.Username)
}
func (s searchOptions) getPassword() string {
	return os.Getenv(s.Password)
}

func BuildJobsFromConfig(configPath string) []job {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	data := make(map[string][]job)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		panic(err.Error())
	}

	return data["jobs"]
}

func (j job) NewSession(httpClient *http.Client) sessionService {
	return session.New(j.Search.Host, j.Search.getUsername(), j.Search.getPassword(), httpClient)
}

func (j job) NewSearch(httpClient *http.Client) queryService {
	return search.New(
		j.Search.Host,
		j.Search.Query,
		j.Search.Streamid,
		j.Frequency,
		j.Search.Fields,
		httpClient,
	)
}

func (j job) NewReport() reportService {
	return report.Report{}
}

func (j job) GetCron(searchService searchService, reportService reportService) func() {
	return func() {
		j := j //MARK

		count, err := searchService.ExecuteSearch(searchService.GetHeader())
		if err != nil {
			panic(err.Error())
		}

		reportService.Send(
			j.TeamsURL,
			j.Name,
			fmt.Sprintf("Alert: %s\nCount: %d\nLink: [GrayLog Query](%s)\n", j.shouldAlertText(count), count, searchService.BuildSearchURL()),
		)
	}
}

func (j job) shouldAlertText(count int) string {
	if count >= j.Threshold {
		return fmt.Sprintf("ALERT %d/%d", count, j.Threshold)
	}

	return fmt.Sprintf("OK %d/%d", count, j.Threshold)
}
