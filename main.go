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
	ID   int64
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
			w := big.NewInt(int64(p))

			r := new(big.Int)
			r.Mul(coef[exp], w)

			accum.Add(accum, r)
		}

		s := new(big.Int)
		s.Set(accum)
		share := Share{Part: s, ID: int64(x)}
		shares = append(shares, share)
	}

	return shares
}

func join(shares []Share) *big.Int {
	accum := big.NewInt(1)
	for formula := 0; formula < len(shares); formula++ {
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		value := new(big.Int)
		value.Set(shares[formula].Part)

		for count := 0; count < len(shares); count++ {
			if formula == count {
				continue
			}

			startposition := big.NewInt(shares[formula].ID)
			nextposition := big.NewInt(shares[count].ID)
			negnextpos := new(big.Int)
			negnextpos.Neg(nextposition)

			numerator.Mul(numerator, negnextpos)
			denominator.Mul(denominator, startposition.Sub(startposition, nextposition))
		}

		value.Mul(value, numerator).Mul(value, big.NewInt(2)).Add(value, big.NewInt(1))
		denominator.Mul(denominator, big.NewInt(2))

		value.Div(value, denominator)
		accum.Add(accum, value)
	}

	return accum
}
