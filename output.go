package main

import (
	"io"
	"time"
)

type resultSet struct {
	searches  []*search
	createdAt time.Time
}

func (rs *resultSet) collect(s *search) {
	if len(rs.searches) == 0 {
		rs.createdAt = time.Now()
	}
	s.extract()
	rs.searches = append(rs.searches, s)
}

type output interface {
	collect(*search)
	print(...io.Writer)
}
