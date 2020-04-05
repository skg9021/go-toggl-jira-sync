package main

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/google/logger"
	"github.com/joho/godotenv"
	"github.com/recoilme/slowpoke"
)

const logPath = "./app.log"

func init() {
	godotenv.Load(".env")
}

func syncTogglJira(session *Session, jiraClient *Client, cfg Config) {
	for _, entry := range session.getTogglEntries(cfg.Days) {
		value := getEntryFromDB(entry.Id)
		duration := strings.Split(value, " ")[0]
		jiraWorklogID := strings.Split(value, " ")[1]

		description := ""

		// if entry.Task.Name != entry.Description {
		// 	description = entry.Description
		// }
		description = entry.Description

		if duration != strconv.FormatInt(entry.Duration, 10) && entry.JiraId != "" {
			if jiraWorklogID != "0" {
				//update
				start := jira.Time(entry.Start)
				worklog, error := jiraClient.updateWorkLog(entry.Task.JiraId, description, jiraWorklogID, entry.Duration, start)

				if error == nil {
					setEntryInDB(entry.Id, strconv.FormatInt(entry.Duration, 10)+" "+worklog.ID)
				} else {
					logger.Fatalf("[jira] Task %s error: %v", entry.Task.JiraId, error)
				}
			} else {
				//new
				start := jira.Time(entry.Start)
				worklog, error := jiraClient.addWorkLog(entry.JiraId, description, entry.Duration, start)

				if error == nil {
					setEntryInDB(entry.Id, strconv.FormatInt(entry.Duration, 10)+" "+worklog.ID)
				} else {
					logger.Errorf("[jira] Task %s error: %v", entry.Task.JiraId, error)
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	lf, error := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if error != nil {
		logger.Fatalf("Failed to open log file: %v", error)
	}

	defer slowpoke.CloseAll()
	defer lf.Close()

	defer logger.Init("Logger", false, true, lf).Close()

	config := NewConfig()

	for _, user := range getUsers() {
		togglSession := getTogglSession(user.TogglToken)
		session := &Session{togglSession}
		jiraClient, _ := getJiraClient(user.JiraLogin, user.JiraToken)
		syncTogglJira(session, jiraClient, config)
	}

}
