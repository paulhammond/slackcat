// Copyright 2014 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"

	"github.com/ogier/pflag"
)

type Config struct {
	WebhookUrl string  `json:"webhook_url"`
	Channel    string  `json:"channel"`
	Username   string  `json:"username"`
	IconEmoji  string  `json:"icon_emoji"`
}

func (c *Config) Load() error {
	err := c.loadConfigFiles()
	if err != nil {
		return err
	}
	c.loadEnvVars()

	if c.WebhookUrl == "" {
		return errors.New("Could not find a WebhookUrl in SLACKCAT_WEBHOOK_URL, /etc/slackcat.conf, /.slackcat.conf, ./slackcat.conf")
	}

	return nil
}

func (c *Config) loadEnvVars() {
	envs := []string{"SLACKCAT_WEBHOOK_URL", "SLACKCAT_CHANNEL", "SLACKCAT_USERNAME", "SLACKCAT_ICON"}
	for _,env := range envs {
		envVal := os.Getenv(env)
		if envVal == "" {
			continue
		}

		switch env {
		case "SLACKCAT_WEBHOOK_URL":
			c.WebhookUrl = envVal
		case "SLACKCAT_CHANNEL":
			c.Channel = envVal
		case "SLACKCAT_USERNAME":
			c.Username = envVal
		case "SLACKCAT_ICON":
			c.IconEmoji = envVal
		}
	}
}

func (c *Config) loadConfigFiles() error {
	homeDir := ""
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	for _, path := range []string{"/etc/slackcat.conf", homeDir + "/.slackcat.conf", "./slackcat.conf"} {
		file, err := os.Open(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}

		err = json.NewDecoder(file).Decode(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) BindFlags() {
	pflag.StringVarP(&c.Channel, "channel", "c", c.Channel, "channel")
	pflag.StringVarP(&c.Username, "name", "n", c.Username, "name")
	pflag.StringVarP(&c.IconEmoji, "icon", "i", c.IconEmoji, "icon")
}

type SlackMsg struct {
	Channel   string `json:"channel"`
	Username  string `json:"username,omitempty"`
	Text      string `json:"text"`
	Parse     string `json:"parse"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

func (m SlackMsg) Encode() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m SlackMsg) Post(WebhookURL string) error {
	encoded, err := m.Encode()
	if err != nil {
		return err
	}

	resp, err := http.PostForm(WebhookURL, url.Values{"payload": {encoded}})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Not OK")
	}
	return nil
}

func defaultUsername() string {
	username := "<unknown>"
	usr, err := user.Current()
	if err == nil {
		username = usr.Username
	}

	hostname := "<unknown>"
	host, err := os.Hostname()
	if err == nil {
		hostname = host
	}
	return fmt.Sprintf("%s@%s", username, hostname)
}

func main() {
	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: slackcat [-c #channel] [-n name] [-i icon] [message]")
	}

	cfg := Config{Username: defaultUsername()}
	err := cfg.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg.BindFlags()
	pflag.Parse()

	// was there a message on the command line? If so use it.
	args := pflag.Args()
	if len(args) > 0 {
		msg := SlackMsg{
			Channel:   cfg.Channel,
			Username:  cfg.Username,
			Parse:     "full",
			Text:      strings.Join(args, " "),
			IconEmoji: cfg.IconEmoji,
		}

		err = msg.Post(cfg.WebhookUrl)
		if err != nil {
			log.Fatalf("Post failed: %v", err)
		}
		return
	}

	// ...Otherwise scan stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := SlackMsg{
			Channel:   cfg.Channel,
			Username:  cfg.Username,
			Parse:     "full",
			Text:      scanner.Text(),
			IconEmoji: cfg.IconEmoji,
		}

		err = msg.Post(cfg.WebhookUrl)
		if err != nil {
			log.Fatalf("Post failed: %v", err)
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
}
