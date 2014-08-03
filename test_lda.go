package gotomo

import "fmt"
import "testing"

// func TestNewLdaModel() {

// }

func TestStuff(t *testing.T) {
	ds1 := NewDocSet("test/")
	fmt.Println(ds1)
	lda1 := NewLdaModel(*ds1, 10, 0.1, 0.1)
	fmt.Println(lda1)
	lda := *lda1
	fmt.Println(lda.Lambda)
}
