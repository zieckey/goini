package wini

import (
	"bytes"
	"errors"
	"io/ioutil"
	//"fmt"
	"log"
)

// Suppress error if they are not otherwise used.
var _ = log.Printf

type Kvmap map[string]string
type SectionMap map[string]Kvmap

type INI struct {
	sections SectionMap
	linesep  string
	kvsep    string
}

func New() *INI {
	ini := &INI{
		sections: make(SectionMap),
		linesep:  "\n",
		kvsep:    "=",
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

	return ini.parseINI(contents, "\n", "=")
}

// Parse parse the data to store the data in the INI
// A successful call returns err == nil
func (ini *INI) Parse(data []byte, linesep, kvsep string) error {
	return ini.parseINI(data, linesep, kvsep)
}

// Get looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (i INI) Get(section, key string) (value string, ok bool) {
	if s := i.sections[section]; s != nil {
		value, ok = (s)[key]
	}
	return
}

func (i INI) GetKvmap(section string) (kvmap Kvmap, ok bool) {
	kvmap, ok = i.sections[section]
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
			// Try to parse INI-Section
			section = string(line[1 : size-1])
			//log.Printf("Got a sction [%v]\n", section)
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
		//log.Printf("Got a key/value pair [%v/%v] for section [%v]\n", string(k), string(v), section)
	}
	return nil
}
