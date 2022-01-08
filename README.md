# mcsweeney v1.0 
media compiler-sharer with efficient editing now employed on youtube.

mcsweeney is designed to pull clips from various video/streaming platforms according to a defined strategy, edit and compile them, and share them for entertainment. 

![mcsweeney](https://i.ibb.co/s6B62S4/Mcsweeney.png) 

## supported platforms
mcsweeney finds and shares content by consuming apis for popular media platforms. It pulls content from a 'source' and shares it to a 'destination'. mcsweeney currently supports the following platforms:

| sources | destinations |
| --- | --- |
| twitch | youtube |


## install
### dependencies
```
- go version 1.17
- ffmpeg
- gcc compiler
```
Once you've installed these dependencies, simply clone the repo and run 'go install'.
```
git clone git@github.com:hrand1005/mcsweeney.git
cd mcsweeney
go install
```
See the following tutorial section to use mcsweeney for the first time.
## Tutorial

mcsweeney operates by consuming external apis. You will need to create developer accounts with the desired platforms, and list the path of the credentials file in your yaml config file. Here's how you should format them:
### Twitch
A simple file with two lines, clientID and token:
```
<clientID>
<token>
```
### Youtube
Generate oauth credentials for this using the [google developer console](https://console.developers.google.com/)

After you've configured your exteral apis, you are ready to define a strategy
for content creation. Each time you run mcsweeney, you will provide a yaml config 
for the strategy you want to run. See mcsweeney/examples/example.yaml for
reference. 

Finally, you are ready to run mcsweeney. 
```
go build
./mcsweeney myconfig.yaml
```

