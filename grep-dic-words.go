package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	NRANKS = 10
	usage  = `
Usage: grep-dic-words FILE WORD
FILE is *.dic holding an ispell dictionary.
WORD is the word to match.
The output are the top words from the dictionary that match characters in WORD.

Example:
  $ grep-dic-words.go en_GB.dic egtoyz
  benzoylmethylecgonine 6
  cyberorganization 6
  demythologize 6
  egyptianization 6
  ethnozoology 6
  etymologize 6
  heterozygote 6
  heterozygous 6
  heterozygousness 6
  homozygote 6
These are the words from en_GB.dic that match the most characters with "egtoyz".
Coincidentally they match all 6 chars.
`
)

func main() {
	check := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if len(os.Args) != 3 {
		check(errors.New(usage))
	}
	ch, err := source(os.Args[1])
	check(err)
	r := newRank()
	for s := range ch {
		r.register(s, score(s, os.Args[2]))
	}
	fmt.Println(r)
}

func source(fname string) (ch chan string, err error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	ch = make(chan string)
	scanner := bufio.NewScanner(file)

	go func() {
		for scanner.Scan() {
			s := strings.ToLower(scanner.Text())
			for _, splitter := range []string{"/", " ", "\t"} {
				if strings.Contains(s, splitter) {
					s = strings.Split(s, splitter)[0]
				}
			}
			ch <- s
		}
		close(ch)
	}()
	return ch, nil
}

func score(s, ref string) (sc int) {
	sc = 0
	for _, ch := range ref {
		if strings.Contains(s, string(ch)) {
			sc++
		}
	}
	return sc
}

type entry struct {
	s string
	m int
}
type ranking struct {
	r []*entry
}

func newRank() *ranking {
	return &ranking{
		r: []*entry{},
	}
}

func (r *ranking) register(s string, m int) {
	e := &entry{s: s, m: m}
	if len(r.r) < NRANKS {
		r.r = append(r.r, e)
	} else {
		if r.r[NRANKS-1].m < m {
			r.r[NRANKS-1] = e
		}
	}
	sort.Slice(r.r, func(i, j int) bool {
		return r.r[i].m > r.r[j].m
	})
}

func (r *ranking) String() string {
	s := []string{}
	for _, e := range r.r {
		s = append(s, fmt.Sprintf("%s %d", e.s, e.m))
	}
	return strings.Join(s, "\n")
}
