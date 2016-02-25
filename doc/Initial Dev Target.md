
### Initial Development Target

    Assume Environment (E) is STATIC. 
    E = Perfect 2D world that agents 'live' in.
    No environmental factors.
    All agent interaction in 2D vector space.
    All agents perfect (constant) health until death.
    
    Visual Predators (VP) STATIC. No variation, no death, no reproduction.
    VP competitive, non-grouping.
    Model simplistic process of visual identification of prey agents, based on distance from VP, colour difference from background environment, and colour distance from VP's imprinted bias towards its "image" of a prey agent.
    
    Colour-Polymorphic Prey (CPP) GENUINE POPULATION. Variation, reproduction, mutation, ageing
    CPP 'death' from consumption by VP and expired lifespan.
    Assume food omnipresent in E for CPP.
    Variations and Mutations reflected in individual agent Colouring.

<br>

From the beginning, as there needed to be some actual end-use, one cannot just hack together some code which only seems to work, unless you want to have something that's going to break all the time — the development process needs to be as organised and orderly as possible: building complete, foundational components which can reliably inter-operate and be used as cogs in the *engine* for the ABM.

####Dealing with Vector Geometries

The first consideration is to decide how a *vector* is to be represented programatically before we can implement the functions necessary for computing movement and relational actions by agents in our 2D vector space. In **Go** the most important data structure for convenient and efficient means of working with sequences of typed data is a *slice*, a dynamic-sized array of elements stored in contiguous memory[^GoBlog-Slices]. By taking advantage of *slices* we can write general vector operations *(cross/dot product, scalar multiplication, vector distance, rotation, etc)* which are suitable for *n-dimensional* vectors. 
By defining a customised slice `type Vector []float64` we can pass a `Vector` of any length *(dimensionality)* to a function which expects a `Vector`, and then test for its length inside the function itself if needed. *(For example, to calculate the cross product of any two vectors, we cannot perform the operation on vectors of two different lengths, nor on vectors that are not three dimensional, so the cross product function must test for `Vector` length equality.)*

I began by constructing the necessary vector operations as a distinct Go *package* named **geometry** made up of small, reliable functions.

	// Vector : Any sized dimension representation of a point of vector space.
	type Vector []float64

	type VectorEquality interface {
		Equal(VectorEquality) bool
	}

	// Equal method implements an Equality comparison between vectors.
	func (v Vector) Equal(u VectorEquality) bool {
		if len(v) != len(u.(Vector)) {
			return false
		}
		for i := 0; i < len(v); i++ {
			if v[i] != u.(Vector)[i] {
				return false
			}
		}
		return true
	}
	
	// VecScalarMultiply - performs scalar multiplication on a vector v,
	// s.t. scalar * v = [scalar*e1, scalar*e2, scalar*e3]
	
	func VecScalarMultiply(v Vector, scalar float64) (Vector, error) {
		if len(v) == 0 {
			return nil, errors.New("v is an empty vector")
		}
		var sm Vector
		for i := range v {
			sm = append(sm, (v[i] * scalar))
		}
		return sm, nil
	}




[^GoBlog-Slices]:  Gerrand, Andrew (2011). *The Go Blog, Go Slices: usage and internals.* http://blog.golang.org/go-slices-usage-and-internals
