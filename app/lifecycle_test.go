package app

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
)

func TestLifecycle(t *testing.T) {
	app := CreateTestApp()
	testLC := NewLifecycle(
		app,
		false,
	)

	for i := 2; i < 4; i++ {
		fmt.Println("Injecting block", i)
		block := types.Block{}
		blockdoc, _ := ioutil.ReadFile(fmt.Sprintf("../test/block%d.json", i))
		utils.UnmarshalBlockResponseFromLCD(blockdoc, &block)

		runResponse := testLC.Inject(block)
		fmt.Println(runResponse)
	}
}
