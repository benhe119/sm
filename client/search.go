package client

import (
	"fmt"

	"github.com/prologic/sm"
)

// SearchFilter ...
type SearchFilter struct {
	ID    string
	Name  string
	State string
}

// SearchOptions ...
type SearchOptions struct {
	Filter *SearchFilter
}

// Search ...
func (c *Client) Search(options *SearchOptions) (res []*sm.Event, err error) {
	url := fmt.Sprintf("%s/search", c.url)

	filter := options.Filter

	switch {
	case filter.ID != "":
		url += fmt.Sprintf("/%s", filter.ID)
	case filter.Name != "":
		url += fmt.Sprintf("?q=name:%s", filter.Name)
	case filter.State != "":
		url += fmt.Sprintf("?q=state:%d", sm.ParseState(filter.State))
	}

	return c.request("GET", url, nil)
}
