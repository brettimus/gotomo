package gotomo

import "fmt"
// import "os"
import "testing"

func TestStuff(t *testing.T) {
	var ds1 = NewDocSet("test/")
	fmt.Println(ds1)
}