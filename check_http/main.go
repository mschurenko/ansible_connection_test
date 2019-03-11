package main

import (
	"../checks"
)

func main() {
	var chks []checks.Checker
	for _, i := range checks.NewHTTP().Checks {
		chks = append(chks, i)
	}

	checks.RunChecks(chks)
}
