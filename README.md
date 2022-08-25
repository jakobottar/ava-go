# Ava
A Discord bot written in Go. Named after the AI Ava from *Ex Machina*

[![Docker Image](https://github.com/jakobottar/ava-go/actions/workflows/stable-image.yml/badge.svg)](https://github.com/jakobottar/ava-go/actions/workflows/stable-image.yml) 
[![Docker Image - Development](https://github.com/jakobottar/ava-go/actions/workflows/dev-image.yml/badge.svg)](https://github.com/jakobottar/ava-go/actions/workflows/dev-image.yml)

## Deployment Instructions
1. Fill out `.env.example` with the required info, rename it to `.env`. 

2. Invite your bot to your server. Make sure to invite it with an OAuth2 URL that includes `applications.commands` permissions. 

3. Copy the `docker-compose.yml` file to your directory and use 
    ```
    # docker compose up
    ``` 
    to launch the bot. 
