package main

import "github.com/sarojkhanal/typegenctl/cmd"

func main() {
	if cmd.Execute() != nil {
		return
	}
}
