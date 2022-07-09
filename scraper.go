package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/nicklaw5/helix"
)

var ErrInvalidClient = errors.New("NewClipScraper: Client must not be nil")

// scraper uses a client and query to scrape clips, and maintains a cursor using
// the twitch API's paging system
type clipScraper struct {
	client *helix.Client
	err    error
	page   helix.Pagination
	query  helix.ClipsParams
}

// NewClipScraper configures a scraper with the provided client options, clip query params, and file
// to retrieve and write the api app access token
func NewClipScraper(c *helix.Client, q helix.ClipsParams) (*clipScraper, error) {
	if c == nil {
		return nil, ErrInvalidClient
	}

	return &clipScraper{
		client: c,
		query:  q,
	}, nil
}

// Filter interface must contain CheckClip, which takes a clip and return true if the
// clip passes the filter check, else false.
type ClipFilter func(helix.Clip) bool

// Scrape returns a clips channel, which it pushes clips to until it recieves a
// done signal. Sets error if encountered. It is the responsibility of the client
// to handle the exit conditions of the Scrape call, such as timeouts.
func (s *clipScraper) Scrape(f ClipFilter, done <-chan bool) <-chan helix.Clip {
	clipChan := make(chan helix.Clip)
	// continue scraping until the client is done
	go func() {
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
					err := setTwitchToken(s.client)
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
			filteredClips := filterClips(f, cResp.Data.Clips)

			// if no done signal is recieved, set the pagination and continue scraping
			for _, c := range filteredClips {
				select {
				case clipChan <- c:
					log.Printf("Scrape: Sent clip with URL: %v", c.URL)
				case <-done:
					log.Println("Scrape: Recieved done signal, returning...")
					return
				}
			}

			select {
			case <-done:
				log.Println("Scrape: Recieved done signal, returning...")
				return
			default:
				s.page.Cursor = cResp.Data.Pagination.Cursor
			}
		}
	}()

	return clipChan
}

// filterClips returns the subset of clips that pass the filter
func filterClips(f ClipFilter, clips []helix.Clip) []helix.Clip {
	filtered := make([]helix.Clip, 0, len(clips))
	for _, c := range clips {
		if f(c) {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

// Err returns the last error encountered during scraping, if any
func (s *clipScraper) Err() error {
	return s.err
}
