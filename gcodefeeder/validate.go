package gcodefeeder

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const maxNozzleTemp = 220

var nozzleTempRegexp = regexp.MustCompile(`^M10[49]\b`)

func ValidateGcode(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(line, ";"); idx >= 0 {
			line = line[:idx]
		}
		if !nozzleTempRegexp.MatchString(line) {
			continue
		}
		temp, ok := parseSParam(line)
		if ok && temp > maxNozzleTemp {
			return errors.New("too high temperature set, please use PLA presets only")
		}
	}
	return scanner.Err()
}

func parseSParam(line string) (int, bool) {
	for _, field := range strings.Fields(line) {
		if strings.HasPrefix(field, "S") || strings.HasPrefix(field, "s") {
			val, err := strconv.Atoi(field[1:])
			if err == nil {
				return val, true
			}
		}
	}
	return 0, false
}
