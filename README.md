# slackcat

slackcat is a command line tool that posts messages to [Slack].

    $ echo "hello" | slackcat

## Installing

If you have a working go installation run `go get github.com/paulhammond/slackcat`. Prebuilt binaries will be provided soon.

## Configuring

First, create a [new Slack Incoming Webhook integration][new-webhook].

Then create a `/etc/slackcat.conf` file, and add your new webhook url:

    {
        "webhook_url":"https://my.slack.com/services/hooks/incoming-webhook?token=token"
    }

If you don't have permission to create `/etc/slackcat.conf` then you can create `~/.slackcat.conf` instead.

## Usage

slackcat will take each line from stdin and post it as a message to Slack:

    tail -F logfile | slackcat

If you'd prefer to provide a message on the command line, you can:

    sleep 300; slackcat "done"

By default slackcat will post each message as coming from "user@hostname". If you want a different username, use the `--name` flag:

    echo ":coffee:" | slackcat --name "coffeebot"

Slackcat will use the channel specified when you set up the incoming webhook. You can override this in the config file by adding a "channel" option, or you can use the `--channel` flag:

    echo "testing" | slackcat --channel #test



[Slack]: http://slack.com/
[new-webhook]: https://my.slack.com/services/new/incoming-webhook
