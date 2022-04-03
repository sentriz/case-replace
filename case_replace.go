package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const sep = "/"

func main() {
	if len(os.Args) < 2 || len(os.Args)%2 != 1 {
		log.Fatalf("usage: %s from/pattern to/pattern from/pattern to/pattern ...", os.Args[0])
	}

	pairs, err := parsePairs(os.Args[1:])
	if err != nil {
		log.Fatalf("error parsing pairs: %v", err)
	}

	replacer, err := makeReplacer(pairs, casers)
	if err != nil {
		log.Fatalf("error making replacer: %v", err)
	}

	input, _ := io.ReadAll(os.Stdin)
	replaced := replacer.Replace(string(input))
	os.Stdout.Write([]byte(replaced))
}

type pair struct {
	from, to []string
}

func parsePairs(in []string) ([]pair, error) {
	var pairs []pair
	for i := 0; i < len(in)-1; i += 2 {
		from, to := in[i], in[i+1]
		fromParts := strings.Split(from, sep)
		toParts := strings.Split(to, sep)
		if len(fromParts) == 0 || len(toParts) == 0 {
			return nil, fmt.Errorf("empty pair")
		}
		pairs = append(pairs, pair{fromParts, toParts})
	}
	return pairs, nil
}

func makeReplacer(pairs []pair, casers []caser) (*strings.Replacer, error) {
	var replacements []string
	for _, pair := range pairs {
		for _, csr := range casers {
			replacements = append(replacements, csr(pair.from))
			replacements = append(replacements, csr(pair.to))
		}
	}
	return strings.NewReplacer(replacements...), nil
}

type caser func([]string) string

var casers = []caser{
	// pascal
	func(in []string) string { return maps(in, "", strings.Title) },
	// camel
	func(in []string) string { return in[0] + maps(in[1:], "", strings.Title) },
	// snake
	func(in []string) string { return maps(in, "_", strings.ToLower) },
	// kebab
	func(in []string) string { return maps(in, "-", strings.ToLower) },
	// screaming
	func(in []string) string { return maps(in, "_", strings.ToUpper) },
	// joined lower
	func(in []string) string { return maps(in, "", strings.ToLower) },
	// joined upper
	func(in []string) string { return maps(in, "", strings.ToUpper) },
	// space
	func(in []string) string { return maps(in, " ", strings.ToLower) },
}

func maps(in []string, sep string, f func(string) string) string {
	var out []string
	for _, i := range in {
		out = append(out, f(i))
	}
	return strings.Join(out, sep)
}
