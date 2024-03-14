package yamldeck

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/Xiangze-Li/deckconvert/internal"
	"github.com/Xiangze-Li/deckconvert/internal/tomldeck"
)

const (
	metaTitle   = "name"
	metaAuthor  = "author"
	metaVersion = "version"
	metaCommand = "command"
	metaLicense = "license"
	metaDesc    = "desc"
	metaInfo    = "info"
	metaInclude = "include"
	metaDefault = "default"
)

type File struct {
	raw    map[string]any
	decks  map[string][]string
	export map[string]bool
}

var _ internal.DeckFile = &File{}

func (y *File) Read(r io.Reader) (internal.DeckFile, error) {
	err := yaml.NewDecoder(r).Decode(&(y.raw))
	if err != nil {
		return y, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	return y, nil
}

func (y File) Convert(logger *log.Logger) tomldeck.File {
	y.decks = make(map[string][]string)
	t := tomldeck.File{
		Decks:    make(map[string][]string),
		SplDecks: make(map[string]tomldeck.SpecialDeck),
	}

	for k, v := range y.raw {
		switch k {
		case metaTitle:
			t.Meta.Title = assertString(k, v, logger)
		case metaAuthor:
			t.Meta.Author = assertString(k, v, logger)
		case metaVersion:
			t.Meta.Version = strconv.Itoa(assertInteger(k, v, logger))
		case metaCommand:
			// handled when metaDefault is processed
		case metaLicense:
			t.Meta.License = assertString(k, v, logger)
		case metaDesc:
			t.Meta.Desc = assertString(k, v, logger)
		case metaInfo, metaInclude:
			// no-op
		case metaDefault:
			cmdAny, cmdExist := y.raw[metaCommand]
			if !cmdExist {
				logger.Printf("YAML field %q is defined but %q is missing", metaDefault, metaCommand)
				continue
			}
			cmdName := assertString(metaCommand, cmdAny, logger)
			if len(cmdName) == 0 {
				continue
			}
			cmdDeck := assertSliceStr(k, v, logger)
			if len(cmdDeck) == 0 {
				continue
			}
			y.decks[cmdName] = cmdDeck
			y.export = map[string]bool{cmdName: true}
		default:
			vv := assertSliceStr(k, v, logger)
			if len(vv) == 0 {
				continue
			}
			y.decks[k] = vv
		}
	}

	if y.export == nil {
		y.export = make(map[string]bool, len(y.decks))
		for k := range y.decks {
			y.export[k] = true
		}
	}

	for k, v := range y.decks {
		v = slices.Clone(v)
		if y.export[k] == !strings.HasPrefix(k, "__") {
			t.Decks[k] = v
		} else {
			t.SplDecks[k] = tomldeck.SpecialDeck{
				Export:  y.export[k],
				Visible: y.export[k],
				Options: v,
			}
		}
	}

	return t
}

func assertString(k string, v any, logger *log.Logger) string {
	s, ok := v.(string)
	if !ok {
		logger.Printf("YAML field %q has invalid type %T", k, v)
	}
	return s
}

func assertInteger(k string, v any, logger *log.Logger) int {
	i, ok := v.(int)
	if !ok {
		logger.Printf("YAML field %q has invalid type %T", k, v)
	}
	return i
}

func assertSliceStr(k string, v any, logger *log.Logger) []string {
	s, ok := v.([]any)
	if !ok {
		logger.Printf("YAML field %q has invalid type %T", k, v)
		return nil
	}
	ss := make([]string, 0, len(s))
	for _, vv := range s {
		s, ok := vv.(string)
		if !ok {
			logger.Printf("YAML field %q has invalid element type %T", k, vv)
			return nil
		}
		ss = append(ss, s)
	}
	return ss
}
