package main

import (
	"testing"
)

var testPairs = []pair{
	{from: []string{"one", "two", "three"}, to: []string{"four", "five", "six"}},
	{from: []string{"cat", "sat"}, to: []string{"small", "dog", "dog"}},
}

func TestNoReplace(t *testing.T) {
	check(t, `one`, `one`)
	check(t, `two`, `two`)
	check(t, `three`, `three`)
	check(t, ``, ``)
}

func TestReplace(t *testing.T) {
	check(t,
		`
			oneTwoThree
			OneTwoThree
			one_two_three
			one-two-three
			ONE_TWO_THREE
			onetwothree
			ONETWOTHREE
			one two three
		`,
		`
			fourFiveSix
			FourFiveSix
			four_five_six
			four-five-six
			FOUR_FIVE_SIX
			fourfivesix
			FOURFIVESIX
			four five six
		`,
	)
}

func TestMulti(t *testing.T) {
	check(t, `oneTwoThree catSat`, `fourFiveSix smallDogDog`)
	check(t, `one_two_three cat-sat`, `four_five_six small-dog-dog`)
}

func check(t *testing.T, in, exp string) {
	replacer, err := makeReplacer(testPairs, casers)
	if err != nil {
		t.Fatalf("error making replacer: %v", err)
	}

	if out := replacer.Replace(in); out != exp {
		t.Fatalf("\n> in\n%s\n> exp\n%s\n> out\n%s", in, exp, out)
	}
}
