# go-opml

[![Build Status](http://img.shields.io/travis/kaorimatz/go-opml.svg?style=flat-square)](https://travis-ci.org/kaorimatz/go-opml)

Go library for parsing and rendering OPML

## Usage

```go
// Parsing
reader, _ := os.Open("input.opml")
inputOPML, _ := opml.Parse(reader)

// Rendering
writer, _ := os.Create("output.opml")
opml.Render(writer, inputOPML)
```

## License

MIT
