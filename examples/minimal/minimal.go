// +build slurp

//Bare minimum Slurp file.
package main //Anything, even main.

import "github.com/omeid/slurp"

func Slurp(b *slurp.Build) {
	b.Task("default", nil, func(c *slurp.C) error {
		return nil
	})
}
