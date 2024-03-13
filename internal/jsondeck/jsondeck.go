package jsondeck

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"

	"github.com/golang-module/carbon"
	"github.com/tailscale/hujson"

	"github.com/Xiangze-Li/deckconvert/internal"
	"github.com/Xiangze-Li/deckconvert/internal/tomldeck"
)

const (
	jsonMetaTitle      = "_title"
	jsonMetaAuthor     = "_author"
	jsonMetaDate       = "_date"
	jsonMetaUpdateDate = "_updateDate"
	jsonMetaBrief      = "_brief"
	jsonMetaVersion    = "_version"
	jsonMetaLicense    = "_license"
	jsonMetaUpdateUrls = "_updateUrls"
	jsonMetaEtag       = "_etag"
	jsonMetaKeys       = "_keys"
	jsonMetaExport     = "_export"
	jsonMetaExports    = "_exports"
	jsonMetaSchema     = "$schema"
)

type JsonFile struct {
	raw     map[string]any
	decks   map[string][]string
	visible map[string]bool
	export  map[string]bool
}

func (j *JsonFile) Read(r io.Reader) (*JsonFile, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return j, fmt.Errorf("failed to read JSON data: %w", err)
	}
	huJ, err := hujson.Parse(b)
	if err != nil {
		return j, fmt.Errorf("failed to parse JSON data: %w", err)
	}
	huJ.Standardize()

	err = json.Unmarshal(huJ.Pack(), &(j.raw))
	if err != nil {
		return j, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}
	return j, nil
}

func (j JsonFile) Convert(logger *log.Logger) tomldeck.TomlFile {
	j.decks = make(map[string][]string)
	t := tomldeck.TomlFile{
		Decks:    make(map[string][]string),
		SplDecks: make(map[string]tomldeck.TomlSpecialDeck),
	}

	for k, v := range j.raw {
		if k == jsonMetaSchema {
			continue
		}
		s, ok := assertSliceStr(v)
		if !ok {
			logger.Printf("JSON field %q has invalid type", k)
			continue
		}
		switch k {
		case jsonMetaTitle:
			t.Meta.Title = strings.Join(s, " / ")
		case jsonMetaAuthor:
			switch len(s) {
			case 0: // no-op
			case 1:
				t.Meta.Author = s[0]
			default:
				t.Meta.Authors = slices.Clone(s)
			}
		case jsonMetaDate:
			date := strings.Join(s, "/")
			parsed := carbon.Parse(date)
			if parsed.Error != nil {
				logger.Printf("JSON field %q is not a valid datetime: %v", k, parsed.Error)
			}
			time := parsed.ToStdTime()
			t.Meta.Date = &time
		case jsonMetaUpdateDate:
			date := strings.Join(s, "/")
			parsed := carbon.Parse(date)
			if parsed.Error != nil {
				logger.Printf("JSON meta field %q is not a valid datetime: %v", k, parsed.Error)
			}
			time := parsed.ToStdTime()
			t.Meta.UpdateDate = &time
		case jsonMetaBrief:
			t.Meta.Desc = strings.Join(s, "\n")
		case jsonMetaVersion:
			t.Meta.Version = strings.Join(s, " / ")
		case jsonMetaLicense:
			t.Meta.License = strings.Join(s, " / ")
		case jsonMetaUpdateUrls:
			t.Meta.UpdateUrls = slices.Clone(s)
		case jsonMetaEtag:
			if len(s) > 0 {
				t.Meta.Etag = s[0]
			}
		case jsonMetaKeys:
			j.visible = internal.VisMap(s)
		case jsonMetaExport, jsonMetaExports:
			j.export = internal.VisMap(s)
		default:
			j.decks[k] = slices.Clone(s)
		}
	}

	if j.export == nil {
		j.export = make(map[string]bool, len(j.decks))
		for k := range j.decks {
			j.export[k] = true
		}
	}
	if j.visible == nil {
		j.visible = make(map[string]bool, len(j.decks))
		for k := range j.decks {
			j.visible[k] = !strings.HasPrefix(k, "_")
		}
	}

	for k, v := range j.decks {
		if !j.export[k] {
			if strings.HasPrefix(k, "__") {
				t.Decks[k] = v
			} else {
				t.SplDecks[k] = tomldeck.TomlSpecialDeck{Options: v}
			}
			continue
		}
		if !j.visible[k] {
			if strings.HasPrefix(k, "_") {
				t.Decks[k] = v
			} else {
				t.SplDecks[k] = tomldeck.TomlSpecialDeck{Options: v, Export: true}
			}
			continue
		}
		if strings.HasPrefix(k, "_") {
			t.SplDecks[k] = tomldeck.TomlSpecialDeck{Options: v, Export: true, Visible: true}
		} else {
			t.Decks[k] = v
		}
	}

	return t
}

func assertSliceStr(v any) ([]string, bool) {
	v1, ok := v.([]any)
	if !ok {
		return nil, false
	}
	s := make([]string, 0, len(v1))
	for _, v2 := range v1 {
		v3, ok := v2.(string)
		if !ok {
			return nil, false
		}
		s = append(s, v3)
	}
	return s, true
}
