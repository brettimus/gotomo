package main

import(
		"fmt"
		"strings"
		"os"
		"io"
		"io/ioutil"
		"bufio"
)

func main() {
	var ds1 = DocSet{Path:"test/", GlobalWordMap: make(map[string]int)}
	var ds2 = DocSet{Path:"test/", GlobalWordMap: make(map[string]int)}
	var wd, _ = os.Getwd()
  var stopWords = Document{File: wd + "/engStopWords.txt", WordMap: make(map[string]int)}
	stopWords.ReadFile()
	// Pass GetFiles() a Nil Document to remove no words. 
	ds1.GetFiles(NewDocument())
	// Pass a Document containing stop words to remove them from docs. 
	ds2.GetFiles(&stopWords)
	fmt.Println(ds1)
	fmt.Println(ds2)
}

type DocSet struct {
	Path           string // path to directory containing the 'Docs'. 
	Docs           []Document
	GlobalWordMap  map[string]int
}

func (ds DocSet) String() string {
	const str = "< DocSet: %d Documents in %s, Vocab has %d unique words. >"
	return fmt.Sprintf(str, len(ds.Docs), ds.Path, len(ds.GlobalWordMap))
}

// This is the 'batch' method to populate a DocSet
func (ds *DocSet) GetFiles(sw *Document) {
	dir := ds.Path
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range files {
		if !file.IsDir() {
			d := NewDocument()
			d.File = dir + file.Name()
			d.ReadFile()
			d.rmStopWords(sw)
			for word, count := range d.WordMap { // update global counts. 
				ds.GlobalWordMap[word] += count
			}
			ds.Docs = append(ds.Docs, *d)
		}
	}
}

// Update method for populating a DocSet
func (ds *DocSet) Update(other *DocSet) {
	for _, oDoc := range other.Docs {
		ds.Docs = append(ds.Docs, *oDoc)
		for word, count := range oDoc.WordMap {
			ds.GlobalWordMap[word] += count;
		}
	}
}

type Document struct {
	File      string // tag this as: "full path"
	WordCount int
	WordMap   map[string]int
}

// loop through folders and create new documents
func NewDocument() *Document {
	var d Document
	d.WordMap = make(map[string]int)
	return &d
}

func (d Document) String() string {
	const str = "<< Document: \"%s\" [%d Words, %d Unique] >>"
	return fmt.Sprintf(str,d.File, d.WordCount,len(d.WordMap))
}

// does it need an error?
// it's reading '\n' as a word
func (d *Document) ReadLine(s string) {
	for _, word := range strings.Split(s, " ") {
		d.WordMap[parse(word)]++
		d.WordCount++
	}
}

func (d *Document) ReadFile() (err error) {
	iFile := os.Stdin
	if iFile, err = os.Open(d.File); err != nil {
		return err
	}
	defer iFile.Close()
	reader := bufio.NewReader(iFile)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			return err
		}
		if err != nil {
			return err
		}
		d.ReadLine(line)
	}

}

func parse(s string) string {
	return strings.ToLower(strings.Trim(s,"\n"))
}

// Remove words that appear in a stopwords doc from d.  
func (d *Document) rmStopWords(sw *Document) {
	for key, _ := range sw.WordMap {
		delete(d.WordMap, key)
	}
}
