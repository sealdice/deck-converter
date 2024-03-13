package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Xiangze-Li/deckconvert/internal"
	"github.com/jessevdk/go-flags"
)

func parseFlags() {
	parser := flags.NewParser(&flag, flags.HelpFlag|flags.PassDoubleDash)
	parser.FindOptionByShortName('o').Description =
		"Output files. Must be empty or of same number as InputFiles.\n" +
			"If empty, <input>.toml is used for each input file.\n" +
			"If specified, each input file will be converted into corresponding output file."
	parser.FindOptionByShortName('p').Description =
		"Output Parent Directory. Ignored if -o is specified.\n" +
			"If specified, all output file will be in such directory.\n" +
			"Else, each output file will be in the same directory as the input file."

	_, err := parser.Parse()

	help := &flags.Error{}
	if errors.As(err, &help) && help.Type == flags.ErrHelp {
		fmt.Println(help.Message)
		os.Exit(0)
	}

	if flag.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if err != nil {
		log.Fatalln(err)
	}

	if len(flag.Parent) > 0 {
		err = os.MkdirAll(flag.Parent, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func checkOutput() {
	if len(flag.Output) == 0 {
		flag.Output = make([]string, 0, len(flag.Args.Input))
		for _, in := range flag.Args.Input {
			dir, file := filepath.Split(in)
			ext := filepath.Ext(file)
			base := strings.TrimSuffix(file, ext)
			if len(flag.Parent) > 0 {
				dir = flag.Parent
			}
			flag.Output = append(flag.Output, filepath.Join(dir, base+".toml"))
		}
	}

	if len(flag.Output) != len(flag.Args.Input) {
		log.Fatalf(
			"number of output files (%d) is diffrent from number of input files (%d)",
			len(flag.Output), len(flag.Args.Input),
		)
	}

	abs := make([]string, 0, len(flag.Args.Input))
	for _, in := range flag.Args.Input {
		absFn, _ := filepath.Abs(in)
		abs = append(abs, absFn)
	}
	absVis := internal.VisMap(abs)
	if len(absVis) != len(abs) {
		log.Fatalf("detected duplication in input files")
	}
}