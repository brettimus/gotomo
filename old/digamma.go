package gotomo

import "math"

// Ported from C, thanks be to Mark Johnson. 
// Link here: http://web.science.mq.edu.au/~mjohnson/code/digamma.c
//
// Accuracy up to 7 sigfigs checked against Wolfram for digamma(.1), digamma(1.1), ..., digamma(9.1)

func digamma(x float64) float64 {
	var result, xx, xx2, xx4 float64
	for x < 7.0 {
		result -= 1/x 
		x++
	}
	x -= 1/2.0
	xx = 1/x
	xx2 = xx * xx
	xx4 = xx2 * xx2
	result += math.Log(x) + (1.0/24.0)*xx2 - (7.0/960.0)*xx4 + (31.0/8064.0)*xx4*xx2 - (127.0/30720.0)*xx4*xx4
	return result
}
