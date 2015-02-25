package web

import (
	"github.com/omeid/slurp"
	"github.com/omeid/slurp/tools/http"
)

// Gets  the list of urls and passes the results to output channel.
// It reports the progress to the Context using a ReadProgress proxy.
func Get(c *slurp.C, urls ...string) slurp.Pipe {

	out := make(chan slurp.File)

	go func() {
		defer close(out)

		for _, url := range urls {

			c.Infof("Downloading %s", url)

			file, err := http.Get(url)
			if err != nil {
				c.Error(err)
				continue
			}

			s, _ := file.Stat()
			file.Reader = c.ReadProgress(file.Reader, "Downloading "+file.Path, s.Size())
			out <- file
		}
	}()

	return out
}

/*
func Put(url url.URL) slurp.Stage {
	return func(files <-chan slurp.File, out chan<- slurp.File) {
		for file := range files {
			_ = file
			/*
			// */ /*
		}
	}
}

*/
