package content

type ContentStatus int

const (
    UNKNOWN ContentStatus = 0
    RAW ContentStatus = 1
    PROCESSED ContentStatus = 2
)

type ContentObj struct {
    CreatorName string
    Title string
    Description string
    Path string
    Status ContentStatus
    Url string
}
