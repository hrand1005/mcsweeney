package twitch

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nicklaw5/helix"
)

// Scraper defines an interface for scraping twitch clips from the web
type Scraper interface {
	Scrape(ClipFilter, chan<- helix.Clip, <-chan bool)
}

// scraper uses a client and query to scrape clips, and maintains a cursor using
// the twitch API's paging system
type scraper struct {
	client *helix.Client
	err    error
	page   helix.Pagination
	query  helix.ClipsParams
	tokenFile string
}

var ErrInvalidOptions = errors.New("NewScraper: options must contian ClientID and ClientSecret")

// NewScraper configures a scraper with the provided client options, clip query params, and file
// to retrieve and write the api app access token
func NewScraper(o *helix.Options, q helix.ClipsParams, f string) (Scraper, error) {
	if o == nil || o.ClientID == "" || o.ClientSecret == "" {
		return nil, ErrInvalidOptions
	}
	
	c, err := helix.NewClient(o)
	if err != nil {
		return nil, fmt.Errorf("NewScraper: failed to create new twitch client: %v", err)
	}

	c.SetAppAccessToken(os.Getenv(AppTokenEnvKey))

	return &scraper{
		client: c,
		query:  q,
		tokenFile: f,
	}, nil
}

// Filter interface must contain CheckClip, which takes a clip and return true if the
// clip passes the filter check, else false.
type ClipFilter func(helix.Clip) bool

// Scrape pushes clips that pass the filter to the provided clip channel, until it
// recieves a done signal. Sets error if encountered. It is the responsibility of
// the client to handle the exit conditions of the Scrape call, such as timeouts
func (s *scraper) Scrape(f ClipFilter, ch chan<- helix.Clip, done <-chan bool) {
	// continue scraping until the client is done
	for {
		s.query.After = s.page.Cursor
		cResp, err := s.client.GetClips(&s.query)
		if err != nil {
			s.err = fmt.Errorf("Encountered error scraping clips: %v\n", err)
			return
		}

		// check GetClips response, if 401 generate new token, else set error and exit
		if cResp.StatusCode != http.StatusOK {
			if cResp.StatusCode == http.StatusUnauthorized {
				log.Println("Scrape: Got 401 Status Code, generating new access token")
				err := UpdateAppToken(s.client, s.tokenFile)
				if err != nil {
					s.err = fmt.Errorf("Scrape: failed to update app token: %v", err)
					return
				}
				// if the AppToken is successfully updated, start anew and get clips
				continue
			} else {
				s.err = fmt.Errorf("Response returned status %v\nError message: %s", cResp.StatusCode, cResp.ErrorMessage)
				log.Println(s.err)
				return
			}
		}

		// filter clips and push valid responses on the clips channel
		for _, c := range cResp.Data.Clips {
			if f(c) {
				// blocks until either client is ready to recieve another clip or sends done signal
				select {
				case ch <- c:
					log.Printf("Scrape: Sent clip: %+v\n", c)
					continue
				case <-done:
					log.Println("Scrape: Recieved done signal, returning")
					return
				}
			}
		}

		// if no done signal is recieved, set the pagination and continue scraping
		select {
		case <-done:
			log.Println("Scrape: Recieved done signal, returning")
			return
		default:
			s.page.Cursor = cResp.Data.Pagination.Cursor
		}
	}
}

// Err returns the last error encountered during scraping, if any
func (s *scraper) Err() error {
	return s.err
}
