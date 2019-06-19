package main

// Message to post to a Jandi topick.
// See also: http://bit.ly/31Bxeeh
type Message struct {
	MessageBody string       `json:"body,omitempty"`
	Color       string       `json:"connectColor"`
	Attachments []Attachment `json:"connectInfo"`
}

// Attachment adds more context to a Jandi chat message.
// See also: http://bit.ly/31Bxeeh
type Attachment struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

// Constants for webhook
const (
	headerAccept      = "application/vnd.tosslab.jandi-v2+json"
	headerContentType = "application/json"
)
