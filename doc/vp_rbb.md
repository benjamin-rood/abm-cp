
### Rule-Based Behaviour: Visual Predator ###

#### 1. Begin ####
	Increment agent hunger, fertility
	
	IF PREDATOR STARVATION is true:
		IF hunger causes PANIC:
			Increase range of γ (visual acuity)
			Decrease ετ (colour specialisation)
	
	IF PREDATOR AGEING is true:
		Decrement agent lifespan
		IF lifetime ≤ zero: Jump to DEATH
	
	IF PREDATOR STARVATION ⋀ hunger > limit: Jump to DEATH

	IF fertility is zero: Jump to SPAWN

	IF fertility ⋀ hunger > PREDATOR ENERGY REQUIREMENT: Jump to MATE SEARCH
	
	ELSE: Jump to PREY SEARCH

#### 5.	 Mate Search
	
	Search neighbouring PREDATOR agents
	IF neighbour fertility >  PREDATOR ENERGY REQUIREMENT:
		INTERCEPT target
		Attempt to COPULATE, decrement neighbour fertility by SEXUAL COST
		IF copulation successful:
			agent fertility ⟵ -
			
			

#### 2. Prey (Visual) Search ####
    SEARCH PREY population of SECTORS within vsr 
    A given PREY agent is seen based on the formula:

$$f (\chi) = ce^{c\chi} + \ell$$
	
*Where $\chi$ is the colour difference between PREY agent colouration and PREDATOR agent colour target $\, \tau \,$ , $\,\  c\, $ is the colour specialisation strength $\epsilon\tau \,$ ,  $\,$ and $\,\ \ell = \, \frac{1}{c}\, \,  ∀\, \ c\, \geq\,  1 $*
	
	IF ƒ(χ) > γ for some agent in PREY population:
		Jump to INTERCEPT
		
	ELSE: Jump to PATROL
	
	
#### 3. Intercept
	
	Let Ψ be the difference in angle between the current
	the current relative position of the target.
	
	Calculate the distance from the current agent position
	and the target.
	
	IF dist < movS: 
		move agent to target position
		Jump to ATTACK
	
	ELSE: 
		TURN agent heading by Ψ radians
		MOVE forward according to new heading, at distance
		(movS·movA) from current position.
		Jump to END
	

####4. Attack ####
	Let α be a random real value in [0, 1]
	
	IF α > (1 - vpAttackChance): 	
		Attack Success ⟵ true
		τ ⟵ (prey colouration - τ)·CAF
		Decrement hunger by ƒ(χ)·bg
		Increment ετ
		Remove prey agent from population.
		
	Jump to END

####6. Patrol ####
	Let ϴ be a random angle in (-turnRate, +turnRate)
	 
	IF ϴ > zero: agent turns left by ϴ radians
	
	If ϴ < zero: agent turns right by ϴ radians
	
	TURN agent heading by ϴ radians

	MOVE forward according to new heading, at distance (movS·movA) from current position.
	
	Jump to END
	

####7. Death####
	Remove agent from PREDATOR population and current sector record.
	
	Jump to END

####8. End####
	Increment action counter for model time.
	
	EXIT
	
####List of abbreviations,  variables:####

	E = model environment (subset of Rⁿ space)
		In a 2D environment, E is a field of (-d,d) × (-d,d).
		
	d = the unit measurement of E
		e.g. 2*d = length of E, if E is a 2D space.
		
	vp = Visual Predator
	
	∂ = distance between vp and target
	
	γ = visual acuity of vp (operating in E)
	
	χ = colour distance (difference)
	

	bg = minimum energy gained by vp consuming prey

	ετ = Colour specialisation strength
	
	movS, movA = speed, acceleration of vp
	
	vpAttackChance = 
		constant values to bias likelihood of
		successfully finding, eating PREY