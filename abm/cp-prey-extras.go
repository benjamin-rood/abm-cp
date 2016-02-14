package abm

// extra CPP functions for testing/benchmarking:

func cppTesterAgent(xPos float64, yPos float64) (tester ColourPolymorphicPrey) {
	tester = cppTestPop(1)[0]
	tester.pos[x] = xPos
	tester.pos[y] = yPos
	return
}

func newCppTesterAgent(xPos float64, yPos float64) *ColourPolymorphicPrey {
	tester := cppTestPop(1)[0]
	tester.pos[x] = xPos
	tester.pos[y] = yPos
	return &tester
}

func cppTestPop(size int) []ColourPolymorphicPrey {
	return GeneratePopulationCPP(size, 0, 0, TestContext, testStamp)
}
