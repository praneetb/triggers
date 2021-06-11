package jiraclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jira "github.com/andygrunwald/go-jira"
	cache "github.com/patrickmn/go-cache"
	jwt "github.com/rbriski/atlassian-jwt"
	"github.com/sirupsen/logrus"
)

// SecurityConfig holds the information from the installation handshake
type SecurityConfig struct {
	Key            string `json:"key"`
	ClientKey      string `json:"client_key"`
	PublicKey      string `json:"public_key"`
	SharedSecret   string `json:"shared_secret"`
	ServerVersion  string `json:"server_version"`
	PluginsVersion string `json:"plugins_version"`
	BaseURL        string `json:"base_url"`
	ProductType    string `json:"product_type"`
	Description    string `json:"description"`
	EventType      string `json:"event_type"`
}

// QueryResult result of a jira query
type QueryResult struct {
	StartAt    int          `json:"start_at"`
	MaxResults int          `json:"max_results"`
	Total      int          `json:"total"`
	Issues     []jira.Issue `json:"issues"`
}

// JiraClient the jira client object
type JiraClient struct {
	BaseURL string
	Project string
	Cfg     jwt.Config
	client  *jira.Client
	Cache   *cache.Cache
}

//NewJiraClient client to talk to jira server
func NewJiraClient(url, file, project string) (*JiraClient, error) {
	jc := &JiraClient{
		BaseURL: url,
		Project: project,
		Cache:   cache.New(cache.NoExpiration, time.Millisecond),
	}

	content, err := ioutil.ReadFile(file)
	if err == nil {
		cfg := &jc.Cfg
		err1 := json.Unmarshal([]byte(content), cfg)
		if err1 != nil {
			logrus.Errorf("cannot unmarshall jwt config, error: %v", err)
			return nil, err1
		}
	}

	// http handlers
	http.HandleFunc("/atlassian-connect.json", ConnectHandler(url))

	http.HandleFunc("/installed", InstallHandler(&jc.Cfg, file))

	http.HandleFunc("/uninstalled", UninstallHandler(&jc.Cfg))

	http.HandleFunc("/event", EventHandler(&jc.Cfg))

	jc.GetAllIssues()

	return jc, nil
}

// GetAllIssues get all issues in the Jira
func (jc *JiraClient) GetAllIssues() {

	client, _ := jira.NewClient(jc.Cfg.Client(), jc.Cfg.BaseURL)

	// Check number of issues, then loop through  and get them all
	query := fmt.Sprintf("rest/api/2/search?jql=project=%s&maxResults=0", jc.Project)
	req, err := client.NewRequest("GET", query, nil)
	if err != nil {
		logrus.Errorf("Read issues failed, error: %v", err)
		return
	}

	results := new(QueryResult)
	_, err = client.Do(req, results)
	if err != nil {
		logrus.Errorf("Read issues request failed, error: %v", err)
		return
	}

	count := 0
	for count < results.Total {
		query := fmt.Sprintf("rest/api/2/search?jql=project=%s&fields=key,summary,description&maxResults=1000&startAt=%d", jc.Project, count)
		count += 1000
		req, err := client.NewRequest("GET", query, nil)
		if err != nil {
			logrus.Errorf("Read issues failed, error: %v", err)
			return
		}

		results = new(QueryResult)
		_, err = client.Do(req, results)
		if err != nil {
			logrus.Errorf("Read issues request failed, error: %v", err)
			return
		}
		for _, issue := range results.Issues {
			logrus.Debugf("Issue: %v", issue.Key)
			jc.Cache.Set(issue.Key, issue, cache.NoExpiration)
		}
	}

}
