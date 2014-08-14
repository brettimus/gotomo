package gotomo

import(
	"fmt"
	"math/rand"
	"math"
)

type LdaModel struct {
	Dset DocSet
	Alpha, Beta float64
	K int
	VarPhi []map[string][]float64
	Gamma [][]float64
	Lambda []map[string]float64
}

// will this be called often? 
// i.e., would it be more efficient to store a running total in memory?
// - BB
// Methinks you're right. - JM
func sumSlice(s []float64) (out float64) {
	for _, val := range s {
		out += val
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

func NewLdaModel(ds DocSet, K int, alpha, beta float64) *LdaModel {
	m := len(ds.Docs)
	
	// initialize Lambda Randomly
	initLambda := make([]map[string]float64, K)
	for i, _ := range initLambda {
		initLambda[i] = make(map[string]float64)
		for term, _ := range ds.GlobalWordMap {
			initLambda[i][term] = rand.Float64()
		}
	}

	// initialize gamma with all ones. 
	initGamma := make([][]float64, m)
	for i, _ := range initGamma {
		initGamma[i] = allOnesSlice(K)
	}

	// initialize VarPhi with all zeroes. 
	initVarPhi := make([]map[string][]float64, m)
	for i, _ := range initVarPhi {
		initVarPhi[i] = make(map[string][]float64)
		for term, _ := range ds.GlobalWordMap {
			initVarPhi[i][term] = make([]float64, K)
		}
	}

	ldam :=  LdaModel{Dset: ds, Alpha: alpha, Beta: beta, K: K, Lambda: initLambda, Gamma: initGamma, VarPhi: initVarPhi}
	return &ldam
}

// To be called on initialization.
func (ldam *LdaModel) BatchInfer() {
	// This is the last thing that needs to be written for batch. 
}

// online inference is a method on an Lda Model for updating. 
// func (ldam *LdaModel) OnlineInfer(ds *DocSet, kappa, tau float64, batchSize int) {}

func (ldam *LdaModel) ThetaExpectation(d, k int) float64{
	return digamma(ldam.Gamma[d][k]) - digamma(sumSlice(ldam.Gamma[d]))
}

func (ldam *LdaModel) PhiExpectation(k int, t string) (out float64) {
	for term, _ := range ldam.Dset.GlobalWordMap {
		out += ldam.Lambda[k][term]
	}
	return digamma(ldam.Lambda[k][t]) - out
}

// E-Step helpers
func (ldam *LdaModel) VarPhiUpdateBatch(d, k int, t string) (new float64) {
	new = math.Exp(ldam.ThetaExpectation(d, k) + ldam.PhiExpectation(k, t))
	return new
}

func (ldam *LdaModel) GammaUpdateBatch(d, k int) (new float64) {
	var sum float64
	for term, _ := range ldam.Dset.GlobalWordMap {
		sum += float64(ldam.Dset.Docs[d].WordMap[term])*ldam.VarPhi[d][term][k]
	}
	new = ldam.Alpha + sum
	return new
}

// M-Step helper
func (ldam *LdaModel) LambdaUpdateBatch(k int, t string) (out float64) {
	for index, doc := range ldam.Dset.Docs {
		out += float64(doc.WordMap[t])*ldam.VarPhi[index][t][k]
	}
	return ldam.Beta + out
}

func (ldam *LdaModel) EStepBatch() (diff float64) {
// Update VarPhi
	for d,_ := range ldam.Dset.Docs {
		for term, _ := range ldam.Dset.GlobalWordMap {
			for k := 0; k < ldam.K; k++ {
				ldam.VarPhi[d][term][k] = ldam.VarPhiUpdateBatch(d, k, term)				
			}
		}
	}

	// Update Gamma
	for d,_ := range ldam.Dset.Docs {
		for k:=0; k < ldam.K; k++ {
			old, new := ldam.Gamma[d][k], ldam.GammaUpdateBatch(d, k)
			diff += math.Abs(new - old)
		}
	}

	return diff/float64(ldam.K)
}

func (ldam *LdaModel) MStepBatch() {
	for k := 0; k < ldam.K; k++ {
		for term, _ := range ldam.Dset.GlobalWordMap {
			ldam.Lambda[k][term] = ldam.LambdaUpdateBatch(k, term)
		}
	}
}

func (ldam *LdaModel) EstParams() ([]map[string]float64, [][]float64) {
	// returns Topic-Term Probabilities (Phi) and Doc-Topic Mixture Proportions (Theta)
	M := len(ldam.Dset.Docs)
	numTopics := ldam.K
	Phi, Theta := make([]map[string]float64, numTopics), make([][]float64, M) // kept vars here for clarity's sake, but it'd be more efficient just to return without declaring (BB)

	// Phi contains slices of length v
	for i, _ := range Phi {
		Phi[i] = make(map[string]float64)
	}

	// Theta contains slices of length 'numTopics'
	for i, _ := range Theta {
		Theta[i] = make([]float64, numTopics)
	}

	// Assign each entry of Phi to its expecation.
	for k :=0; k < numTopics; k++ {
		for term, _ := range ldam.Dset.GlobalWordMap {
			Phi[k][term] = ldam.PhiExpectation(k, term)
		}
	}

	// Assign each entry of Theta to its expectation.
	for d:=0; d < M; d++ {
		for k:=0; k < numTopics; k++ {
			Theta[d][k] = ldam.ThetaExpectation(d, k)
		}
	}
	
	return Phi, Theta
}

func (ldam LdaModel) String() string {
	const str ="< LdaModel: Model with %d topics, and %d documents. >"
	return fmt.Sprintf(str, ldam.K, len(ldam.Dset.Docs))
}
