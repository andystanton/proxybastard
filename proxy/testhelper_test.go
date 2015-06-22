package proxy

import "strings"

type textFile struct {
	lines []string
}

type shellFile struct {
	shellLines []ShellLine
}

func createShellFile() shellFile {
	return shellFile{[]ShellLine{}}
}

func (sf shellFile) addShellLine(line ShellLine) shellFile {
	sf.shellLines = append(sf.shellLines, line)
	return sf
}

func (sf shellFile) toSlice() []ShellLine {
	return sf.shellLines
}

func (sf shellFile) toString() string {
	output := []string{}
	for _, shellLine := range sf.shellLines {
		output = append(output, shellLine.lines...)
	}
	return strings.Join(output, "\n")
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
