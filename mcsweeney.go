package main

/* TODO:
- Consistent error handling
- Testing
- Architecture
    - Think about params, pointers, interfaces, etc.
- Goroutines
- Decide on a standard for where to use pointers vs struct values
*/

//TODO: maybe we don't need the entire packages?
import (
	"fmt"
	"github.com/nicklaw5/helix"
	"google.golang.org/api/youtube/v3"
	"log"
	"mcsweeney/auth"
    "mcsweeney/config"
	"mcsweeney/db"
    "mcsweeney/get"
	"os"
	"os/exec"
	"strings"
	//"sync"
	"time"
)

const (
	RawVidsDir       = "tmp/raw"
	ProcessedVidsDir = "tmp/processed"
)

func main() {
	// Get config from yaml file
	c, err := config.NewConfig("config/example.yaml")
	if err != nil {
		fmt.Println("Couldn't load config.")
		log.Fatal(err)
	}

    // TODO: let's init the db here instead of later
	dbIntf, err := db.NewContentDB(c.Source)
	if err != nil {
		fmt.Println("Couldn't create content-db.")
		log.Fatal(err)
	}

    // TODO: change to .Init()
	err = dbIntf.Create()
	if err != nil {
		fmt.Println("Couldn't create DB.")
		log.Fatal(err)
	}

    getIntf, err := get.NewContentGetter(*c, dbIntf)
    if err != nil {
        fmt.Println("Couldn't create content-getter.")
        log.Fatal(err)
    }

    content, err := getIntf.GetContent()
    if err != nil {
        fmt.Println("Couldn't get new content.")
        log.Fatal(err)
    }

	//s.EditContent()
	editClipsTimer := clipFuncTimer(editClips)
	err = editClipsTimer(content)
	if err != nil {
		fmt.Println("Couldn't edit some clips")
		log.Fatal(err)
	}

	//s.CompileContent()
	err = compileClips()
	if err != nil {
		fmt.Println("Couldn't compile clips")
		log.Fatal(err)
	}

	//s.ShareContent()
	uploadArgs := uploadArgs{
		"compiled-vid.mp4",
		"McSweeney's title",
		"McSweeney's description",
		"McSweeney's category",
		"McSweeney's keywords",
		"private",
	}

	resp, err := uploadVideo(uploadArgs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Upload response: %v", resp)

	return
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
		fmt.Printf("Overlay text: %s\n", overlayText)
		filename := getClipPath(&v)
		rawPath := RawVidsDir + filename
		overlayArg := fmt.Sprintf(`drawtext=fontfile=/usr/share/fonts/TTF/DejaVuSans.ttf:text='%s':fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=0:y=0`, overlayText)
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

type uploadArgs struct {
	Filename    string
	Title       string
	Description string
	Category    string
	Keywords    string
	Privacy     string
}

func uploadVideo(args uploadArgs) (*youtube.Video, error) {
	if args.Filename == "" {
		return nil, fmt.Errorf("cannot upload nil file")
	}
	//TODO: perform checks on the inputs

	// TODO: figure out this
	client := auth.GetClient(youtube.YoutubeUploadScope)

	service, err := youtube.New(client)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create youtube service: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       args.Title,
			Description: args.Description,
			//TODO: this might be nice :)
			//CategoryId:  args.Category,
			Tags: strings.Split(args.Keywords, ","),
		},
		Status: &youtube.VideoStatus{PrivacyStatus: args.Privacy},
	}

	insertArgs := []string{"snippet", "status"}
	call := service.Videos.Insert(insertArgs, upload)

	file, err := os.Open(args.Filename)
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("Couldn't open file: %s, %v", args.Filename, err)
	}

	response, err := call.Media(file).Do()
	if err != nil {
		return nil, fmt.Errorf("Couldn't upload file: %v", err)
	}

	fmt.Printf("%s uploaded successfully!\nTitle: %s\n", args.Filename, args.Title)

	return response, nil
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
