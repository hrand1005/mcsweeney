# mcsweeney v1.0
media compiler-sharer with efficient editing now employed on youtube
mcsweeney is designed to pull clips from various video/streaming platforms according to a defined strategy, edit and compile them, and share them for entertainment. 

## supported platforms
mcsweeney finds and shares content by consuming apis for popular media platforms. It pulls content from a 'source' and shares it to a 'destination'. mcsweeney currently supports the following platforms:

### sources
- twitch

### destinations
- youtube

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
Each time you run mcsweeney, you will provide a yaml config file:
```
go build
./mcsweeney myconfig.yaml
```
Your config file will define a strategy for mcsweeney to retrieve, compile, edit, and finally share your content. The skeleton looks like this:
```
name: <string>

# mcsweeney can optionally prepend an intro to your compiled video
intro: 
  path: "/path/to/intro/video.mp4"
  duration: <float seconds>
  # mcsweeney may overlay text on your intro. The text will be of the following form:
  # Start Date - End Date, Year
  overlay-start: <float seconds, defining the start time of a text overlay>
  font: "/path/to/font.ttf"

# mcsweeney can optionally append an outro to your compiled video
outro:
  path: "/path/to/outro.mp4"
  duration: <float seconds>

# define the source platform for your content. v1.0 supports twitch
source: 
  platform: "twitch"
  credentials: "path/to/credentials/file"
  query:
    game-id: "twitch api game id"
    first: <int number of clips>
    days: <int number of days backwards from today to search for clips>

# define language and blacklist filters -- only clips satisfying the language, 
# and creators off the blacklist, will be retrieved
filters:
  language: "en" 
  blacklist:
    - "blacklisted creator"

# define the destination platform, where your compiled video will be shared. v1.0 supports youtube
destination:
  platform: "youtube"
  credentials: "/path/to/youtube/credentials.json"
  title: "name of your title"
  category-id: "youtube category id"
  description: "your video's description"
  keywords: "comma,separated,keywords"
  privacy: "public"
  token-cache: "path to token0cache, /home/username/.credentials/youtube-go.json by default"

# define options for your custom overlay. this overlay will appear on your retrieved content, 
# and is made up of a background image and foreground text, crediting the creator
options:
  overlay:
    font: "path/to/overlay/font"
    color: "6 letter color string 'ffffff'"
    size: <int fontsize>
    fade: <float fade duration for the text>
    duration: <float duration that the overlay stays on screen>
    background: "/path/to/background/image"
```
mcsweeney operates by consuming external apis. You will need to create developer accounts with the desired platforms, and list the path of the credentials file in your yaml config file. Here's how you should format them:
### Twitch
A simple file with two lines, clientID and token:
```
<clientID>
<token>
```
### Youtube
Generate oauth credentials for this using the (google developer console)[https://console.developers.google.com/]

