package client

import (
	"fmt"
	"net/url"

	"github.com/prologic/sm"
)

// Create ...
func (c *Client) Create(title string, level int) (res []*sm.Event, err error) {

	url := fmt.Sprintf(
		"%s/create?title=%s&level=%d",
		c.url, url.QueryEscape(title), level,
	)

	return c.request("POST", url, nil)
}
