# Ava
A Discord bot written in Go. Named after the AI Ava from *Ex Machina*

Add the environment variable `BOT_TOKEN` in `.env` as your bot's secret token. Copy the `docker-compose.yml` file to your directory and use 
```
# docker compose up
```
to launch the bot. 

[![Docker Image CI](https://github.com/jakobottar/ava-go/actions/workflows/docker-image.yml/badge.svg)](https://github.com/jakobottar/ava-go/actions/workflows/docker-image.yml)

## ToDo List
* fix error handling/wrapping
* Role manager with reactions (like carl)
* !help page
* database for multi-server support
* ~~Mention sender in "pong" response~~
  * Add (real) ping stats to ping command
* ~~!echo command~~
* ~~Voice channel creator~~
* ~~Timer~~
  * ~~"I'll be on in 5 min" - everybody ever~~
  * ~~Better arg parsing~~
* ~~Package into container~~
* ~~logger~~
