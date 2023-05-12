package scheduler

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/nicksedov/sbconn-bot/pkg/cli"
	"github.com/nicksedov/sbconn-bot/pkg/settings"
)

type Secrets struct {
	BotToken string    `yaml:"BotToken"`
	OpenAIToken string `yaml:"OpenAIToken"`
}

func TestSchedule(t *testing.T) {
	*cli.FlagConfig = "../../sbconn-settings.yaml"
	secrets := getSecrets()
	*cli.FlagBotToken = secrets.BotToken
	*cli.FlagOpenAIToken = secrets.OpenAIToken
	cfg := settings.GetSettings()
	cfg.Events.Once = cfg.Events.Once[0:2] // Drop all records except two
	
	now := time.Now()
	fmt.Printf("Time = %s", now.Local().Format(time.RFC3339))
	event := &cfg.Events.Once[0]
	event.Destination = "testgroup-prompt"
	event.Moment = now.Add(500 * time.Millisecond)
	event.PromptRef = "have.a.nice.day"
	event = &cfg.Events.Once[1]
	event.Destination = "testgroup-prompt"
	event.Moment = now.Add(1500 * time.Millisecond)
	event.PromptRef = "current.date"
	event.MessageArgs = time.Now().Format("January 02")
	Schedule()
}

func getSecrets() Secrets {
	secrets := Secrets{}
	yfile, ioErr := ioutil.ReadFile("../../secrets.yaml")
	if ioErr == nil {
		ymlErr := yaml.Unmarshal(yfile, &secrets)
		if ymlErr != nil {
			log.Fatal(ymlErr)
		}
	}
	return secrets
}
