gotomo
======
lda in golang

## TODO
1. Implement DiGamma Function.
2. Implement Batch and Online Variational Inference
3. (lower priority) Have parse function remove punctuation. 

3. Across the program, draw distinction between 'batch' and 'update' routines.
  * For instance, an 'update' method would add documents to a DocSet and increment the appropriate counts.
	* Similarly, for the eventual inference package, we will have functionality for 'batch' and 'online' inference.
  * Presumably, the 'online' or 'update' routine will take as input a pointer to a set of update documents.
  * It will then update the model parameters without iterating through the entire corpus to date. 

4. Should DocSet.Docs be a []Document or a []*Document?