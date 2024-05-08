# Deck Converter

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sealdice/deck-converter)
[![GoDoc](https://godoc.org/github.com/sealdice/deck-converter?status.svg)](https://pkg.go.dev/mod/github.com/sealdice/deck-converter)
![GitHub License](https://img.shields.io/github/license/sealdice/deck-converter)
![GitHub Release](https://img.shields.io/github/v/release/sealdice/deck-converter)
![GitHub Tag](https://img.shields.io/github/v/tag/sealdice/deck-converter)

A library / CLI tool to convert deck of other formats to SealDice TOML format.

## Install as a CLI tool

```shell
go install github.com/sealdice/deck-converter/cli/deck-converter@latest
```

### Usage

#### Convert deck file(s)

Feed path(s) of input file(s) to the tool.

```shell
deck-converter path/to/input.json path/to/another/input.yaml
```

This converts each input file to a toml file in the same directory, or `path/to/input.toml` and `path/to/another/input.toml`.

#### Specify output file(s)

Use `-o` flag to sepcify output file for each input. This flag must occur 0 times or the same time as number of input files.

Specified output files are used one-by-one. The first input will be converted to the first output, and so on.

```shell
deck-converter path/to/input.json -o output/from_json.toml path/to/another/input.yaml -o output/from_yaml.toml
```

Output files will be `output/from_json.toml` and `output/from_yaml.toml`

#### Specify output directory

Use `-p` flag to command all output files to be put under that directory.
Each output will have the same file name as input, with extension name changed to `toml`.

This flag overrides all `-o` flags.

```shell
deck-converter -p output/ path/to/input_json.json path/to/another/input_yaml.yaml
```

Output files will be `output/input_json.toml` and `output/input_yaml.toml`

#### Allow overwrites existing file(s)

Use `-O` flag to enable overwriting exist files.

`deck-converter input.json -o output.toml` fails if `output.toml` exist.

If it's wished to overwrites `output.toml`, use:

```shell
deck-converter -O input.json -o output.toml
```

## Use as Go package

```go
import "github.com/sealdice/deck-converter"
```
