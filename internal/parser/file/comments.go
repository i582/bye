package file

import (
	"bytes"
	"strings"
)

func isEmptyLine(line []byte) bool {
	line = bytes.TrimLeft(line, " \t\r\n")
	line = bytes.TrimRight(line, " \t\r\n")
	return len(line) == 0
}

func isAllowedComment(line []byte) bool {
	line = bytes.TrimLeft(line, " \t")

	switch {
	case bytes.HasPrefix(line, []byte("//")):
		return bytes.HasPrefix(line, []byte("//!"))
	case bytes.HasPrefix(line, []byte("<!--")):
		return bytes.HasPrefix(line, []byte("<!--!"))
	case bytes.HasPrefix(line, []byte("#")) && !bytes.HasPrefix(line, []byte("#[")):
		return bytes.HasPrefix(line, []byte("#!"))
	}

	return false
}

func trimComment(line []byte) []byte {
	line = bytes.TrimRight(line, " \t\n\r")
	line = bytes.TrimLeft(line, " #\t\n\n/")
	line = bytes.TrimPrefix(line, []byte("<!--"))
	line = bytes.TrimSuffix(line, []byte("-->"))

	return line
}

func joinLines(lines [][]byte) string {
	var res string

	for _, line := range lines {
		if isEmptyLine(line) {
			res += "<br>"
			continue
		}

		res += string(line) + " "
	}

	res = strings.TrimSuffix(res, "<br>")

	return res
}
