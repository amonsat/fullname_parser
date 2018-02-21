# Fullname parser

## Description

ParseFullname() is designed to parse large batches of full names in multiple inconsistent formats, as from a database.

parseFullName():

1. accepts a string containing a person's full name, in any format,
2. analyzes and attempts to detect the format of that name,
3. (if possible) parses the name into its component parts, and
4. returns an object containing all individual parts of the name:
    - title (string): title(s) (e.g. "Ms." or "Dr.")
    - first (string): first name or initial
    - middle (string): middle name(s) or initial(s)
    - last (string): last name or initial
    - nick (string): nickname(s)
    - suffix (string): suffix(es) (e.g. "Jr.", "II", or "Esq.")

## Use

### Instalation
```
	go get github.com/amonsat/fullname_parser
```

### Basic Use

```go
package main
import fp "github.com/amonsat/fullname_parser"

func main() {
    parsedFullname := fp.ParseFullname("Mr. David Davis")
    println(parsedFullname.Title)
    println(parsedFullname.First)
    println(parsedFullname.Last)
}
```

### Parsedname struct

```go
type ParsedName struct {
	Title  string
	First  string
	Middle string
	Last   string
	Nick   string
	Suffix string
}
```

## Credits and precursors

My thanks to David Schnell-Davis for sharing his work.

David Schnell-Davis's parse-full-name
https://github.com/dschnelldavis/parse-full-name