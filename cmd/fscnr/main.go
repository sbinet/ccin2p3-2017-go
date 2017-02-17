package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(2)
	}

	fname, beg, end, err := parseCmd(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("fname: %v\n", fname)

	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	if end == -1 {
		end = len(buf)
	}

	if end < beg {
		log.Fatalf("invalid file range: start=%d end=%d\n", beg, end)
	}

	log.Printf("start: %d\n", beg)
	log.Printf("end:   %d\n", end)

	log.Printf("=== %v ===\n%s\n", fname, hex.Dump(buf))
	log.Printf("buf[%d:%d]: %s\n", beg, end, buf[beg:end])
}

func parseCmd(cmd string) (fname string, beg int, end int, err error) {
	const errfmt = "invalid arg: %q. (want <path-to-file>:#<beg>[,#<end>]\n"
	toks := strings.Split(cmd, ":")
	if len(toks) != 2 {
		err = fmt.Errorf(errfmt, cmd)
		return
	}

	fname = toks[0]
	toks = strings.Split(toks[1], ",")
	switch len(toks) {
	case 2:
		var v int64
		v, err = parsePos(toks[0])
		if err != nil {
			return
		}
		beg = int(v)

		v, err = parsePos(toks[1])
		if err != nil {
			return
		}
		end = int(v)
	case 1:
		var v int64
		v, err = parsePos(toks[0])
		if err != nil {
			return
		}
		beg = int(v)
		end = -1
	default:
		err = fmt.Errorf(errfmt, cmd)
		return
	}

	return
}

func parsePos(str string) (int64, error) {
	const errfmt = "invalid arg: %q. (want \"#<position>\""
	if !strings.HasPrefix(str, "#") || len(str) <= 1 {
		return 0, fmt.Errorf(errfmt, str)
	}
	return strconv.ParseInt(str[1:], 10, 64)
}
