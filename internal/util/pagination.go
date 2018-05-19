package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

const DefaultLimit = 10

//Pagination represents data for retrieving time related structures.
type Pagination struct {
	limit  int
	before int
	after  int
}

//NewPagination returns a Pagination.
func NewPagination(c *gin.Context) *Pagination {
	var pagination Pagination

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = DefaultLimit
	}
	pagination.limit = limit

	before, err := strconv.Atoi(c.Query("before"))
	if err == nil {
		pagination.before = before
	}
	after, err := strconv.Atoi(c.Query("after"))
	if err == nil {
		pagination.after = after
	}

	return &pagination
}

//GetLimit returns limit.
func (p *Pagination) GetLimit() int {
	return p.limit
}

//GetBefore return before.
func (p *Pagination) GetBefore() int {
	return p.before
}

//GetAfter returns after.
func (p *Pagination) GetAfter() int {
	return p.after
}
