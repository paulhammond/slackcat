# slackcat

slackcat is a command line tool that posts messages to [Slack].

    $ echo "hello" | slackcat

## Installing

If you have a working go installation run `go get github.com/paulhammond/slackcat`. Prebuilt binaries will be provided soon.

## Configuring

To use slackcat you must create a [Slack Incoming Webhook integration][new-webhook].

You can configure slackcat through a config file or environment variables (the latter overriding the former).

Only `webhook_url` is required, `channel`, `username`, and `icon_emoji` are optional, and will be used if not supplied as a command line option.

### Environment Variable

    $ export SLACKCAT_WEBHOOK_URL=https://my.slack.com/services/hooks/incoming-webhook?token=token

### Config File

    {
        "webhook_url":"https://my.slack.com/services/hooks/incoming-webhook?token=token"
    }

In either `/etc/slackcat.conf`, `~/.slackcat.conf`, or `./slackcat.conf`.

## Usage

slackcat will take each line from stdin and post it as a message to Slack:

    tail -F logfile | slackcat

If you'd prefer to provide a message on the command line, you can:

    sleep 300; slackcat "done"

By default slackcat will post each message as coming from "user@hostname". If you want a different username, use the `--name` flag:

    echo ":coffee:" | slackcat --name "coffeebot"

Slackcat will use the channel specified when you set up the incoming webhook. You can override this in the config file by adding a "channel" option, or you can use the `--channel` flag:

    echo "testing" | slackcat --channel #test

You can set an avatar using an [emoji string](http://www.emoji-cheat-sheet.com/):

    echo "we're watching you" | slackcat --icon=:family:



[Slack]: http://slack.com/
[new-webhook]: https://my.slack.com/services/new/incoming-webhook
