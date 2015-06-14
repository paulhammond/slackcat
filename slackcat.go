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
	Username   *string `json:"username"`
	Proxy      string  `json:"proxy"`
}

func ReadConfig() (*Config, error) {
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
			return nil, err
		}

		json.NewDecoder(file)
		conf := Config{}
		err = json.NewDecoder(file).Decode(&conf)
		if err != nil {
			return nil, err
		}
		return &conf, nil
	}

	return nil, errors.New("Config file not found")
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

func (m SlackMsg) Post(WebhookURL string, Proxy string) error {
	encoded, err := m.Encode()
	if err != nil {
		return err
	}

	if len(Proxy) > 0 {
		proxyUrl, err := url.Parse(Proxy)
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

		if err != nil {
			return err
		}

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

func username() string {
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

	cfg, err := ReadConfig()
	if err != nil {
		log.Fatalf("Could not read config: %v", err)
	}

	// By default use "user@server", unless overridden by config. cfg.Username
	// can be "", implying Slack should use the default username, so we have
	// to check if the value was set, not just for a non-empty string.
	defaultName := username()
	if cfg.Username != nil {
		defaultName = *cfg.Username
	}

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: slackcat [-c #channel] [-n name] [-i icon] [message]")
	}

	channel := pflag.StringP("channel", "c", cfg.Channel, "channel")
	name := pflag.StringP("name", "n", defaultName, "name")
	icon := pflag.StringP("icon", "i", "", "icon")
	pflag.Parse()

	// was there a message on the command line? If so use it.
	args := pflag.Args()
	if len(args) > 0 {
		msg := SlackMsg{
			Channel:   *channel,
			Username:  *name,
			Parse:     "full",
			Text:      strings.Join(args, " "),
			IconEmoji: *icon,
		}

		err = msg.Post(cfg.WebhookUrl, cfg.Proxy)
		if err != nil {
			log.Fatalf("Post failed: %v", err)
		}
		return
	}

	// ...Otherwise scan stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := SlackMsg{
			Channel:   *channel,
			Username:  *name,
			Parse:     "full",
			Text:      scanner.Text(),
			IconEmoji: *icon,
		}

		err = msg.Post(cfg.WebhookUrl, cfg.Proxy)
		if err != nil {
			log.Fatalf("Post failed: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
}
