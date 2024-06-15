package main

import (
	"fmt"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrstatus/mrflow"
)

func main() {
	status := mrenum.ItemStatusEnabled

	fmt.Printf("STATUS: %s\n", status.String())

	flow := mrflow.ItemStatusFlow()

	fmt.Printf("check: %#v\n", flow.Check(mrenum.ItemStatusDraft, mrenum.ItemStatusEnabled))
	fmt.Printf("check: %#v\n", flow.Check(mrenum.ItemStatusEnabled, mrenum.ItemStatusDisabled))
	fmt.Printf("check: %#v\n", flow.Check(mrenum.ItemStatusDisabled, mrenum.ItemStatusDraft))
}
