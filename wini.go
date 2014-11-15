package wini

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
)

// Suppress error if they are not otherwise used.
var _ = log.Printf

type Kvmap map[string]string
type SectionMap map[string]Kvmap

const (
	DefaultSection           = ""
	DefaultLineSeperator     = "\n"
	DefaultKeyValueSeperator = "="
)

type INI struct {
	sections SectionMap
	linesep  string
	kvsep    string
}

func New() *INI {
	ini := &INI{
		sections: make(SectionMap),
		linesep:  DefaultLineSeperator,
		kvsep:    DefaultKeyValueSeperator,
	}
	return ini
}

// ParseFile reads the INI file named by filename and parse the contents to store the data in the INI
// A successful call returns err == nil
func (ini *INI) ParseFile(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return ini.parseINI(contents, DefaultLineSeperator, DefaultKeyValueSeperator)
}

// Parse parse the data to store the data in the INI
// A successful call returns err == nil
func (ini *INI) Parse(data []byte, linesep, kvsep string) error {
	return ini.parseINI(data, linesep, kvsep)
}

// Get looks up a value for a key in the default section
// and returns that value, along with a boolean result similar to a map lookup.
func (ini *INI) Get(key string) (value string, ok bool) {

	return ini.SectionGet(DefaultSection, key)
}

// Get looks up a value for a key in a section
// and returns that value, along with a boolean result similar to a map lookup.
func (ini *INI) SectionGet(section, key string) (value string, ok bool) {
	if s := ini.sections[section]; s != nil {
		value, ok = s[key]
	}
	return
}

func (ini *INI) GetKvmap(section string) (kvmap Kvmap, ok bool) {
	kvmap, ok = ini.sections[section]
	return kvmap, ok
}

//////////////////////////////////////////////////////////////////////////

func (ini *INI) parseINI(data []byte, linesep, kvsep string) error {
	ini.linesep = linesep
	ini.kvsep = kvsep

	// Insert the default section
	var section string
	kvmap := make(Kvmap)
	ini.sections[section] = kvmap

	lines := bytes.Split(data, []byte(linesep))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		size := len(line)
		if size == 0 {
			// Skip blank lines
			continue
		}
		if line[0] == ';' || line[0] == '#' {
			// Skip comments
			continue
		}
		if line[0] == '[' && line[size-1] == ']' {
			// Parse INI-Section
			section = string(line[1 : size-1])
			kvmap = make(Kvmap)
			ini.sections[section] = kvmap
			continue
		}

		pos := bytes.Index(line, []byte(kvsep))
		if pos < 0 {
			// ERROR happened when passing
			err := errors.New("Came accross an error : " + string(line) + " is NOT a valid key/value pair")
			return err
		}

		k := bytes.TrimSpace(line[0:pos])
		v := bytes.TrimSpace(line[pos+len(kvsep):])
		kvmap[string(k)] = string(v)
	}
	return nil
}
