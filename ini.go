// Copyright 2014 zieckey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goini

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
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
	sections     SectionMap
	linesep      string
	kvsep        string
	parseSection bool
	skipCommits  bool
}

func New() *INI {
	ini := &INI{
		sections:     make(SectionMap),
		linesep:      DefaultLineSeperator,
		kvsep:        DefaultKeyValueSeperator,
		parseSection: false,
		skipCommits:  false,
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
	ini.parseSection = true
	ini.skipCommits = true
	return ini.parseINI(contents, DefaultLineSeperator, DefaultKeyValueSeperator)
}

// Parse parse the data to store the data in the INI
// A successful call returns err == nil
func (ini *INI) Parse(data []byte, linesep, kvsep string) error {
	return ini.parseINI(data, linesep, kvsep)
}

// Reset clear all the data hold by INI
func (ini *INI) Reset() {
	ini.sections = make(SectionMap)
	//FIXME effective optimize
}

// SetSkipCommits set INI.skipCommits whether skip commits when parsing
func (ini *INI) SetSkipCommits(skipCommits bool) {
	ini.skipCommits = skipCommits
}

// SetParseSection set INI.parseSection whether process the INI section when parsing
func (ini *INI) SetParseSection(parseSection bool) {
	ini.parseSection = parseSection
}

// Get looks up a value for a key in the default section
// and returns that value, along with a boolean result similar to a map lookup.
func (ini *INI) Get(key string) (value string, ok bool) {
	return ini.SectionGet(DefaultSection, key)
}

// GetInt get value as int
func (ini *INI) GetInt(key string) (value int, ok bool) {
	return ini.SectionGetInt(DefaultSection, key)
}

// GetFloat get value as float64
func (ini *INI) GetFloat(key string) (value float64, ok bool) {
	return ini.SectionGetFloat(DefaultSection, key)
}

// GetBool returns the boolean value represented by the string.
// It accepts "1", "t", "T", "true", "TRUE", "True", "on", "ON", "On", "yes", "YES", "Yes" as true
// and "0", "f", "F", "false", "FALSE", "False", "off", "OFF", "Off", "no", "NO", "No" as false
// Any other value returns false.
func (ini *INI) GetBool(key string) (value bool, ok bool) {
	return ini.SectionGetBool(DefaultSection, key)
}

// Get looks up a value for a key in a section
// and returns that value, along with a boolean result similar to a map lookup.
func (ini *INI) SectionGet(section, key string) (value string, ok bool) {
	if s := ini.sections[section]; s != nil {
		value, ok = s[key]
	}
	return
}

// SectionGetInt get value as int
func (ini *INI) SectionGetInt(section, key string) (value int, ok bool) {
	v, ok := ini.SectionGet(section, key)
	if ok {
		v, err := strconv.Atoi(v)
		if err == nil {
			return v, true
		}
	}

	return 0, ok
}

// SectionGetFloat get value as float64
func (ini *INI) SectionGetFloat(section, key string) (value float64, ok bool) {
	v, ok := ini.SectionGet(section, key)
	if ok {
		v, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return v, true
		}
	}

	return 0.0, ok
}

// SectionGetBool get a value as bool. See GetBool for more detail
func (ini *INI) SectionGetBool(section, key string) (value bool, ok bool) {
	v, ok := ini.SectionGet(section, key)
	if ok {
		switch v {
		case "1", "t", "T", "true", "TRUE", "True", "on", "ON", "On", "yes", "YES", "Yes":
			return true, true
		case "0", "f", "F", "false", "FALSE", "False", "off", "OFF", "Off", "no", "NO", "No":
			return false, true
		}
	}

	return false, false
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
		if ini.skipCommits && line[0] == ';' || line[0] == '#' {
			// Skip comments
			continue
		}
		if ini.parseSection && line[0] == '[' && line[size-1] == ']' {
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
