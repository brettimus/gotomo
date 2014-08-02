package main

import(
		"fmt"
		"strings"
		"os"
		"io"
		"io/ioutil"
		"bufio"
		"math" // for digamma
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
	ds2.Update(&ds1)
	fmt.Println(ds2)
}

type LdaModel struct {
	// model parameters
	alpha, beta float64
	K int

	// Below are variational parameters. See 'Online learning for LDA', page 3.
	varPhi [][][]float64 // named varPhi so as to distinguish from Phi, the topic-term multinom params. 
	gamma [][]float64
	lambda [][]float64

	// Need to keep track of number of docs LDA has seen so far (for online).
	numDocs int
}

// batch inference takes an initial DocSet and parameters, returns ptr to LdaModel
func batchInfer(ds *DocSet, K int, alpha, beta float64) *LdaModel {
	// dummy return value for now. 
	return new(LdaModel)
}

// online inference is a method on an Lda Model for updating. 
func (ldam *LdaModel) onlineInfer(ds *DocSet, kappa, tau float64, batchSize int) {
	// See 'Online Learning for Lda', page 4, for an explanation of kappa and tau parameters. 
	// It looks like this method can be parallelized over per-document calculations.
	// In the online algorithm, estimate of lambda is a weighted average over mini batches
	// It looks like a (boring) decision will have to be made as to wheter or not the new DocSet provided 
	// to this method should itself be considered a miniBatch, or if this method should split it further.
	// Yeah, in fact, it would be pretty damn sweet if a different procedure split it up into minibatches (also DocSets)
	// and called this function in parallel.  DO IT. 

	// *be sure to update ldam.numDocs*
}

func (ldam *LdaModel) estParams() ([][]float64, [][]float64) {
	// returns Topic-Term Probabilities and Doc-Topic Mixture Proportions
	// Will do so via the expectations listed on page 3 of 'Online Learning'

	// Dummy return value for now. 
	return (make([][]float64, 0), make([][]float64, 0))
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



// digamma.go
//
// Ported from C code by Mark Johnson. 
// Link here: http://web.science.mq.edu.au/~mjohnson/code/digamma.c
//
// checked for accuracy on wolfram alpha for digamma(.1), digamma(1.1), ..., digamma(9.1)


func digamma(x float64) float64 {
	var result, xx, xx2, xx4 float64
	for x < 7.0 {
		result -= 1/x 
		x++
	}
	x -= 1.0/2.0
	xx = 1.0/x
	xx2 = xx*xx
	xx4 = xx2*xx2
	result += math.Log(x) +(1.0/24.0)*xx2-(7.0/960.0)*xx4+(31.0/8064.0)*xx4*xx2-(127.0/30720.0)*xx4*xx4
	return result
}


