// Copyright 2014 zieckey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goini

import (
	"bufio"
	"bytes"
	"errors"
	"io"
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
func (ini *INI) Parse(data []byte, lineSep, kvSep string) error {
	return ini.parseINI(data, lineSep, kvSep)
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

// SectionGet looks up a value for a key in a section
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

// GetKvmap gets all keys under section as a Kvmap (map[string]string).
// The first return value will get the value that corresponds to the key
// (or the map’s value type’s zero value if the key isn’t present),
// and the second will get true(or false if the key isn’t present).
func (ini *INI) GetKvmap(section string) (kvmap Kvmap, ok bool) {
	kvmap, ok = ini.sections[section]
	return kvmap, ok
}

// GetAll gets the section map and its key/value pairs.
func (ini *INI) GetAll() SectionMap {
	return ini.sections
}

// Set store the key/value pair to the default section of this INI,
// creating it if it wasn't already present.
func (ini *INI) Set(key, value string) {
	ini.SectionSet(DefaultSection, key, value)
}

// SetInt store the key/value pair to the default section of this INI,
// creating it if it wasn't already present.
func (ini *INI) SetInt(key string, value int) {
	ini.SectionSetInt(DefaultSection, key, value)
}

// SetFloat store the key/value pair to the default section of this INI,
// creating it if it wasn't already present.
func (ini *INI) SetFloat(key string, value float64) {
	ini.SectionSetFloat(DefaultSection, key, value)
}

// SetBool store the key/value pair to the default section of this INI,
// creating it if it wasn't already present.
func (ini *INI) SetBool(key string, value bool) {
	ini.SectionSetBool(DefaultSection, key, value)
}

// SectionSetInt store the section/key/value triple to this INI,
// creating it if it wasn't already present.
func (ini *INI) SectionSetInt(section, key string, value int) {
	ini.SectionSet(section, key, strconv.Itoa(value))
}

// SectionSetFloat store the section/key/value triple to this INI,
// creating it if it wasn't already present.
func (ini *INI) SectionSetFloat(section, key string, value float64) {
	ini.SectionSet(section, key, strconv.FormatFloat(value, 'f', 8, 64))
}

// SectionSetBool store the section/key/value triple to this INI,
// creating it if it wasn't already present.
func (ini *INI) SectionSetBool(section, key string, value bool) {
	var s = "false"
	if value {
		s = "true"
	}
	ini.SectionSet(section, key, s)
}

// SectionSet store the section/key/value triple to this INI,
// creating it if it wasn't already present.
func (ini *INI) SectionSet(section, key, value string) {
	kvmap, ok := ini.sections[section]
	if !ok {
		kvmap = make(Kvmap)
		ini.sections[section] = kvmap
	}
	kvmap[key] = value
}

// Delete deletes the key in given section.
func (ini *INI) Delete(section, key string) {
	kvmap, ok := ini.GetKvmap(section)
	if ok {
		delete(kvmap, key)
	}
}

// Write try to write the INI data into an output.
func (ini *INI) Write(w io.Writer) error {
	buf := bufio.NewWriter(w)

	//write the default section first
	if kv, ok := ini.GetKvmap(DefaultSection); ok {
		ini.write(kv, buf)
	}

	for section, kv := range ini.sections {
		if section == DefaultSection {
			continue
		}
		buf.WriteString("[" + section + "]" + ini.linesep)
		ini.write(kv, buf)
	}
	return buf.Flush()
}

//////////////////////////////////////////////////////////////////////////
func (ini *INI) write(kv Kvmap, buf *bufio.Writer) {
	for k, v := range kv {
		buf.WriteString(k)
		buf.WriteString(ini.kvsep)
		buf.WriteString(v)
		buf.WriteString(ini.linesep)
	}
}

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
