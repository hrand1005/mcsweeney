# mcsweeney v2.0 
mcsweeney creates Twitch-clip compilations, and uploads them to Youtube.

![mcsweeney](https://i.ibb.co/s6B62S4/Mcsweeney.png) 

## Setup
### Dependencies

Install ffmpeg. The default Ubuntu install works:
```
sudo apt update 
sudo apt install ffmpeg
```
If you choose to compile ffmpeg on your machine, be sure to configure your installation with the following options:
```
ffmpeg version 4.4.git Copyright (c) 2000-2022 the FFmpeg developers
  built with gcc 9 (Ubuntu 9.3.0-17ubuntu1~20.04)
  configuration: --enable-libx264 --enable-gpl --enable-gnutls
  libavutil      57. 27.100 / 57. 27.100
  libavcodec     59. 35.100 / 59. 35.100
  libavformat    59. 25.100 / 59. 25.100
  libavdevice    59.  6.100 / 59.  6.100
  libavfilter     8. 41.100 /  8. 41.100
  libswscale      6.  6.100 /  6.  6.100
  libswresample   4.  6.100 /  4.  6.100
  libpostproc    56.  5.100 / 56.  5.100
```
mcsweeney uses a sqlite database to track scraped clips. You need a SQLite driver for this, for example:
```
go get github.com/mattn/go-sqlite3
go install github.com/mattn/go-sqlite3
``` 
You may need to install or reinstall gcc. On Ubuntu:
```
sudo apt install --reinstall build-essential
```

### Credentials
mcsweeney uses Twitch and Youtube APIs. To configure mcsweeney with your credentials, create a dotfile of the following form:
```
TWITCH_CLIENT_ID="<twitch client id>"
TWITCH_CLIENT_TOKEN="<twitch client token>"
TWITCH_TOKEN_FILE="<path/to/twitch/token/file>"

YOUTUBE_CLIENT_ID="<youtube client id>"
YOUTUBE_CLIENT_TOKEN="<youtube client token>"
YOUTUBE_TOKEN_FILE="<path/to/youtube/token/file>"
```
The path to this file should be provided as ```--env=<filepath>``` when executing the mcsweeney binary. It's worth noting
that the token file paths don't need to contain files with working tokens -- if a working token isn't found at the token file path,
then a file will be created after mcsweeney acquires a new token using the given credentials.

### Configs
Check out ```configs/``` for example configurations for mcsweeney. The yaml config file you provide to mcsweeney defines the content
you will create. For example, you might create a config called ```melee.yaml```:
```
title: "Top Melee Clips of the Week"
game-id: "16282"
first: 10
days: 7
database: "melee-clips.sqlite"
```
The ```game-id```, ```first```, and ```days``` fields define the category and number of clips that mcsweeney will scrape using the
twitch API. In this case, mcsweeney will pull the first (top) 10 clips from created in the last 7 days for the provided game-id "16282", 
which happens to represent Super Smash Bros. Melee. ```database``` defines a sqlite file which will store scraped clips so as to prevent
duplicate scraping in subsequent runs of mcsweeney. If the provided database file doesn't exist, one will be created. Finally, ```title``` 
defines the title that will be assigned to the uploaded youtube video.

## Run
```
Usage of ./mcsweeney:
  -env string
        Path to file defining environment variables, may be overwritten
  -max-encoders int
        Maximum number of video encodings that can occur concurrently (default 1)
  -config string
        Path to mcsweeney config

Example:
    ./mcsweeney --env=.env --max-encoders=2 --config=configs/melee.yaml 
```
