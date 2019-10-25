package main

import (
	"github.com/klovack/bambi/pkg/cli/bambi"
	"github.com/klovack/bambi/pkg/util"
)

func main() {
	err := bambi.NewCommand().Execute()

	util.CheckErrorP(err)
}
