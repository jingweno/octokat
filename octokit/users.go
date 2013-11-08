package octokit

import (
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/url"
	"time"
)

var (
	CurrentUserURL = Hyperlink("/user")
	UserURL        = Hyperlink("/users{/user}{?since}")
)

// Create a UsersService with the base Hyperlink and the params M to expand the Hyperlink
// If no Hyperlink is passed in, it will use CurrentUserHyperlink.
func (c *Client) Users(link *Hyperlink, m M) (users *UsersService, err error) {
	if link == nil {
		link = &CurrentUserURL
	}

	url, err := link.Expand(m)
	if err != nil {
		return
	}

	users = &UsersService{client: c, URL: url}
	return
}

// A service to return user records
type UsersService struct {
	client *Client
	URL    *url.URL
}

// Get a user based on UserService#URL
func (u *UsersService) Get() (user *User, result *Result) {
	result = u.client.Get(u.URL, &user)
	return
}

// Update a user based on UserService#URL
func (u *UsersService) Update(params interface{}) (user *User, result *Result) {
	result = u.client.Put(u.URL, params, &user)
	return
}

// Get a list of users based on UserService#URL
func (u *UsersService) GetAll() (users []User, result *Result) {
	result = u.client.Get(u.URL, &users)
	return
}

type User struct {
	*hypermedia.HALResource

	Login             string     `json:"login,omitempty"`
	ID                int        `json:"id,omitempty"`
	AvatarURL         string     `json:"avatar_url,omitempty"`
	GravatarID        string     `json:"gravatar_id,omitempty"`
	URL               string     `json:"url,omitempty"`
	Name              string     `json:"name,omitempty"`
	Company           string     `json:"company,omitempty"`
	Blog              string     `json:"blog,omitempty"`
	Location          string     `json:"location,omitempty"`
	Email             string     `json:"email,omitempty"`
	Hireable          bool       `json:"hireable,omitempty"`
	Bio               string     `json:"bio,omitempty"`
	PublicRepos       int        `json:"public_repos,omitempty"`
	PublicGists       int        `json:"public_gists,omitempty"`
	Followers         int        `json:"followers,omitempty"`
	Following         int        `json:"following,omitempty"`
	HTMLURL           string     `json:"html_url,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	Type              string     `json:"type,omitempty"`
	FollowingURL      Hyperlink  `json:"following_url,omitempty"`
	FollowersURL      Hyperlink  `json:"followers_url,omitempty"`
	GistsURL          Hyperlink  `json:"gists_url,omitempty"`
	StarredURL        Hyperlink  `json:"starred_url,omitempty"`
	SubscriptionsURL  Hyperlink  `json:"subscriptions_url,omitempty"`
	OrganizationsURL  Hyperlink  `json:"organizations_url,omitempty"`
	ReposURL          Hyperlink  `json:"repos_url,omitempty"`
	EventsURL         Hyperlink  `json:"events_url,omitempty"`
	ReceivedEventsURL Hyperlink  `json:"received_events_url,omitempty"`
}
