package jcrypt

import (
	"fmt"
	"math/rand"
	"strings"
)

var family = []string{
	"QzF3n6/9CAtpu0O",
	"B1IREhcSyrleKvMW8LXx",
	"7N-dVbwsY2g4oaJZGUDj",
	"iHkq.mPf5T",
}

var magic = "$9$"

var encoding = [][]int{
	{
		1, 4, 32,
	},
	{
		1, 16, 32,
	},
	{
		1, 8, 32,
	},
	{
		1, 64,
	},
	{
		1, 32,
	},
	{
		1, 4, 16, 128,
	},
	{
		1, 32, 64,
	},
}

var numAlpha map[int]rune
var alphaNum map[rune]int

var extra map[rune]int

func init() {
	numAlpha = make(map[int]rune)
	alphaNum = make(map[rune]int)
	extra = make(map[rune]int)
	for i, r := range []rune(strings.Join(family, "")) {
		numAlpha[i] = r
		alphaNum[r] = i
	}

	for i, fam := range family {
		for _, c := range []rune(fam) {
			extra[c] = 3 - i
		}
	}

}

func Encrypt(password string, salt rune) string {

	rand := randc(extra[salt])
	prev := salt
	crypt := fmt.Sprintf("%s%c%s", magic, salt, rand)

	for pos, c := range []rune(password) {
		encode := encoding[pos%len(encoding)]
		crypt += gapEncode(c, prev, encode)
		x := []rune(crypt)
		prev = x[len(x)-1]
	}

	return crypt
}

func gapEncode(pc rune, prev rune, enc []int) string {
	crypt := ""
	gaps := make([]int, 0)

	for _, mod := range reverse(enc) {
		gaps = append([]int{int(pc) / mod}, gaps...)
		pc = pc % rune(mod)
	}

	for _, gap := range gaps {
		gap += alphaNum[prev] + 1
		c := numAlpha[gap%len(numAlpha)]
		prev = c
		crypt += string(c)
	}

	return crypt
}

func randc(cnt int) string {
	r := ""
	randNumGen := rand.New(rand.NewSource(99))

	for cnt > 0 {
		r += string(numAlpha[int(randNumGen.Int63())%len(numAlpha)])
		cnt--
	}

	return r
}

func reverse(in []int) []int {
	out := make([]int, len(in))
	j := 0
	for i := len(in) - 1; i >= 0; i-- {
		out[j] = in[i]
		j++
	}
	return out
}
