package gotomo

import(
		"fmt"
		"os"
)

func Test() {
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







