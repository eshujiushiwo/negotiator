package negotiator

import (
	"net/http"
	"strings"
)

const (
	defaultQ float64 = 1.0

	// HeaderAccept is the HTTP "Accept" Header.
	HeaderAccept = "Accept"
)

type spec struct {
	val string
	q   float64
}

// Specs it the shorthand for []Spec.
type specs []spec

// Len is used to impelement sort.Interface for Specs.
func (ss specs) Len() int {
	return len(ss)
}

// Swap is used to impelement sort.Interface for Specs.
func (ss specs) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

// Less is used to impelement sort.Interface for Specs.
func (ss specs) Less(i, j int) bool {
	if ss[i].q > ss[j].q {
		return true
	}
	if ss[i].q == ss[j].q {
		if ss[i].val == "*" || strings.HasSuffix(ss[i].val, "/*") {
			return true
		}

		if ss[j].val == "*" || strings.HasSuffix(ss[j].val, "/*") {
			return false
		}

		return i < j
	}

	return false
}

// Negotiator repensents the HTTP negotiator.
type Negotiator struct {
	req *http.Request
}

// New creates an instance of Negotiator.
func New(req *http.Request) Negotiator {
	return Negotiator{req}
}

// Accept returns the most preferred content types from the HTTP Accept header.
func (n Negotiator) Accept(offers []string) (bestOffer string, matched bool) {
	specs := parseAccept(n.req.Header)

	bestQ, bestWild := 0.0, 3

	for _, offer := range offers {
		offer = strings.ToLower(offer)

		for _, spec := range specs {
			switch {
			case spec.q < bestQ:
				continue
			case spec.val == "*/*" && (spec.q < bestQ || bestWild > 2):
				matched, bestOffer, bestQ, bestWild = true, offer, spec.q, 2
			case strings.HasSuffix(spec.val, "/*") &&
				strings.HasPrefix(offer, spec.val[:len(spec.val)-1]) &&
				(spec.q < bestQ || bestWild > 1):
				matched, bestOffer, bestQ, bestWild = true, offer, spec.q, 1
			case spec.val == offer && (spec.q < bestQ || bestWild > 0):
				matched, bestOffer, bestQ, bestWild = true, offer, spec.q, 0
			}
		}
	}

	return
}