package util

import "os"

// TeeLogger acts like the GNU tee command - logs to stdout & file.
type TeeLogger struct {
	Filename string
	file     *os.File
	ToStdout bool
	ToFile   bool
}

func (c TeeLogger) Write(p []byte) (n int, err error) {
	if c.ToFile {
		if c.file == nil {
			file, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic("Unable to open log file")
			}
			defer file.Close()
			c.file = file
		}

		c.file.Write(p)
	}
	if c.ToStdout {
		os.Stdout.Write(p)
	}
	return 0, nil
}
