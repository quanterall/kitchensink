package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func main() {

	// The following are the strings that differ between encoder.go and decoder.go
	subs := [][]string{
		{"Encode", "Decode"},
		{"encChan", "decChan"},
		{"encRes", "decRes"},
		{"waitingEnc", "waitingDec"},
		{"//go:generate go run ./gen/.", "// generated code DO NOT EDIT"},
	}
	bytes, err := ioutil.ReadFile("encoder.go")
	if err != nil {
		log.Println(err)
		return
	}
	file := string(bytes)
	for i := range subs {
		file = strings.ReplaceAll(file, subs[i][0], subs[i][1])
	}
	err = ioutil.WriteFile("decoder.go", []byte(file), 0755)
	if err != nil {
		return
	}
}
