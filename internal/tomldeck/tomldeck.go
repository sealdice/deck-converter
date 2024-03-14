package tomldeck

import (
	"fmt"
	"io"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Meta struct {
	Title      string     `toml:"title"`
	Author     string     `toml:"author,omitempty"`
	Authors    []string   `toml:"authors,omitempty"`
	Version    string     `toml:"version,omitempty"`
	License    string     `toml:"license,omitempty"`
	Date       *time.Time `toml:"date,omitempty"`
	UpdateDate *time.Time `toml:"update_date,omitempty"`
	Desc       string     `toml:"desc,omitempty,multiline"`
	UpdateUrls []string   `toml:"update_urls,omitempty"`
	Etag       string     `toml:"etag,omitempty"`

	FormatVersion int64 `toml:"format_version,omitempty"`
}

type SpecialDeck struct {
	Export  bool     `toml:"export"`
	Visible bool     `toml:"visible"`
	Options []string `toml:"options"`
}

type File struct {
	Meta     Meta
	Decks    map[string][]string
	SplDecks map[string]SpecialDeck
}

func (t File) Output(out io.Writer) error {
	enc := toml.NewEncoder(out)

	err := enc.Encode(map[string]Meta{"meta": t.Meta})
	if err != nil {
		return fmt.Errorf("failed to encode TOML data: %w", err)
	}

	enc = enc.SetArraysMultiline(true)

	if len(t.Decks) > 0 {
		_, _ = out.Write([]byte("\n"))
		err = enc.Encode(map[string]map[string][]string{
			"decks": t.Decks,
		})
		if err != nil {
			return fmt.Errorf("failed to encode TOML data: %w", err)
		}
	}

	if len(t.SplDecks) > 0 {
		_, _ = out.Write([]byte("\n"))
		err = enc.Encode(t.SplDecks)
		if err != nil {
			return fmt.Errorf("failed to encode TOML data: %w", err)
		}
	}

	return nil
}
