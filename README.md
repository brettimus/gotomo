gotomo
======
lda in golang

## TODO
1. ~~Implement DiGamma Function.~~   
2. Implement Batch and Online Variational Inference
3. (lower priority) Have parse function remove punctuation. 

4. Across the program, draw distinction between 'batch' and 'update' routines.
  * For instance, an 'update' method would add documents to a DocSet and increment the appropriate counts.
	* Similarly, for the eventual inference package, we will have functionality for 'batch' and 'online' inference.
  * Presumably, the 'online' or 'update' routine will take as input a pointer to a set of update documents.
  * It will then update the model parameters without iterating through the entire corpus to date. 

5. Should DocSet.Docs be a []Document or a []*Document?
6. Global Word Count (GWC)
  * The way JM implemented it, the GWC map is populated when files are read in.
	* It is updated by merging with another docSet. 
	* There maybe is a better way to do this... 

