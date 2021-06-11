package jiraclient

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"

	jira "github.com/andygrunwald/go-jira"
	jwt "github.com/rbriski/atlassian-jwt"
	"github.com/sirupsen/logrus"
)

// IssueEvent holds all issue change data
type IssueEvent struct {
	Timestamp          int64  `json:"timestamp"`
	WebhookEvent       string `json:"webhookEvent"`
	IssueEventTypeName string `json:"issue_event_type_name"`
	Issue              struct {
		ID   string `json:"id"`
		Self string `json:"self"`
		Key  string `json:"key"`
	} `json:"issue"`
}

// ConnectHandler default endpoint handler
func ConnectHandler(baseURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"url path": r.URL.Path,
			"base url": baseURL,
		}).Info("GET for ConnectHandler URL path")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		fp := path.Join("./templates", "atlassian-connect.json")
		templateVals := map[string]string{
			"BaseURL": baseURL,
		}
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			logrus.Errorf("Template Parsing failed, error: %v", err)
			return
		}
		tmpl.Execute(w, templateVals)
		if err != nil {
			logrus.Errorf("Template execution failed, error: %v", err)
			return
		}
		return
	}
}

// InstallHandler default endpoint handler
func InstallHandler(cfg *jwt.Config, contextFile string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"url path": r.URL.Path,
		}).Info("GET for InstallHandler URL path")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("http read failed, error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sc := SecurityConfig{}
		err = json.Unmarshal(body, &sc)
		if err != nil {
			logrus.Errorf("json unmarshal failed, error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cfg.Key = sc.Key
		cfg.ClientKey = sc.ClientKey
		cfg.SharedSecret = sc.SharedSecret
		cfg.BaseURL = sc.BaseURL

		file, _ := json.MarshalIndent(cfg, "", " ")

		err = ioutil.WriteFile(contextFile, file, 0644)
		if err != nil {
			logrus.Errorf("jira context write to file failed, error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode([]string{"OK"})
	}
}

// UninstallHandler default endpoint handler
func UninstallHandler(cfg *jwt.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"url path": r.URL.Path,
		}).Info("GET for UninstallHandler URL path")

		w.WriteHeader(http.StatusOK)
	}
}

// EventHandler default endpoint handler
func EventHandler(cfg *jwt.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"url path": r.URL.Path,
		}).Info("GET for EventHandler URL path")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("http read failed, error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var ie IssueEvent
		json.Unmarshal(body, &ie)

		logrus.Debugf("base url :%v\n", cfg.BaseURL)
		logrus.Debugf("ISSUE EVENT:%v\n", ie)

		jiraClient, _ := jira.NewClient(cfg.Client(), cfg.BaseURL)
		issue, _, err := jiraClient.Issue.Get(ie.Issue.Key, nil)
		if err != nil {
			logrus.Errorf("jira issue get failed, error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logrus.Debug("ISSUE INFO:\n")
		issueJSON, err := json.MarshalIndent(issue, "", "    ")
		if err != nil {
			logrus.Errorf("json marshal failed, error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logrus.Debug(string(issueJSON))

		json.NewEncoder(w).Encode([]string{"OK"})

	}
}
