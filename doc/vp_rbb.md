
### Rule-Based Behaviour: Visual Predator ###

#### 1. Begin ####
	Increment vp.hunger
	
	IF PREDATOR AGEING is true:
		Decrement vp.lifetime
		IF lifetime ≤ 0: Jump to DEATH
	
	IF hunger > HUNGER LIMIT and STARVATION is true:
		Jump to DEATH
	
	IF hunger < PREDATOR SEXUAL REQUIREMENTS:
		Jump to MATE SEARCH
	
	ELSE: Jump to PREY SEARCH
	
#### 2. (Visual) Prey Search ####
	SEARCH PREY population of SECTORS within vsr:
		where FIND PREY ⟮ ∂ * γ * colour∂ ⟯ > vsrSearchChance
	
	IF FIND PREY: Jump to ATTACK
	
	ELSE: Jump to PATROL
	
####3. Attack ####
	
	Let α be a random real value in [0, 1)
	
	IF α > vpAttackChance: 
		MOVE to PREY position
		Jump to EAT
	
	ELSE: 
		MOVE to fuzzy offset of PREY position (missed target)
		Jump to END

####4.  Eat####

	Remove PREY from model population and sector location.
	Decrement hunger by ℮
	Jump to END

####5. Patrol ####
	Let ϴ be a random angle (radians) in [-π/2, π/2] 
	IF ϴ > zero: vp turns left by ϴ radians.

	MOVE forward according to new heading, at distance 
	(movS * movA) from current position.
	Jump to END

	Further Details:
	ϴ to be added to the existing heading of vp, s.t. heading = U 
	Where (heading + ϴ) ≣ U (mod 2π)
	heading ⟵ U
	direction ⟵ Vec2f{cosU, sineU}

####6. Death####
	Remove vp from model population and sector location.
	Jump to END

####7. End####
	Increment action counter for model time.
	?
	EXIT
	
####List of abbreviations,  variables:####

	E = model environment (subset of Rⁿ space)
		In a 2D environment, E is a field of (-d,d) × (-d,d).
		
	d = the unit measurement of E
		e.g. 2*d = length of E, if E is a 2D space.
		
	vp = Visual Predator
	
	∂ = distance between vp and target
	
	γ = visual acuity of vp (operating in E)
	
	colour∂ = colour distance (difference)
	
	℮ = energy gained by vp consuming prey
	
	movS, movA = speed, acceleration of vp
	
	vsrSearchChance, vpAttackChance = 
		constant values to bias likelihood of
		successfully finding, eating PREY