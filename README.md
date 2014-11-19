## goini

This is a Go package to interact with arbitrary INI

[![Build Status](https://secure.travis-ci.org/zieckey/goini.png)](http://travis-ci.org/zieckey/goini)

## Importing

    import github.com/zieckey/goini

## Usage

### Example 1 : Parsing an INI file

The simplest example code is :
```go
import github.com/zieckey/goini

ini := goini.New()
err := ini.ParseFile(filename)
if err != nil {
	fmt.Printf("parse INI file %v failed : %v\n", filename, err.Error())
	return
}

v, ok := ini.Get("the-key")
//...
```

### Example 2 : Parsing the memory data like the format of INI

```go
raw := []byte("a:av|b:bv||c:cv|||d:dv||||||")
ini := goini.New()
err := ini.Parse(raw, "|", ":")
if err != nil {
	fmt.Printf("parse INI memory data failed : %v\n", err.Error())
	return
}

v, ok := ini.Get("a")
//...
```