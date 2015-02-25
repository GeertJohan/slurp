package util

import (
	"bytes"

	"github.com/omeid/slurp"
)

// Concatenates all the files from the input channel
// and passes them to output channel with the given name.
func Concat(c *slurp.C, name string) slurp.Stage {
	return func(files <-chan slurp.File, out chan<- slurp.File) {

		var (
			size    int64
			bigfile = new(bytes.Buffer)
		)

		for f := range files {
			c.Infof("Adding %s to %s", f.Path, name)
			n, err := bigfile.ReadFrom(f)
			if err != nil {
				c.Error(err)
				return
			}
			bigfile.WriteRune('\n')
			size += n + 1

			f.Close()
		}

		sf := slurp.File{
			Reader: bigfile,
			Dir:    "",
			Path:   name,
		}
		stat := &slurp.FileInfo{}
		stat.SetSize(size)
		stat.SetName(name)

		sf.SetStat(stat)

		out <- sf
	}
}

//For The Glory of Debugging.
func List(c *slurp.C) slurp.Stage {
	return func(files <-chan slurp.File, out chan<- slurp.File) {
		for f := range files {
			s, err := f.Stat()
			if err != nil {
				c.Error("Can't get File Stat name.")
			} else {
				c.Infof("slurp.File: %+v Name: %s", f, s.Name())
			}
			out <- f
		}
	}
}
