gotomo
======
lda in golang

## TODO
1. Implement a global word count
  * Aggregate after the fact, or keep track of it as we go along?
	* JM - just implemented so that it populates the GlobalWordMap for a DocSet when GetFiles is called. 

2. Have parse function remove punctuation. 