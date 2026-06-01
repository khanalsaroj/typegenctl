package app

import (
	"fmt"

	"github.com/khanalsaroj/typegenctl/internal/result"
	"github.com/khanalsaroj/typegenctl/util"
)

func RenderAndExit(r *result.ValidationResult, opts *Options) error {
	if opts.JSONOutput {
		return util.PrintJSON(r)
	}

	util.PrintHuman(r)

	fmt.Println("")
	return nil
}
