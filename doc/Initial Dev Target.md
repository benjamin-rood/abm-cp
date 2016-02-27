
## Initial Development Target

    Assume Environment (E) is STATIC. 
    E = Perfect 2D world that agents 'live' in.
    No environmental factors.
    All agent interaction in 2D vector space.
    All agents perfect (constant) health until death.
    
    Visual Predators (VP) STATIC. No variation, no death, no reproduction.
    VP competitive, non-grouping.
    Model simplistic process of visual identification of prey agents, based on distance from VP, and colour distance/differentiation from VP's imprinted bias towards its target "image" of a prey agent.
    
    Colour-Polymorphic Prey (CPP) GENUINE POPULATION. Variation, reproduction, mutation, ageing
    CPP 'death' from consumption by VP and expired lifespan.
    Assume food omnipresent in E for CPP.
    Variations and Mutations reflected in individual agent Colouring.

<br>

From the beginning, given that modularity was a known requirement, the development process needed to be as organised and orderly as possible: building complete, foundational components which can reliably inter-operate and be used as cogs in the *engine* for the ABM.

One of my personal requirements was that the project be *open-source* for two reasons:

 1. When developing a program to test a research hypothesis it would seem essential that anybody should have free access interrogate the code base as a means of verifying the resultant data used in a publication.
 
 2. By using **[git](https://en.wikipedia.org/wiki/Git_%28software%29)** with the source code hosted by **[GitHub](www.github.com)** I would have free, standardised source code management. This is particularly a natural route by the fact that the command `go get` to copy and install a Go project uses **git** .

###Dealing with Vector Geometries

The first consideration was to decide how a *vector* is to be represented programmatically before we can implement the functions necessary for computing movement and relational actions by agents in our 2D vector space. In **Go** the most important data structure for convenient and efficient means of working with sequences of typed data is a *slice*, a dynamic-sized array of elements stored in contiguous memory[^GoBlog-Slices]. By taking advantage of *slices* we can write general vector operations *(cross/dot product, scalar multiplication, vector distance, rotation, etc)* which are suitable for *n-dimensional* vectors. 
By defining a customised slice `type Vector []float64` we can pass a `Vector` of any length *(dimensionality)* to a function which expects a `Vector`, and then test for its length inside the function itself if needed. *(For example, to calculate the cross product of any two vectors, we cannot perform the operation on vectors of two different lengths, nor on vectors that are not three dimensional, so the cross product function must test for `Vector` length equality.)*

I began by constructing the necessary vector operations as a distinct Go *package*[^GoBlog-packages] named **`geometry`** made up of small, reliable functions, covering the gamut of operations I would need, and which I could easily add to later.

<pre class="prettyprint"><code class="go"><font size=2>
// Vector : Any sized dimension representation of a point of vector space.
type Vector []float64

// VecScalarMultiply - performs scalar multiplication on a vector v, s.t. scalar * v = [scalar*e1, scalar*e2, scalar*e3]
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
<br></font></code></pre>
<center><font size = 2>Example function from geometry package</font>
<br>

I created unit tests for each function, providing static and random inputs to test against when making any changes to the codebase.
<br>


###Simple Behaviours

Following this, I wrote pseudocode algorithms defining the *Rule-Based Behaviour* for both agent types:




###Creating the Model Foundations

Next, I started work on my **`abm`** package to hold all of the modelling-specific constructs: 

 - **Timeline** structure to represent in-model clock as `TURNS :: PHASES :: ACTIONS`
 
 - **Context**  which stores modelling *(control)* parameters set by the user.
 
 - **Environment** structure which specifies the boundary / dimensions of the working model.
 
 - **Population** structures for both Colour-Polymorphic Prey and Visual Predators which hold:
	 - *slice* of the agent type
	 - *slice* of strings which define the characteristics of agent type

 
 - **Model** structure, with each instance acting as the working encapsulation of the modelling session, collection unique instances of all the above structures and the accompanying operations and data in memory. This allows for multiple, concurrent modelling instances on one host machine *(as would be expected from a client of a web app)*.


<br>

[^GoBlog-Slices]: <font size=1>$\ $Gerrand, Andrew (2011). *The Go Blog, Go Slices: usage and internals.* http://blog.golang.org/go-slices-usage-and-internals</font>

[^GoBlog-packages]: <font size=1>$\ $Ajmani, Sameer (2015). *The Go Blog: Package Names.* https://blog.golang.org/package-names</font>
><font size=1>Go code is organized into packages. Within a package, code can refer to any identifier (name) defined within, while clients of the package may only reference the package's exported types, functions, constants, and variables. Such references always include the package name as a prefix: foo.Bar refers to the exported name Bar in the imported package named foo. Good package names make code better. A package's name provides context for its contents, making it easier for clients to understand what the package is for and how to use it. The name also helps package maintainers determine what does and does not belong in the package as it evolves. Well-named packages make it easier to find the code you need.</font>