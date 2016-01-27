# About slackcat

Slackcat is a command line tool that posts messages to [Slack].

    $ echo "hello" | slackcat
    $ slackcat "world"

## Installing

To build them at your own, you need to define `GOPATH`:

    $ export GOPATH="$HOME/go/"

Then you have to download dependency:

    $ go get github.com/ogier/pflag

After this you can build slackcat:

    $ go build slackcat.go

## Configuring

First, create a [new Slack Incoming Webhook integration][new-webhook].

Then create a `/etc/slackcat.conf` file, and add your new webhook url:

```json
{
    "webhook_url":"https://hooks.slack.com/services/.../.../token",
    "channel":"#myslackchannel",
    "username":"SlackCat"
}
```

If you don't have permission to create `/etc/slackcat.conf` then you can create `~/.slackcat.conf` or `./slackcat.conf` instead.

## Usage

Slackcat will take each line from stdin and post it as a message to Slack:

    tail -F logfile | slackcat

If you'd prefer to provide a message on the command line, you can:

    sleep 300; slackcat "done"

By default slackcat will post each message as coming from "user@hostname". If you want a different username, use the `--name` flag:

    echo ":coffee:" | slackcat --name "coffeebot"

Slackcat will use the channel specified when you set up the incoming webhook. You can override this in the config file by adding a "channel" option, or you can use the `--channel` flag:

    echo "testing" | slackcat --channel #test



[Slack]: http://slack.com/
[new-webhook]: https://my.slack.com/services/new/incoming-webhook
