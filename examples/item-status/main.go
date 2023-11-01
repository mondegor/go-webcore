package main

import (
    "fmt"

    "github.com/mondegor/go-webcore/mrenum"
)

func main() {
    status := mrenum.ItemStatusEnabled

    fmt.Printf("STATUS: %s\n", status.String())

    fmt.Printf("check: %#v\n", mrenum.ItemStatusFlow.Check(mrenum.ItemStatusEnabled, mrenum.ItemStatusDisabled))
    fmt.Printf("check: %#v\n", mrenum.ItemStatusFlow.Check(mrenum.ItemStatusRemoved, mrenum.ItemStatusDisabled))

    fmt.Printf("check: %#v\n", mrenum.OnlyRemoveStatusFlow.Check(mrenum.OnlyRemoveStatusEnabled, mrenum.OnlyRemoveStatusRemoved))
    fmt.Printf("check: %#v\n", mrenum.OnlyRemoveStatusFlow.Check(mrenum.OnlyRemoveStatusRemoved, mrenum.OnlyRemoveStatusEnabled))
}
