## wini

This is a Go package to interact with arbitrary INI

[![Build Status](https://secure.travis-ci.org/zieckey/wini.png)](http://travis-ci.org/zieckey/wini)

## Importing

    import github.com/zieckey/wini

## Usage

### Example 1 : Parsing an INI file

The simplest example code is :
```go
import github.com/zieckey/wini

ini := wini.New()
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
ini := wini.New()
err := ini.Parse(raw, "|", ":")
if err != nil {
	fmt.Printf("parse INI memory data failed : %v\n", err.Error())
	return
}

v, ok := ini.Get("a")
//...
```