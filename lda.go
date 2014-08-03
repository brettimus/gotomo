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

// batch inference takes an initial DocSet and parameters, returns ptr to LdaModel
func NewLdaModel(ds DocSet, k int, alpha, beta float64) *LdaModel {
	m := len(ds.Docs)
	v := len(ds.GlobalWordMap)
	
	// initialize Lambda Randomly
	initLambda := make([][]float64, k)
	for i, _ := range initLambda {
		initLambda[i] = unifRandomSlice(v)
	}

	// initialize gamma with all ones. 
	initGamma := make([][]float64, m)
	for i, _ := range initGamma {
		initGamma[i] = allOnesSlice(k)
	}

	// initialize VarPhi with all zeroes. 
	initVarPhi := make([][][]float64, m)
	for i, _ := range initVarPhi {
		initVarPhi[i] = make([][]float64, v)
		for j, _ := range initVarPhi[i] {
			initVarPhi[i][j] = make([]float64, k)
		}
	}

	ldam :=  LdaModel{Dset: ds, Alpha: alpha, Beta: beta, K: k, Lambda: initLambda, Gamma: initGamma, VarPhi: initVarPhi}
	return &ldam
}

// To be called on initialization.
func (ldam *LdaModel) BatchInfer() {
}

// online inference is a method on an Lda Model for updating. 
func (ldam *LdaModel) OnlineInfer(ds *DocSet, kappa, tau float64, batchSize int) {}

func (ldam *LdaModel) thetaExpectation(d, k int) float64{
	return digamma(ldam.Gamma[d][k]) - digamma(sumSlice(ldam.Gamma[d]))
}

func (ldam *LdaModel) phiExpectation(k, t int) float64{
	return digamma(ldam.Lambda[k][t]) - digamma(sumSlice(ldam.Lambda[k]))
}

func (ldam *LdaModel) varPhiUpdate(d, t, k int) {}

func (ldam *LdaModel) EstParams() ([][]float64, [][]float64) {
	// returns Topic-Term Probabilities (Phi) and Doc-Topic Mixture Proportions (Theta)
	m, v := len(ldam.Dset.Docs), len(ldam.Dset.GlobalWordMap)
	k := ldam.K
	Phi, Theta := make([][]float64, k), make([][]float64, m) // kept vars here for clarity's sake, but it'd be more efficient just to return without declaring (BB)

	// Phi contains slices of length v
	for i, _ := range Phi {
		Phi[i] = make([]float64, v)
	}

	// Theta contains slices of length k
	for i, _ := range Theta {
		Theta[i] = make([]float64, k)
	}

	return Phi, Theta
}

func (ldam LdaModel) String() string {
	const str ="< LdaModel: Model with %d topics, and %d documents. >"
	return fmt.Sprintf(str, ldam.K, len(ldam.Dset.Docs))
}
