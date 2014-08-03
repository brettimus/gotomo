package gotomo

import(
		"fmt"
		"strings"
		"os"
		"io"
		"io/ioutil"
		"bufio"
)

type DocSet struct {
	Path           string // path to directory containing the 'Docs'. // needs trailing slash
	Docs           []Document
	GlobalWordMap  map[string]int
	StopWordMap    map[string]bool
	DocsCount      int    // this could be useful for allocating size of Docs + eventual channel (BB)
}

// path requires trailing slash-- parse it?
func NewDocSet(path string) *DocSet {
	sw := NewStopWordMap()
	ds := DocSet{Path: path, StopWordMap: *sw, GlobalWordMap: make(map[string]int)}
	if err := ds.setDocsCount(); err != nil {
		panic(err)
	}
	// pass a channel to GetFiles if we want to read them concurrently (BB)
	if err := ds.GetFiles(); err != nil {
		panic(err) // this is panic-worthy, right? (BB)
	}
	// TODO merge all of the documents' word maps into the Global wordmap one (BB)
	for _, doc := range ds.Docs {
		for word, count := range doc.WordMap {
			ds.GlobalWordMap[word] += count
		}
	}
	return &ds
}

// I don't like the redundancy between the next two functions.
// should use closures. halp!! (BB)
// *** SERIOUSLY USE SOME CLOUSRE *** //
func (ds *DocSet) setDocsCount() (err error) {
	var files []os.FileInfo 
	if files, err = ioutil.ReadDir(ds.Path); err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			ds.DocsCount++
		}
	}
	return err
}

// This is the 'batch' method to populate a DocSet
//
func (ds *DocSet) GetFiles() (err error) {
	files, err := ioutil.ReadDir(ds.Path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() { // is this necessary? 
			d := NewDocument(ds)
			d.File = ds.Path + file.Name()
			// TODO pass this to go routine as anon func (BB)
			d.ReadFile()

			// for word, count := range d.WordMap { // update global counts. 
			// 	ds.GlobalWordMap[word] += count
			// }
			ds.Docs = append(ds.Docs, *d)
		} // handle else of having folders of files?
	}
	return err
}

// Update method for populating a DocSet
//
// Consider making this a function that just merges two maps? (BB)
func (ds *DocSet) Update(other *DocSet) {
	for _, oDoc := range other.Docs {
		ds.Docs = append(ds.Docs, oDoc)
		for word, count := range oDoc.WordMap {
			ds.GlobalWordMap[word] += count;
		}
	}
}

type Document struct {
	File      string // tag this as: "full path"
	WordCount int
	WordMap   map[string]int
	Parent    *DocSet
}

// loop through folders and create new documents
func NewDocument(ds *DocSet) *Document {
	d := Document{Parent: ds}
	d.WordMap = make(map[string]int)
	return &d
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

// DO we need error handling here? (BB)
// it's reading '\n' as a word
func (d *Document) ReadLine(s string) {
	for _, word := range strings.Split(s, " ") {
		d.Parse(word)
	}
}

func (d *Document) Parse(s string) {
	word := strings.ToLower(strings.Trim(s,"\n"))
	if !d.Parent.StopWordMap[word] {
		d.WordMap[word]++
		d.WordCount++
	}
}



/*** Parsing Helpers ***/



// Remove words that appear in a stopwords doc from d.
//  
// *** this can be incorporated into the parse *** // -BB
func (d *Document) rmStopWords(sw *Document) {
	for key, _ := range sw.WordMap {
		delete(d.WordMap, key)
	}
}


/*** String() methods ***/

func (ds DocSet) String() string {
	const str = "< DocSet: %d Documents in %s, Vocab has %d unique words. >"
	return fmt.Sprintf(str, len(ds.Docs), ds.Path, len(ds.GlobalWordMap))
}

func (d Document) String() string {
	const str = "<< Document: \"%s\" [%d Words, %d Unique] >>"
	return fmt.Sprintf(str,d.File, d.WordCount,len(d.WordMap))
}