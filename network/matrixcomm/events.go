package matrixcomm

import (
	"html"
	"regexp"
)

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

func (event *Event) ViewContent(tip string) (body string, ok bool) {
	value, exists := event.Content[tip]
	if !exists {
		return
	}
	body, ok = value.(string)
	return
}

func (event *Event) Body() (body string, ok bool) {
	value, exists := event.Content["body"]
	if !exists {
		return
	}
	body, ok = value.(string)
	return
}

func (event *Event) MessageType() (msgtype string, ok bool) {
	value, exists := event.Content["msgtype"]
	if !exists {
		return
	}
	msgtype, ok = value.(string)
	return
}

type TextMessage struct {
	MsgType string `json:"msgtype"`
	Body    string `json:"body"`
}

type ImageInfo struct {
	Height   uint   `json:"h,omitempty"`
	Width    uint   `json:"w,omitempty"`
	Mimetype string `json:"mimetype,omitempty"`
	Size     uint   `json:"size,omitempty"`
}

type VideoInfo struct {
	Mimetype      string    `json:"mimetype,omitempty"`
	ThumbnailInfo ImageInfo `json:"thumbnail_info"`
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`
	Height        uint      `json:"h,omitempty"`
	Width         uint      `json:"w,omitempty"`
	Duration      uint      `json:"duration,omitempty"`
	Size          uint      `json:"size,omitempty"`
}

type VideoMessage struct {
	MsgType string    `json:"msgtype"`
	Body    string    `json:"body"`
	URL     string    `json:"url"`
	Info    VideoInfo `json:"info"`
}

type ImageMessage struct {
	MsgType string    `json:"msgtype"`
	Body    string    `json:"body"`
	URL     string    `json:"url"`
	Info    ImageInfo `json:"info"`
}

type HTMLMessage struct {
	Body          string `json:"body"`
	MsgType       string `json:"msgtype"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
}

var htmlRegex = regexp.MustCompile("<[^<]+?>")

func GetHTMLMessage(msgtype, htmlText string) HTMLMessage {
	return HTMLMessage{
		Body:          html.UnescapeString(htmlRegex.ReplaceAllLiteralString(htmlText, "")),
		MsgType:       msgtype,
		Format:        "org.matrix.custom.html",
		FormattedBody: htmlText,
	}
}
