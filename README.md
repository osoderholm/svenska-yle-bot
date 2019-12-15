# svenska-yle-bot

[![Go Report Card](https://goreportcard.com/badge/github.com/osoderholm/svenska-yle-bot)](https://goreportcard.com/report/github.com/osoderholm/svenska-yle-bot)

Telegram bot for [Svenska Yle news feed](https://svenska.yle.fi).

[Chat with SvenskaYleBot on Telegram](https://t.me/SvenskaYleBot)

## Why

WhatsApp ended support for sending mass-group messages, or something like that. 
Svenska Yle had a [questionnaire](https://svenska.yle.fi/artikel/2019/12/12/whatsapp-stoppade-grupputskicken-hjalp-oss-att-hitta-nya-satt-att-na-dig) regarding other messaging platforms, 
where they could continue sending out news updates to readers.

Telegram offers channels, that allow you to send out messages to your "followers", 
without having to use a chaotic group where everyone can type anything. 
This is fine, but it requires work. Of course in hindsight, someone probably gets paid for doing posting these news.

Anyway, this bot can do even more! You can add it to channels (contact me), 
but it can also be added to groups and you can even get news directly from it as a user.

## Features

The bot can act on different commands sent by users.

| Command | Description | Private | Group | Channel |
|:---------|:------------|---------|-------|---------|
| /latest  | Returns 5 latest articles| Yes | Yes | No |
| /subscribe | Subscribes you to news feed. You will receive latest news every hour-ish | Yes | No | No |
| /unsubscribe | Ends your subscription if you had one | Yes | No | No |

You can use `/help` to see available commands.

## Developers

### Running

#### Docker

**DISCLAIMER! Don't trust GitHub docker registry `latest` tag!**

There is a `docker-compose.yml` file containing an example of how to run the bot using Docker.

The only thing you need is a Telegram API token. You can get it using Telegram BotFather.

#### Local run
In order to run this locally, you must have: 

* Go 1.13
* Postgres 12 database
* Environmental variables set

Fetch dependencies:

    go get -u ./...
    
Run:

    go run .

### Database

The bot needs a database and currently that needs to be a Postgres database. 
This database can run anywhere, as long as the bot can access it. 
The connection params are specified using environmental variables.

The bot will automatically perform migrations and at the same time create a schema `bot` where all tables will go.

### Environmental variables

| Name | Description |
|:-----|:------------|
| TG_API_TOKEN | Telegram API Token. Get it from BotFather. |
| CHANNEL_SUBSCRIBERS | A comma separated list of channels that subscribe to news. Channels begin with @, no spaces. Example: `"@channelname,@channel2"`. The bot must have the right to post to the channel. Can be empty or left out completely. |
| DB_HOST | Database host address. |
| DB_PORT | Database host port. |
| DB_NAME | Database name. |
| DB_USER | Database username. |
| DB_PASSWORD | Database password for given DB_USER. |    
