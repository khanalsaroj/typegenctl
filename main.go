package main

import "github.com/khanalsaroj/typegenctl/cmd"

func main() {
	if cmd.Execute() != nil {
		return
	}
}
