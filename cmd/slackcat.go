package main

import (
	"bufio"
	"fmt"
	"log"
	"github.com/whosonfirst/slackcat"
	"github.com/ogier/pflag"
	"os"
	"os/user"
	"strings"
)

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

	cfg, err := slackcat.ReadConfig()

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
		msg := slackcat.SlackMsg{
			Channel:   *channel,
			Username:  *name,
			Parse:     "full",
			Text:      strings.Join(args, " "),
			IconEmoji: *icon,
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
		msg := slackcat.SlackMsg{
			Channel:   *channel,
			Username:  *name,
			Parse:     "full",
			Text:      scanner.Text(),
			IconEmoji: *icon,
		}

		err = msg.Post(cfg.WebhookUrl)
		if err != nil {
			log.Fatalf("Post failed: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
}
