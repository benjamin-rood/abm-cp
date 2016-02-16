package abm

import (
	"bytes"
	"fmt"
)

// String returns a clear textual presentation the internal values of the CPP agent
func (c *ColourPolymorphicPrey) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("pos=(%v,%v)\n", c.pos[x], c.pos[y]))
	buffer.WriteString(fmt.Sprintf("movS=%v\n", c.movS))
	buffer.WriteString(fmt.Sprintf("movA=%v\n", c.movA))
	buffer.WriteString(fmt.Sprintf("𝚯=%v\n", c.𝚯))
	buffer.WriteString(fmt.Sprintf("dir=(%v,%v)\n", c.dir[x], c.dir[y]))
	buffer.WriteString(fmt.Sprintf("tr=%v\n", c.tr))
	buffer.WriteString(fmt.Sprintf("sr=%v\n", c.sr))
	buffer.WriteString(fmt.Sprintf("lifespan=%v\n", c.lifespan))
	buffer.WriteString(fmt.Sprintf("hunger=%v\n", c.hunger))
	buffer.WriteString(fmt.Sprintf("fertility=%v\n", c.fertility))
	buffer.WriteString(fmt.Sprintf("gravid=%v\n", c.gravid))
	buffer.WriteString(fmt.Sprintf("colouration=%v\n", c.colouration))
	return buffer.String()
}

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
