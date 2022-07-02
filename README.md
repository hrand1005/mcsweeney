# mcsweeney v1.0 
mcsweeney creates Twitch-clip compilations, and uploads them to Youtube. 

![mcsweeney](https://i.ibb.co/s6B62S4/Mcsweeney.png) 

## install
### dependencies

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
mcsweeney leverages the Twitch API. To configure your twitch credentials, create a dotfile of the following form:
```
CLIENT_ID="<client id>"
CLIENT_TOKEN="<client token>"
TWITCH_APP_TOKEN="<twitch app token>"
```
The path to this file should be provided as ```--env=<filepath>``` when executing the mcsweeney binary. It's worth noting
that the ```TWITCH_APP_TOKEN``` isn't strictly required -- if a token isn't available, mcsweeney will request a new token, and 
write it to this file. 

### Configs
Check out ```configs/``` for example configurations for mcsweeney. 

### How to Run
```
Usage of ./mcsweeney:
  -env string
        Path to file defining environment variables, may be overwritten
  -max-encoders int
        Maximum number of video encodings that can occur concurrently (default 1)
  -twitch-config string
        Path to twitch scraper configuration file

Example:
    ./mcsweeney --env=.env --max-encoders=2 --twitch-config=configs/melee.yaml 
```
