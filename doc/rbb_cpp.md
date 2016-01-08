### Rule-Based Behaviour: CP Prey ###
####1. Begin####
	IF PREY AGEING is true:
		decrement lifetime
		IF lifetime ≤ 0: Jump to DEATH
	Increment fertility
	IF fertility = 0:
		gravid ⟵ false
		Jump to SPAWN
	ELSE IF fertility ≥ period: Jump to MATE SEARCH
	ELSE: Jump to EXPLORE

####2.  Mate Search ⧸⧸(local) ####⧸⧸####⧸⧸
	SEARCH PREY population of MICRO SECTORS within range
	IF FOUND MATE: 
		position ⟵ MATE(position)
		Jump to ATTEMPT REPRODUCE
	ELSE: Jump to EXPLORE

####3. Attempt Reproduce ####
	Let ω be a random real number in [0, 1)
	IF ω ≤ sReproduceChance:
		fertility ⟵ -gestation
		gravid ⟵ true
	Jump to END

####4. Explore####
	MOVE to random position within radius (movS*movA)
	Jump to END

####5. Spawn####
	Let n be a random integer in [1, b]
	Add n children inheriting colour from  to PREY population to a random proximate position.
	Jump to END

####6. Death#####
	Remove cpp object from PREY population and MICRO SECTOR location
	Jump to END

####7. End####
	Increment action counter for model time
	?
	EXIT 
