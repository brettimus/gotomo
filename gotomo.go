package main

import(
		"fmt"
		"strings"
		"os"
		"io"
		"io/ioutil"
		"bufio"
		"math" // for digamma
	  "math/rand"
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
	lda1 := NewLdaModel(ds2, 10, 0.1, 0.1)
	fmt.Println(lda1)
	lda := *lda1
	fmt.Println(lda.Lambda)
}

type LdaModel struct {
	Dset DocSet
	Alpha, Beta float64
	K int
	VarPhi [][][]float64 // named varPhi so as to distinguish from Phi, the topic-term multinom params. 
	Gamma [][]float64
	Lambda [][]float64
}

func sumSlice(sl []float64) float64 {
	out := 0
	for _, val := range sl {
		out += val
	}
	return out
}

func unifRandomSlice(size int) []float64 {
	out := make([]float64, size)
	for i, _ := range out {
		out[i] = rand.Float64()
	}
	return out
}

func allOnesSlice(size int) []float64 {
	out := make([]float64, size)
	for i, _ := range out {
		out[i] = 1
	} 
	return out
}

// batch inference takes an initial DocSet and parameters, returns ptr to LdaModel
func NewLdaModel(ds DocSet, k int, alpha float64, beta float64) *LdaModel {
	m := len(ds.Docs)
	v := len(ds.GlobalWordMap)
	
	// initialize Lambda Randomly
	initLambda := make([][]float64, k)
	for i, _ := range initLambda {
		initLambda[i] = unifRandomSlice(v)
	}

	// initialize gamma with all ones. 
	initGamma := make([][]float64, m)
	for i, _ := range initGamma {
		initGamma[i] = allOnesSlice(k)
	}

	// initialize VarPhi with all zeroes. 
	initVarPhi := make([][][]float64, m)
	for i, _ := range initVarPhi {
		initVarPhi[i] = make([][]float64, v)
		for j, _ := range initVarPhi[i] {
			initVarPhi[i][j] = make([]float64, k)
		}
	}

	ldam :=  LdaModel{Dset: ds, Alpha: alpha, Beta: beta, K: k, Lambda: initLambda, Gamma: initGamma, VarPhi: initVarPhi}
	return &ldam
}

// To be called on initialization.
func (ldam *LdaModel) BatchInfer() {
}

// online inference is a method on an Lda Model for updating. 
func (ldam *LdaModel) OnlineInfer(ds *DocSet, kappa, tau float64, batchSize int) {}

func (ldam *LdaModel) thetaExpectation(d, k int) float64{
	return digamma(ldam.Gamma[d][k]) - digamma(sumSlice(ldam.Gamma[d]))
}

func (ldam *LdaModel) phiExpectation(k, t int) float64{
	return digamma(ldam.Lambda[k][t]) - digamma(sumSlice(ldam.Lambda[k]))
}

func (ldam *LdaModel) varPhiUpdate(d, t, k int) {}

func (ldam *LdaModel) EstParams() ([][]float64, [][]float64) {
	// returns Topic-Term Probabilities (Phi) and Doc-Topic Mixture Proportions (Theta)
	m, v := len(ldam.Dset.Docs), len(ldam.Dset.GlobalWordMap)
	k := ldam.K
	Phi, Theta := make([][]float64, k), make([][]float64, m)
}

func (ldam LdaModel) String() string {
	const str ="< LdaModel: Model with %d topics, and %d documents. >"
	return fmt.Sprintf(str, ldam.K, len(ldam.Dset.Docs))
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


