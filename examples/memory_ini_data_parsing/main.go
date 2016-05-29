// Copyright 2014 zieckey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/zieckey/goini"
)

func f1() {
	raw := []byte("a:av|b:bv||c:cv|||d:dv||||||")
	ini := goini.New()
	err := ini.Parse(raw, "|", ":")
	if err != nil {
		fmt.Printf("parse INI memory data failed : %v\n", err.Error())
		return
	}

	key := "a"
	v, ok := ini.Get(key)
	if ok {
		fmt.Printf("The value of %v is [%v]\n", key, v) // Output : The value of a is [av]
	}

	key = "c"
	v, ok = ini.Get(key)
	if ok {
		fmt.Printf("The value of %v is [%v]\n", key, v) // Output : The value of c is [cv]
	}
}

func f2() {
	raw := []byte("a:av||b:bv||c:cv||||d:dv||||||")
	ini := goini.New()
	err := ini.Parse(raw, "||", ":")
	if err != nil {
		fmt.Printf("parse INI memory data failed : %v\n", err.Error())
		return
	}

	key := "a"
	v, ok := ini.Get(key)
	if ok {
		fmt.Printf("The value of %v is [%v]\n", key, v) // Output : The value of a is [av]
	}

	key = "c"
	v, ok = ini.Get(key)
	if ok {
		fmt.Printf("The value of %v is [%v]\n", key, v) // Output : The value of c is [cv]
	}
}

func main() {
	f1()
	fmt.Print("\n\n\n")
	f2()
}
