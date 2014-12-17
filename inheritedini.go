package goini

import (
	"path/filepath"
	"errors"
	"log"
)

// Suppress error if they are not otherwise used.
var _ = log.Printf

const (
	InheritedFrom = "inherited_from" // The key of the INI path which will be inherited from
)


func LoadInheritedINI(filename string) (*INI, error) {
	ini := New()
	err := ini.ParseFile(filename)
	if err != nil {
		return nil, err
	}
	
	inherited, ok := ini.Get(InheritedFrom)
	if !ok {
		return ini, nil
	}
	
	inherited = realpath(filename, inherited)
	inheritedINI, err := LoadInheritedINI(inherited)
	if err != nil {
		return nil, errors.New(err.Error() + " " + inherited)
	}
	
	ini.Merge(inheritedINI, false)
	return ini, nil
}

// Merge merges the data in another INI (from) to this INI (ini), and
// from INI will not be changed
func (ini *INI) Merge(from *INI, override bool) {
	for section, kv := range from.sections {
		for key, value := range kv {
			_, found := ini.SectionGet(section, key)
			if override || !found {
				ini.SectionSet(section, key, value)
			}
		}
	}
}

func realpath(currentPath, inheritedPath string) string {
	if filepath.IsAbs(inheritedPath) {
		return inheritedPath
	}
	
	dir, _ := filepath.Split(currentPath)
	return filepath.Join(dir, inheritedPath)
}