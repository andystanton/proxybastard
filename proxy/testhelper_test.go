package proxy

import "strings"

type textFile struct {
	lines []string
}

func createTextFile() textFile {
	return textFile{[]string{}}
}

func (sf textFile) addLine(line string) textFile {
	sf.lines = append(sf.lines, line)
	return sf
}

func (sf textFile) toSlice() []string {
	return sf.lines
}

func (sf textFile) toString() string {
	return strings.Join(sf.lines, "\n")
}
