package store

import (
	"fmt"
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	Search string `json:"search" validate:"omitempty"`
	// additional filters can be added here, e.g. tags, date range, etc.
	// should use a separate struct for filters and embed it here if there are many fields to avoid bloating this struct with too many fields that are not related to pagination
	Tags []string `json:"tags" validate:"max=4"`
}

const PaginationQueryLimit = 20

func NewPaginatedFeedQuery() *PaginatedFeedQuery {
	return &PaginatedFeedQuery{
		Limit:  PaginationQueryLimit,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
	}
}

// Parse extracts pagination parameters from the query string and populates the PaginatedFeedQuery struct.
func (p *PaginatedFeedQuery) Parse(r *http.Request) error {

	qs := r.URL.Query()

	if limitStr := qs.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return fmt.Errorf("invalid limit parameter")
		}
		p.Limit = limit
	}

	if offsetStr := qs.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return fmt.Errorf("invalid offset parameter")
		}
		p.Offset = offset
	}

	p.Search = qs.Get("search")

	if sort := qs.Get("sort"); sort != "" {
		p.Sort = sort
	}
	return nil
}
