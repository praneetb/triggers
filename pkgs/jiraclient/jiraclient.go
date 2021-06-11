package jiraclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	jwt "github.com/rbriski/atlassian-jwt"
	"github.com/sirupsen/logrus"
)

// SecurityConfig holds the information from the installation handshake
type SecurityConfig struct {
	Key            string `json:"key"`
	ClientKey      string `json:"clientKey"`
	PublicKey      string `json:"publicKey"`
	SharedSecret   string `json:"sharedSecret"`
	ServerVersion  string `json:"serverVersion"`
	PluginsVersion string `json:"pluginsVersion"`
	BaseURL        string `json:"baseUrl"`
	ProductType    string `json:"productType"`
	Description    string `json:"description"`
	EventType      string `json:"eventType"`
}

// JiraClient the jira client object
type JiraClient struct {
	BaseURL string
	Cfg     jwt.Config
	client  *jira.Client
}

//NewJiraClient client to talk to jira server
func NewJiraClient(url, file string) (*JiraClient, error) {
	jc := &JiraClient{
		BaseURL: url,
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

	return jc, nil
}
