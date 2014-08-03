package gotomo

import(
	"fmt"
    "math/rand"
)

type LdaModel struct {
	Dset DocSet
	Alpha, Beta float64
	K int
	VarPhi [][][]float64 // named varPhi so as to distinguish from Phi, the topic-term multinom params. 
	Gamma [][]float64
	Lambda [][]float64
}

func sumSlice(sl []float64) (out float64) {
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

func NewLdaModel(ds DocSet, K int, alpha, beta float64) *LdaModel {
	m := len(ds.Docs)
	v := len(ds.GlobalWordMap)
	
	// initialize Lambda Randomly
	initLambda := make([][]float64, K)
	for i, _ := range initLambda {
		initLambda[i] = unifRandomSlice(v)
	}

	// initialize gamma with all ones. 
	initGamma := make([][]float64, m)
	for i, _ := range initGamma {
		initGamma[i] = allOnesSlice(K)
	}

	// initialize VarPhi with all zeroes. 
	initVarPhi := make([][][]float64, m)
	for i, _ := range initVarPhi {
		initVarPhi[i] = make([][]float64, v)
		for j, _ := range initVarPhi[i] {
			initVarPhi[i][j] = make([]float64, K)
		}
	}

	ldam :=  LdaModel{Dset: ds, Alpha: alpha, Beta: beta, K: K, Lambda: initLambda, Gamma: initGamma, VarPhi: initVarPhi}
	return &ldam
}

// To be called on initialization.
func (ldam *LdaModel) BatchInfer() {
	// While not converged, 
	//   E-Step then M-Step
}

// online inference is a method on an Lda Model for updating. 
func (ldam *LdaModel) OnlineInfer(ds *DocSet, kappa, tau float64, batchSize int) {}

func (ldam *LdaModel) ThetaExpectation(d, k int) float64{
	return digamma(ldam.Gamma[d][k]) - digamma(sumSlice(ldam.Gamma[d]))
}

func (ldam *LdaModel) PhiExpectation(k, t int) float64{
	return digamma(ldam.Lambda[k][t]) - digamma(sumSlice(ldam.Lambda[k]))
}

// E-Step helpers
func (ldam *LdaModel) VarPhiUpdateBatch(d, t, k int) {}
func (ldam *LdaModel) GammaUpdateBatch(d, k int) {}

// M-Step helper
func (ldam *LdaModel) LambdaUpdateBatch(k, t int) {}

func (ldam *LdaModel) EStepBatch() {}
func (ldam *LdaModel) MStepBatch() {}

func (ldam *LdaModel) EstParams() ([][]float64, [][]float64) {
	// returns Topic-Term Probabilities (Phi) and Doc-Topic Mixture Proportions (Theta)
	M, V := len(ldam.Dset.Docs), len(ldam.Dset.GlobalWordMap)
	numTopics := ldam.K
	Phi, Theta := make([][]float64, numTopics), make([][]float64, M) // kept vars here for clarity's sake, but it'd be more efficient just to return without declaring (BB)

	// Phi contains slices of length v
	for i, _ := range Phi {
		Phi[i] = make([]float64, V)
	}

	// Theta contains slices of length 'numTopics'
	for i, _ := range Theta {
		Theta[i] = make([]float64, numTopics)
	}

	// Assign each entry of Phi to its expecation.
	for k :=0; k < numTopics; k++ {
		for t := 0; t < V; t++ {
			Phi[k][t] = ldam.PhiExpectation(k, t)
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
