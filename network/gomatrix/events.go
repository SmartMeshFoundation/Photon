package gomatrix

import (
	"html"
	"regexp"
)

// Event represents a single Matrix event.
type Event struct {
	StateKey  *string                `json:"state_key,omitempty"`
	Sender    string                 `json:"sender"`
	Type      string                 `json:"type"`
	Timestamp int64                  `json:"origin_server_ts"`
	ID        string                 `json:"event_id"`
	RoomID    string                 `json:"room_id"`
	Content   map[string]interface{} `json:"content"`
	Redacts   string                 `json:"redacts,omitempty"`
}

// ViewContent returns the value of the "msgtype" key in the event.
func (event *Event) ViewContent(tip string) (body string, ok bool) {
	value, exists := event.Content[tip]
	if !exists {
		return
	}
	body, ok = value.(string)
	return
}

// Body returns the value of the "body" key in the event content if it is
// present and is a string.
func (event *Event) Body() (body string, ok bool) {
	value, exists := event.Content["body"]
	if !exists {
		return
	}
	body, ok = value.(string)
	return
}

// MessageType returns the value of the "msgtype" key in the event content if
// it is present and is a string.
func (event *Event) MessageType() (msgtype string, ok bool) {
	value, exists := event.Content["msgtype"]
	if !exists {
		return
	}
	msgtype, ok = value.(string)
	return
}

// TextMessage is the contents of a Matrix formated message event.
type TextMessage struct {
	MsgType string `json:"msgtype"`
	Body    string `json:"body"`
}

// ImageInfo contains info about an image
type ImageInfo struct {
	Height   uint   `json:"h,omitempty"`
	Width    uint   `json:"w,omitempty"`
	Mimetype string `json:"mimetype,omitempty"`
	Size     uint   `json:"size,omitempty"`
}

// VideoInfo contains info about a video
type VideoInfo struct {
	Mimetype      string    `json:"mimetype,omitempty"`
	ThumbnailInfo ImageInfo `json:"thumbnail_info"`
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`
	Height        uint      `json:"h,omitempty"`
	Width         uint      `json:"w,omitempty"`
	Duration      uint      `json:"duration,omitempty"`
	Size          uint      `json:"size,omitempty"`
}

// VideoMessage is an m.video
type VideoMessage struct {
	MsgType string    `json:"msgtype"`
	Body    string    `json:"body"`
	URL     string    `json:"url"`
	Info    VideoInfo `json:"info"`
}

// ImageMessage is an m.image event
type ImageMessage struct {
	MsgType string    `json:"msgtype"`
	Body    string    `json:"body"`
	URL     string    `json:"url"`
	Info    ImageInfo `json:"info"`
}

// HTMLMessage is the contents of a Matrix HTML formated message event.
type HTMLMessage struct {
	Body          string `json:"body"`
	MsgType       string `json:"msgtype"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
}

var htmlRegex = regexp.MustCompile("<[^<]+?>")

// GetHTMLMessage returns an HTMLMessage with the body set to a stripped version of the provided HTML, in addition
// to the provided HTML.
func GetHTMLMessage(msgtype, htmlText string) HTMLMessage {
	return HTMLMessage{
		Body:          html.UnescapeString(htmlRegex.ReplaceAllLiteralString(htmlText, "")),
		MsgType:       msgtype,
		Format:        "org.matrix.custom.html",
		FormattedBody: htmlText,
	}
}
