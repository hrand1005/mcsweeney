package twitch

import(
  "fmt"
  "log"
  "net/http"

  "github.com/nicklaw5/helix"
)

// scraper uses a client and query to scrape clips, and maintains a cursor using 
// the twitch API's paging system
type scraper struct {
  client *helix.Client
  page *helix.Pagination
  query *helix.ClipsParams
}

// NewScraper configures a scraper with the provided client and ClipsParams
func NewScraper(c *helix.Client, q *helix.ClipsParams) *scraper {
  return &scraper{
    client: c,
    query: q,
  }
}

// Filter interface must contain CheckClip, which takes a clip and return true if the 
// clip passes the filter check, else false
type ClipFilter func(helix.Clip) bool

// Scrape pushes clips that pass the filter to the provided clip channel, returns error
// if encountered
func (s *scraper) Scrape(f ClipFilter, ch chan<- helix.Clip, done <-chan bool) error {
  // continue scraping until the client is done
  for <-done {
    s.query.After = s.page.Cursor
    cResp, err := s.client.GetClips(s.query)
    if err != nil {
      log.Printf("Encountered error scraping clips: %v\n", err)
      return err
    }

    if cResp.StatusCode != http.StatusOK {
      log.Printf("Response returned status %v\nError message: %s", cResp.StatusCode, cResp.ErrorMessage)
      return fmt.Errorf(cResp.ErrorMessage)
    }

    for _, c := range cResp.Data.Clips {
      if f(c) {
        ch <- c
      }
    }
    s.page.Cursor = cResp.Data.Pagination.Cursor
  }
  // TODO: should the scraper close the channel?
  //close(ch)

  return nil
}
