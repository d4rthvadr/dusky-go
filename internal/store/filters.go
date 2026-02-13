package store

import (
	"net/http"
	"strings"

	"github.com/d4rthvadr/dusky-go/internal/utils"
)

type FilterFeedQuery struct {
	Tags   []string `json:"tags" validate:"max=4"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func NewFilterFeedQuery() *FilterFeedQuery {
	return &FilterFeedQuery{
		Tags:   []string{},
		Search: "",
		Since:  "",
		Until:  "",
	}
}

// Parse extracts filter parameters from the query string and populates the FilterFeedQuery struct.
func (f *FilterFeedQuery) Parse(r *http.Request) {
	qs := r.URL.Query()

	tags := qs.Get("tags")

	if tags != "" {
		f.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")

	limit := qs.Get("limit")
	if limit != "" {
		f.Search = search
	}

	if search != "" {
		f.Search = search
	}
	since := qs.Get("since")
	if since != "" {
		parsedSince, err := utils.ParseStrToTime(since)
		if err == nil {
			f.Since = parsedSince
		}
	}
	until := qs.Get("until")
	if until != "" {
		parsedUntil, err := utils.ParseStrToTime(until)
		if err == nil {
			f.Until = parsedUntil
		}
	}
}
