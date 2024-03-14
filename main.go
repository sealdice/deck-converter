package main

import (
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/Xiangze-Li/deckconvert/internal/jsondeck"
)

var flag = struct {
	Output    []string `short:"o" long:"output" value-name:"OutputFile" description:"PLACEHOLDER"`
	Parent    string   `short:"p" long:"parent" value-name:"ParentDir" description:"PLACEHOLDER"`
	Overwrite bool     `short:"O" long:"overwrite" description:"Overwrite the output file if it already exists."`

	Args struct {
		Input []string `positional-arg-name:"InputFiles" required:"1"`
	} `positional-args:"yes" required:"yes"`
}{}

var helpFlags = struct {
	Version bool `short:"v" long:"version" description:"Show version" group:"Help Options"`
	Help    bool `short:"h" long:"help" description:"Show this help message" group:"Help Options"`
}{}

var (
	Version = semver.MustParse(VersionStr)
)

func initialize() {
	initLogger()
	parseFlags()
}

func main() {
	initialize()

	checkOutput()

	wg := &sync.WaitGroup{}
	wg.Add(len(flag.Args.Input))

	for i, in := range flag.Args.Input {
		go func(in, out string) {
			defer wg.Done()

			logger := getLogger(in)

			if in == out {
				logger.Printf("output file %q is the same as input file", out)
				return
			}

			i, errI := openInput(in)
			if errI != nil {
				logger.Println(errI)
				return
			}
			defer i.Close()
			o, errO := openOutput(out)
			if errO != nil {
				logger.Println(errO)
				return
			}
			defer o.Close()

			json, err := (&jsondeck.JsonFile{}).Read(i)
			if err != nil {
				logger.Println(err)
				return
			}
			err = json.Convert(logger).Output(o)
			if err != nil {
				logger.Println(err)
				return
			}

			logger.Printf("finished converting into %q", out)
		}(in, flag.Output[i])
	}

	wg.Wait()
}
