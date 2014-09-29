package main

import "fmt"

func main() {

}

type Doc []string

// fulfills iterator interface
type DocSet []Doc

type Id2Word map[int]string
type Word2Id map[string]int

type Vocab struct {
	Id2Word Id2Word
	Word2Id Word2Id
}

// map from word id to word count
type Bow map[int]int

type Corpus struct { 
	[]Bow
}