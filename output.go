package main

type resultSet struct {
	searches []*search
}

func (rs *resultSet) collect(s *search) {
	s.extract()
	rs.searches = append(rs.searches, s)
}

type output interface {
	collect(*search)
	print()
}
