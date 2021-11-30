package main


import (
	"fmt"
	"github.com/nicklaw5/helix"
	"gopkg.in/yaml.v2"
    "io"
	"io/ioutil"
	"log"
	"net/http"
    "os"
	"os/exec"
	"strings"
    //"sync"
    "time"
)


const (
	RawVidsDir = "tmp/raw"
    ProcessedVidsDir = "tmp/processed"
)


type context struct {
	ClientID string `yaml:"clientID"`
	Token    string `yaml:"token"`
    GameID string `yaml:"gameID"` 
    First int `yaml:"first"`
}


func main() {
	// Get context from yaml file
	c := context{}

	err := loadContext("example.yaml", &c)
	if err != nil {
		fmt.Println("Couldn't load context.")
		log.Fatal(err)
	}

	// Remember, this is a strategy, so it will be more like s.GetContent()
    clips, err := getClips(c.ClientID, c.Token, c.GameID, c.First)
	if err != nil || len(clips) == 0 {
		fmt.Println("Couldn't get content.")
		log.Fatal(err)
	}
    
	//s.EditContent()
    editClipsTimer := clipFuncTimer(editClips)
    err = editClipsTimer(clips)
    if err != nil {
        fmt.Printf("Couldn't some clips: %v\n", err)
    }

    //s.CompileContent()
    err = compileClips()
    if err != nil {
        fmt.Printf("Couldn't compile clips: %v\n", err)
    }
	//s.CompileContent()
	//s.ShareContent()

	return
}


// It may be appropriate to get more information than just a token
func loadContext(path string, c *context) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read %s: %w", path, err)
	}

	err = yaml.Unmarshal(raw, c)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal: %v", err)
	}

	return nil
}


func getClips(clientID string, token string, gameID string, first int) ([]helix.Clip, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID: clientID,
	})
	if err != nil {
		return nil, err
	}

	client.SetUserAccessToken(token)
	defer client.SetUserAccessToken("")

	// Define query for clips
	clipParams := &helix.ClipsParams{
		GameID: gameID,
        First: first,
	}

    // Execute query for clips, TODO: more error checking here?
	twitchResp, err := client.GetClips(clipParams)
	if err != nil {
		return nil, err
	}

	err = downloadNewClips(twitchResp.Data.Clips)
	if err != nil {
		return nil, err
	}

	return twitchResp.Data.Clips, nil
}


func downloadNewClips(manyClips []helix.Clip) error {
    start := time.Now()
    fmt.Println("Download start...")

    //var wg sync.WaitGroup

    // TODO: spawn goroutines for each download
    // NOTE: Preliminary testing indicates not much of a difference around 14
    // clips downloaded using this commented out method. 
    fmt.Printf("Trying to download %v clips.\n", len(manyClips))
	for _, v := range manyClips {
		// TODO: Verify clip here
        //wg.Add(1)
        //v := v
		fmt.Println("Attempting to download a clip...")
        //go func() {
            //defer wg.Done()
        err := downloadClip(&v)
        if err != nil {
            fmt.Println("Failed to download a clip: ", err)
        }
        //}()
	}

    //wg.Wait()
    finish := time.Now()
    elapsed := finish.Sub(start)
    fmt.Printf("finished in %v\n", elapsed)

	return nil
}


func downloadClip(clip *helix.Clip) error {
	thumbURL := clip.ThumbnailURL
	mp4URL := strings.SplitN(thumbURL, "-preview", 2)[0] + ".mp4"
	fmt.Println("MP4 URL: ", mp4URL)

	resp, err := http.Get(mp4URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filename := strings.SplitN(mp4URL, "twitch.tv", 2)[1]
	outFile := RawVidsDir + filename

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}


//TODO: get rid of this
func getClipPath(clip *helix.Clip) string {
	thumbURL := clip.ThumbnailURL
	mp4URL := strings.SplitN(thumbURL, "-preview", 2)[0] + ".mp4"
	filename := strings.SplitN(mp4URL, "twitch.tv", 2)[1]

	return filename
}
    

// TODO: replace this with a ffmpeg library dear god
// TODO: goroutines!
func editClips(clips []helix.Clip) error {
    f, err := os.Create("clips.txt")
    if err != nil {
        return err
    }
    defer f.Close()

    for _, v := range clips {
        overlayText := fmt.Sprintf("%s\n%s", v.Title, v.BroadcasterName)
        filename := getClipPath(&v)
        rawPath := RawVidsDir + filename
        overlayArg := fmt.Sprintf("drawtext=fontfile=/usr/share/fonts/noto/NotoSansTamilUI-Regular.ttf:text='%s':fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=0:y=0", overlayText) 
        processedPath := ProcessedVidsDir + filename
        cmdName := "ffmpeg"
        args := []string{"-i", rawPath, "-vf", overlayArg, "-codec:a", "copy", processedPath}
        ffmpegCmd := exec.Command(cmdName, args...)
        err := ffmpegCmd.Run()
        if err != nil {
            fmt.Printf("Failed to execute ffmpeg cmd: %v\n", err)
        }


        w := fmt.Sprintf("file '%s'\n", processedPath)
        _, err = f.WriteString(w)
        if err != nil {
            return err
        }
    }
    
    return nil
}


func compileClips() error {
    cmdName := "ffmpeg"
    args := []string{"-f", "concat", "-safe", "0", "-i", "clips.txt", "compiled-vid.mp4"}
    cmd := exec.Command(cmdName, args...)
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Failed to execute ffmpeg cmd: %v\n", err)
    }

    return nil
}

// Consider this in the final version
/*
func main(){
    var get, edit, share

    // parse command line flags
    if -g then get = getContent(); etc.
    eachStrategy(get, edit, share)
}
*/


// some of that experimental stuff
type clipFunc func([]helix.Clip) error

func clipFuncTimer(f clipFunc) clipFunc {
    return func(c []helix.Clip) error {
        defer func(t time.Time) {
            fmt.Printf("clipFunc elapsed in %v\n", time.Since(t))
        }(time.Now())

        return f(c)
    }
}
