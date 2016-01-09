package gitter
import "time"

// A Room in Gitter can represent a GitHub Organization, a GitHub Repository, a Gitter Channel or a One-to-one conversation.
// In the case of the Organizations and Repositories, the access control policies are inherited from GitHub.
type Room struct {

	// Room ID
	Id             string `json:"id"`

	// Room name
	Name           string `json:"name"`

	// Room topic. (default: GitHub repo description)
	Topic          string `json:"topic"`

	// Room URI on Gitter
	URI            string `json:"uri"`

	// Indicates if the room is a one-to-one chat
	OneToOne       bool `json:"oneToOne"`

	// Count of users in the room
	UserCount      int `json:"userCount"`

	// Number of unread messages for the current user
	UnreadItems    int `json:"unreadItems"`

	// Number of unread mentions for the current user
	Mentions       int `json:"mentions"`

	// Last time the current user accessed the room in ISO format
	LastAccessTime time.Time `json:"lastAccessTime"`

	// Indicates if the current user has disabled notifications
	Lurk           bool `json:"lurk"`

	// Path to the room on gitter
	Url            string `json:"url"`

	// Type of the room
	// - ORG: A room that represents a GitHub Organization.
	// - REPO: A room that represents a GitHub Repository.
	// - ONETOONE: A one-to-one chat.
	// - ORG_CHANNEL: A Gitter channel nested under a GitHub Organization.
	// - REPO_CHANNEL A Gitter channel nested under a GitHub Repository.
	// - USER_CHANNEL A Gitter channel nested under a GitHub User.
	GithubType     string `json:"githubType"`

	// Tags that define the room
	Tags           []string `json:"tags"`

	RoomMember     bool `json:"roomMember"`

	// Room version.
	Version        int `json:"v"`
}