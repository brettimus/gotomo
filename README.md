gotomo
======
lda in golang

## TODO
1. Implement a global word count
  * Aggregate after the fact, or keep track of it as we go along?
	* JM - just implemented so that it populates the GlobalWordMap for a DocSet when GetFiles is called. 

2. Have parse function remove punctuation. 

3. Across the program, draw distinction between 'batch' and 'update' routines.
  * For instance, an 'update' method would add documents to a DocSet and increment the appropriate counts.
	* Similarly, for the eventual inference package, we will have functionality for 'batch' and 'online' inference.
  * Presumably, the 'online' or 'update' routine will take as input a pointer to a set of update documents.
  * It will then update the model parameters without iterating through the entire corpus to date. 