package internal

import (
	"io"
	"log"

	"github.com/Xiangze-Li/deckconvert/internal/tomldeck"
)

func VisMap[K comparable](s []K) map[K]bool {
	m := make(map[K]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

type DeckFile interface {
	Read(r io.Reader) (DeckFile, error)
	Convert(logger *log.Logger) tomldeck.File
}
