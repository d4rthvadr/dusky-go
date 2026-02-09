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
}

const PaginationQueryLimit = 20

func NewPaginatedFeedQuery() *PaginatedFeedQuery {
	return &PaginatedFeedQuery{
		Limit:  PaginationQueryLimit,
		Offset: 0,
		Sort:   "desc",
	}
}

// Parse extracts pagination parameters from the query string and populates the PaginatedFeedQuery struct.
func (p *PaginatedFeedQuery) Parse(r *http.Request) error {

	qs := r.URL.Query()

	limitStr := qs.Get("limit")
	if limitStr == "" {
		p.Limit = PaginationQueryLimit
	} else {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return fmt.Errorf("invalid limit parameter")
		}
		p.Limit = limit
	}

	offsetStr := qs.Get("offset")
	if offsetStr == "" {
		p.Offset = 0
	} else {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return fmt.Errorf("invalid offset parameter")
		}
		p.Offset = offset
	}

	sort := qs.Get("sort")
	if sort == "" {
		p.Sort = "desc"
	} else {
		p.Sort = sort
	}
	return nil
}
