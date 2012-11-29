package main

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

type Share struct {
	Part *big.Int
	ID   int
}

func main() {
	input, err := terminal.ReadPassword(0)

	if err != nil {
		panic(err)
	}

	shared := new(big.Int)
	shared.SetBytes(input)

	fmt.Println(shared)

	parts := split(shared, 6, 3)

	fmt.Println(parts)

	newshares := []Share{parts[1], parts[2], parts[5]}

	result := join(newshares)
	fmt.Println(result)

	data := result.Bytes()

	fmt.Println(string(data))
}

func split(number *big.Int, available, needed int) []Share {
	coef := make([]*big.Int, 0)
	shares := make([]Share, 0)

	coef = append(coef, number)

	for i := 1; i < needed; i++ {
		c := big.NewInt(rand.Int63())
		coef = append(coef, c)
	}

	for x := 1; x <= available; x++ {
		accum := new(big.Int)
		accum.Set(coef[0])
		for exp := 1; exp < needed; exp++ {
			p := math.Pow(float64(x), float64(exp))
			w := new(big.Int)
			w = w.SetInt64(int64(p))

			r := new(big.Int)
			r.Mul(coef[exp], w)

			accum.Add(accum, r)
		}

		s := new(big.Int)
		s.Set(accum)
		share := Share{Part: s, ID: x}
		shares = append(shares, share)
	}

	return shares
}

func join(shares []Share) *big.Int {
	accum := new(big.Int)
	for formula := 0; formula < len(shares); formula++ {
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		value := big.NewInt(0)

		for count := 0; count < len(shares); count++ {
			if formula == count {
				continue
			}

			startposition := big.NewInt(int64(shares[formula].ID))

			value.Set(shares[formula].Part)

			nextposition := big.NewInt(int64(shares[count].ID))

			negnextpos := new(big.Int)
			negnextpos.Neg(nextposition)

			numerator.Mul(numerator, negnextpos)

			startMinNext := new(big.Int)
			startMinNext.Sub(startposition, nextposition)

			denominator.Mul(denominator, startMinNext)
		}

		fmt.Println(numerator)
		value.Mul(value, numerator).Mul(value, big.NewInt(2)).Add(value, big.NewInt(1))

		denominator.Mul(denominator, big.NewInt(2))
		value.Div(value, denominator)
		accum.Add(accum, value)
	}

	return accum.Add(accum, big.NewInt(1))
}
