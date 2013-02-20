package main

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh/terminal"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"
)

type Share struct {
	Part *big.Int
	ID   int64
}

func ShareToBase64(share Share) string {
	buff := new(bytes.Buffer)
	buff.WriteByte(byte(share.ID))
	buff.Write(share.Part.Bytes())

	return base64.StdEncoding.EncodeToString(buff.Bytes())
}

func Base64ToShare(share string) Share {
	decoded, _ := base64.StdEncoding.DecodeString(share)
	buff := bytes.NewBuffer(decoded).Bytes()
	id := int64(buff[0])
	b := buff[1:]

	part := new(big.Int)
	part.SetBytes(b)
	return Share{ID: id, Part: part}
}

func EasySplit(secret string, available, needed int) []string {
	shared := new(big.Int)
	shared.SetBytes([]byte(secret))
	parts := split(shared, available, needed)

	partsEncoded := make([]string, 0)
	for _, p := range parts {
		partsEncoded = append(partsEncoded, ShareToBase64(p))
	}

	return partsEncoded
}

func EasyJoin(parts []string) string {
	shares := make([]Share, 0)
	for _, v := range parts {
		shares = append(shares, Base64ToShare(v))
	}

	result := join(shares)
	data := result.Bytes()
	return string(data)
}

func main() {
	input, err := terminal.ReadPassword(0)

	if err != nil {
		panic(err)
	}

	secret := string(input)
	parts := EasySplit(secret, 6, 3)

	fmt.Println(parts)

	newshares := []string{parts[1], parts[2], parts[5]}
	result := EasyJoin(newshares)

	fmt.Println(result)
}

func split(number *big.Int, available, needed int) []Share {
	coef := make([]*big.Int, 0)
	shares := make([]Share, 0)

	coef = append(coef, number)

	rand.Seed(time.Now().Unix())
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
	zero := big.NewInt(0)
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
			denominator.Mul(denominator,
				startposition.Sub(startposition, nextposition))

		}

		value.Mul(value, numerator).Mul(value, big.NewInt(2))

		if denominator.Cmp(zero) < 0 {
			value.Add(value, big.NewInt(2))
		} else {
			value.Add(value, big.NewInt(1))
		}

		denominator.Mul(denominator, big.NewInt(2))

		value.Div(value, denominator)
		accum.Add(accum, value)
	}

	return accum
}
