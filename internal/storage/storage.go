package storage

import "strconv"

type Storage struct {
	data map[string]string
	cnt  int
}

func New() *Storage {
	return &Storage{
		data: map[string]string{},
		cnt:  0,
	}
}

func (s *Storage) Save(url string) string {
	id := strconv.Itoa(s.cnt)
	s.data[id] = url
	s.cnt++
	return id
}

func (s *Storage) Load(id string) (string, bool) {
	url, ok := s.data[id]
	return url, ok
}
