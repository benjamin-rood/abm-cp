
##Definition of the ABM Environment

Consider the set of all real numbers within $[\,-1,  \,+1\,]$.$\ $ Let $\,E\ $ be the cross product of such a set with itself,  *s.t.* individual elements in $\,E\,$ are two-dimensional points.
$$\therefore \ E\ = \  \{\ (x,y) \in\, \mathbb{R}^{2} \, \lvert\  (\,-1 \leq x \leq 1 \,)\ \wedge\,(\,-1 \leq y \leq 1 \,)\  \}$$

The range of $E$ will be sufficient to act as a "petrie dish" to observe the behaviour of the two populations of agents, with the point $\,(0,0)\,$ as the absolute center of our surface.

Rather than agents 'bouncing' or 'avoiding' the boundaries of $E$, values $(x,y) \in  E\,$ **wrap** when exceeding the limits of the surface area, *s.t.* the surface is in effect quasi-spherical. 

$\ \forall \ a \in [-1, +1]\,$ we have three cases:

Case $\ 1 < a \ :$
: $ a \longmapsto (a-2)$

<br>
Case $\ -1 < a < 1 \ :$ 
: $a = a$

<br>
Case $\ a < 1 :$
 : $a \longmapsto (a+2)$

 
*For example, if an agent moves from position $\,(-0.9, 0.2)\,$ down by 0.3 it would map from  $\,(-1.2, 0.2)\,$ to $\,(0.8, 0.2)$; similarly, if an agent moves from position $\,(1.0, 0.95)$ to $(1.13, 1.025)\,$ the resulting position would be $\,(-0.87,-0.975)$*


<br>

![enter image description here](http://i.imgur.com/DGgASF3.jpg)

<CENTER><font size=2>*Planar projection sketch of $E$ with various Colour Polymorphic Prey agents *</font></CENTER>