package main

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"
)

func main() {
	status := itemstatus.Enabled

	fmt.Printf("STATUS: %s\n", status.String())

	flow := itemstatus.NewFlowMap()

	fmt.Printf("is possible: %#v\n", flow.IsPossible(itemstatus.Draft, itemstatus.Enabled))
	fmt.Printf("is possible: %#v\n", flow.IsPossible(itemstatus.Enabled, itemstatus.Disabled))
	fmt.Printf("is possible: %#v\n", flow.IsPossible(itemstatus.Disabled, itemstatus.Draft))
}
