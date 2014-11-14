package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	//"testing"

	"github.com/zieckey/wini"
	//"github.com/bmizerany/assert"
)

func main() {
	test1()
	test2()
}

func test1() {

}

func test2() {
	/*
	mid=ac9219aa5232c4e519ae5fcb4d77ae5b
	product=ppp
	combo=ccc
	version=4.4
	#appext=abcd
	appext= abcd
	;a=b
	aa=bb
	*/
	filename := filepath.Join(testDataDir(), "ini_parser_testfile.ini")
	ini := wini.New()
	err := ini.ParseFile(filename)
	if err != nil {
		fmt.Printf("parse INI file %v failed : %v\n", filename, err.Error())
		return
	}

	v, ok := ini.Get("", "mid")
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
	v, ok = ini.Get("sss", "aa")
	if !ok || v != "bb" {
		fmt.Printf("sss/aa value is invalid\n")
		return
	}
	v, ok = ini.Get("sss", "appext")
	if !ok || v != "ab=cd" {
		fmt.Printf("sss/appext value is invalid\n")
		return
	}

}

func testDataDir() string {
	var file string
	var ok bool
	if _, file, _, ok = runtime.Caller(0); ok {
		fmt.Printf("file=%v\n", file)
	}

	curdir := filepath.Dir(file)
	return filepath.Join(curdir, "data")
}
