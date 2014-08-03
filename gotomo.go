package gotomo

import(
		"fmt"
)

func Test() {
	ds1 := NewDocSet("test/")
	ds1.GetFiles()
	fmt.Println(ds1)
	lda1 := NewLdaModel(*ds1, 10, 0.1, 0.1)
	fmt.Println(lda1)
	lda := *lda1
	fmt.Println(lda.Lambda)
}







