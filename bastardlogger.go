package main

import "os"

type bastardLogger struct {
	filename string
	file     *os.File
	toStdout bool
	toFile   bool
}

func (c bastardLogger) Write(p []byte) (n int, err error) {
	if c.toFile {
		if c.file == nil {
			file, err := os.OpenFile(c.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic("Unable to open log file")
			}
			defer file.Close()
			c.file = file
		}

		c.file.Write(p)
	}
	if c.toStdout {
		os.Stdout.Write(p)
	}
	return 0, nil
}
