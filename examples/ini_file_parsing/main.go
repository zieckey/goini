// Copyright 2014 zieckey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"path/filepath"

	"github.com/zieckey/goini"
)

func main() {
	filename := filepath.Join("ini_parser_testfile.ini")
	ini := goini.New()
	err := ini.ParseFile(filename)
	if err != nil {
		fmt.Printf("parse INI file %v failed : %v\n", filename, err.Error())
		return
	}

	v, ok := ini.Get("mid")
	if !ok || v != "ac9219aa5232c4e519ae5fcb4d77ae5b" {
		fmt.Printf("mid value is invalid\n")
		return
	}
	fmt.Printf("the value of [mid] in INI is [%v] ok=[%v]\n", v, ok)

	sss, ok := ini.GetKvmap("sss")
	size := len(sss)
	if size != 2 {
		fmt.Printf("sss size is invalid\n")
		return
	}
	
	v, ok = ini.SectionGet("sss", "aa")
	if !ok || v != "bb" {
		fmt.Printf("sss/aa value is invalid\n")
		return
	}
	
	v, ok = ini.SectionGet("sss", "appext")
	if !ok || v != "ab=cd" {
		fmt.Printf("sss/appext value is invalid\n")
		return
	}
}

