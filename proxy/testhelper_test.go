package proxy

import "strings"

type textFile struct {
	lines []string
}

type shellFile struct {
	shellLines []ShellStatement
}

func createShellFile() shellFile {
	return shellFile{[]ShellStatement{}}
}

func (sf shellFile) addShellStatement(line ShellStatement) shellFile {
	sf.shellLines = append(sf.shellLines, line)
	return sf
}

func (sf shellFile) toSlice() []ShellStatement {
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
