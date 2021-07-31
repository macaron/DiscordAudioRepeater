# Discord Audio Repeater

[![dockeri.co](https://dockeri.co/image/rtlsdr/discord-audio-repeater)](https://hub.docker.com/r/rtlsdr/discord-audio-repeater)

Send audio data to voice channel from your audio file or URL. 

## Usage

### docker-compose

```bash
# Create a application & Invite bot account https://discordpy.readthedocs.io/en/stable/discord.html
$ cp .env.example .env
$ vi .env
$ docker-compose pull
$ docker-compose up -d --no-build
```
