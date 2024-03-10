// Package bn256 implements a particular bilinear group at the 128-bit security
// level.
//
// Bilinear groups are the basis of many of the new cryptographic protocols that
// have been proposed over the past decade. They consist of a triplet of groups
// (G₁, G₂ and GT) such that there exists a function e(g₁ˣ,g₂ʸ)=gTˣʸ (where gₓ
// is a generator of the respective group). That function is called a pairing
// function.
//
// This package specifically implements the Optimal Ate pairing over a 256-bit
// Barreto-Naehrig curve as described in
// http://cryptojedi.org/papers/dclxvi-20100714.pdf. Its output is compatible
// with the implementation described in that paper.
package bn256

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

func randomK(r io.Reader) (k *big.Int, err error) {
	for {
		k, err = rand.Int(r, Order)
		if k.Sign() > 0 || err != nil {
			return
		}
	}
}

// G1 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G1 struct {
	P *curvePoint
}

// RandomG1 returns X and g₁ˣ where X is a random, non-zero number read from r.
func RandomG1(r io.Reader) (*big.Int, *G1, error) {
	k, err := randomK(r)
	if err != nil {
		return nil, nil, err
	}

	return k, new(G1).ScalarBaseMult(k), nil
}

func (g *G1) String() string {
	return "bn256.G1" + g.P.String()
}

// ScalarBaseMult sets e to g*k where g is the generator of the group and then
// returns e.
func (e *G1) ScalarBaseMult(k *big.Int) *G1 {
	if e.P == nil {
		e.P = &curvePoint{}
	}
	e.P.Mul(curveGen, k)
	return e
}

// ScalarMult sets e to a*k and then returns e.
func (e *G1) ScalarMult(a *G1, k *big.Int) *G1 {
	if e.P == nil {
		e.P = &curvePoint{}
	}
	e.P.Mul(a.P, k)
	return e
}

// Add sets e to a+b and then returns e.
func (e *G1) Add(a, b *G1) *G1 {
	if e.P == nil {
		e.P = &curvePoint{}
	}
	e.P.Add(a.P, b.P)
	return e
}

// Neg sets e to -a and then returns e.
func (e *G1) Neg(a *G1) *G1 {
	if e.P == nil {
		e.P = &curvePoint{}
	}
	e.P.Neg(a.P)
	return e
}

// Set sets e to a and then returns e.
func (e *G1) Set(a *G1) *G1 {
	if e.P == nil {
		e.P = &curvePoint{}
	}
	e.P.Set(a.P)
	return e
}

// Marshal converts e to a byte slice.
func (e *G1) Marshal() []byte {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8

	e.P.MakeAffine()
	ret := make([]byte, numBytes*2)
	if e.P.IsInfinity() {
		return ret
	}
	temp := &gfP{}

	montDecode(temp, &e.P.X)
	temp.Marshal(ret)
	montDecode(temp, &e.P.Y)
	temp.Marshal(ret[numBytes:])

	return ret
}

// Unmarshal sets e to the result of converting the output of Marshal back into
// a group element and then returns e.
func (e *G1) Unmarshal(m []byte) ([]byte, error) {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8
	if len(m) < 2*numBytes {
		return nil, errors.New("bn256: not enough data")
	}
	// Unmarshal the points and check their caps
	if e.P == nil {
		e.P = &curvePoint{}
	} else {
		e.P.X, e.P.Y = gfP{0}, gfP{0}
	}
	var err error
	if err = e.P.X.Unmarshal(m); err != nil {
		return nil, err
	}
	if err = e.P.Y.Unmarshal(m[numBytes:]); err != nil {
		return nil, err
	}
	// Encode into Montgomery form and ensure it's on the curve
	montEncode(&e.P.X, &e.P.X)
	montEncode(&e.P.Y, &e.P.Y)

	zero := gfP{0}
	if e.P.X == zero && e.P.Y == zero {
		// This is the point at infinity.
		e.P.Y = *newGFp(1)
		e.P.Z = gfP{0}
		e.P.T = gfP{0}
	} else {
		e.P.Z = *newGFp(1)
		e.P.T = *newGFp(1)

		if !e.P.IsOnCurve() {
			return nil, errors.New("bn256: malformed point")
		}
	}
	return m[2*numBytes:], nil
}

// G2 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G2 struct {
	P *twistPoint
}

// RandomG2 returns X and g₂ˣ where X is a random, non-zero number read from r.
func RandomG2(r io.Reader) (*big.Int, *G2, error) {
	k, err := randomK(r)
	if err != nil {
		return nil, nil, err
	}

	return k, new(G2).ScalarBaseMult(k), nil
}

func (e *G2) String() string {
	return "bn256.G2" + e.P.String()
}

// ScalarBaseMult sets e to g*k where g is the generator of the group and then
// returns out.
func (e *G2) ScalarBaseMult(k *big.Int) *G2 {
	if e.P == nil {
		e.P = &twistPoint{}
	}
	e.P.Mul(twistGen, k)
	return e
}

// ScalarMult sets e to a*k and then returns e.
func (e *G2) ScalarMult(a *G2, k *big.Int) *G2 {
	if e.P == nil {
		e.P = &twistPoint{}
	}
	e.P.Mul(a.P, k)
	return e
}

// Add sets e to a+b and then returns e.
func (e *G2) Add(a, b *G2) *G2 {
	if e.P == nil {
		e.P = &twistPoint{}
	}
	e.P.Add(a.P, b.P)
	return e
}

// Neg sets e to -a and then returns e.
func (e *G2) Neg(a *G2) *G2 {
	if e.P == nil {
		e.P = &twistPoint{}
	}
	e.P.Neg(a.P)
	return e
}

// Set sets e to a and then returns e.
func (e *G2) Set(a *G2) *G2 {
	if e.P == nil {
		e.P = &twistPoint{}
	}
	e.P.Set(a.P)
	return e
}

// Marshal converts e into a byte slice.
func (e *G2) Marshal() []byte {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8

	if e.P == nil {
		e.P = &twistPoint{}
	}

	e.P.MakeAffine()
	ret := make([]byte, numBytes*4)
	if e.P.IsInfinity() {
		return ret
	}
	temp := &gfP{}

	montDecode(temp, &e.P.X.X)
	temp.Marshal(ret)
	montDecode(temp, &e.P.X.Y)
	temp.Marshal(ret[numBytes:])
	montDecode(temp, &e.P.Y.X)
	temp.Marshal(ret[2*numBytes:])
	montDecode(temp, &e.P.Y.Y)
	temp.Marshal(ret[3*numBytes:])

	return ret
}

// Unmarshal sets e to the result of converting the output of Marshal back into
// a group element and then returns e.
func (e *G2) Unmarshal(m []byte) ([]byte, error) {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8
	if len(m) < 4*numBytes {
		return nil, errors.New("bn256: not enough data")
	}
	// Unmarshal the points and check their caps
	if e.P == nil {
		e.P = &twistPoint{}
	}
	var err error
	if err = e.P.X.X.Unmarshal(m); err != nil {
		return nil, err
	}
	if err = e.P.X.Y.Unmarshal(m[numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.X.Unmarshal(m[2*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.Y.Unmarshal(m[3*numBytes:]); err != nil {
		return nil, err
	}
	// Encode into Montgomery form and ensure it's on the curve
	montEncode(&e.P.X.X, &e.P.X.X)
	montEncode(&e.P.X.Y, &e.P.X.Y)
	montEncode(&e.P.Y.X, &e.P.Y.X)
	montEncode(&e.P.Y.Y, &e.P.Y.Y)

	if e.P.X.IsZero() && e.P.Y.IsZero() {
		// This is the point at infinity.
		e.P.Y.SetOne()
		e.P.Z.SetZero()
		e.P.T.SetZero()
	} else {
		e.P.Z.SetOne()
		e.P.T.SetOne()

		if !e.P.IsOnCurve() {
			return nil, errors.New("bn256: malformed point")
		}
	}
	return m[4*numBytes:], nil
}

// GT is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type GT struct {
	P *gfP12
}

// Pair calculates an Optimal Ate pairing.
func Pair(g1 *G1, g2 *G2) *GT {
	return &GT{optimalAte(g2.P, g1.P)}
}

// PairingCheck calculates the Optimal Ate pairing for a set of points.
func PairingCheck(a []*G1, b []*G2) bool {
	acc := new(gfP12)
	acc.SetOne()

	for i := 0; i < len(a); i++ {
		if a[i].P.IsInfinity() || b[i].P.IsInfinity() {
			continue
		}
		acc.Mul(acc, miller(b[i].P, a[i].P))
	}
	return finalExponentiation(acc).IsOne()
}

// Miller applies Miller's algorithm, which is a bilinear function from the
// source groups to F_p^12. Miller(g1, g2).Finalize() is equivalent to Pair(g1,
// g2).
func Miller(g1 *G1, g2 *G2) *GT {
	return &GT{miller(g2.P, g1.P)}
}

func (g *GT) String() string {
	return "bn256.GT" + g.P.String()
}

// ScalarMult sets e to a*k and then returns e.
func (e *GT) ScalarMult(a *GT, k *big.Int) *GT {
	if e.P == nil {
		e.P = &gfP12{}
	}
	e.P.Exp(a.P, k)
	return e
}

// Add sets e to a+b and then returns e.
func (e *GT) Add(a, b *GT) *GT {
	if e.P == nil {
		e.P = &gfP12{}
	}
	e.P.Mul(a.P, b.P)
	return e
}

// Neg sets e to -a and then returns e.
func (e *GT) Neg(a *GT) *GT {
	if e.P == nil {
		e.P = &gfP12{}
	}
	e.P.Conjugate(a.P)
	return e
}

// Set sets e to a and then returns e.
func (e *GT) Set(a *GT) *GT {
	if e.P == nil {
		e.P = &gfP12{}
	}
	e.P.Set(a.P)
	return e
}

// Finalize is a linear function from F_p^12 to GT.
func (e *GT) Finalize() *GT {
	ret := finalExponentiation(e.P)
	e.P.Set(ret)
	return e
}

// Marshal converts e into a byte slice.
func (e *GT) Marshal() []byte {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8

	ret := make([]byte, numBytes*12)
	temp := &gfP{}

	montDecode(temp, &e.P.X.X.X)
	temp.Marshal(ret)
	montDecode(temp, &e.P.X.X.Y)
	temp.Marshal(ret[numBytes:])
	montDecode(temp, &e.P.X.Y.X)
	temp.Marshal(ret[2*numBytes:])
	montDecode(temp, &e.P.X.Y.Y)
	temp.Marshal(ret[3*numBytes:])
	montDecode(temp, &e.P.X.Z.X)
	temp.Marshal(ret[4*numBytes:])
	montDecode(temp, &e.P.X.Z.Y)
	temp.Marshal(ret[5*numBytes:])
	montDecode(temp, &e.P.Y.X.X)
	temp.Marshal(ret[6*numBytes:])
	montDecode(temp, &e.P.Y.X.Y)
	temp.Marshal(ret[7*numBytes:])
	montDecode(temp, &e.P.Y.Y.X)
	temp.Marshal(ret[8*numBytes:])
	montDecode(temp, &e.P.Y.Y.Y)
	temp.Marshal(ret[9*numBytes:])
	montDecode(temp, &e.P.Y.Z.X)
	temp.Marshal(ret[10*numBytes:])
	montDecode(temp, &e.P.Y.Z.Y)
	temp.Marshal(ret[11*numBytes:])

	return ret
}

// Unmarshal sets e to the result of converting the output of Marshal back into
// a group element and then returns e.
func (e *GT) Unmarshal(m []byte) ([]byte, error) {
	// Each value is a 256-bit number.
	const numBytes = 256 / 8

	if len(m) < 12*numBytes {
		return nil, errors.New("bn256: not enough data")
	}

	if e.P == nil {
		e.P = &gfP12{}
	}

	var err error
	if err = e.P.X.X.X.Unmarshal(m); err != nil {
		return nil, err
	}
	if err = e.P.X.X.Y.Unmarshal(m[numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.X.Y.X.Unmarshal(m[2*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.X.Y.Y.Unmarshal(m[3*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.X.Z.X.Unmarshal(m[4*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.X.Z.Y.Unmarshal(m[5*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.X.X.Unmarshal(m[6*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.X.Y.Unmarshal(m[7*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.Y.X.Unmarshal(m[8*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.Y.Y.Unmarshal(m[9*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.Z.X.Unmarshal(m[10*numBytes:]); err != nil {
		return nil, err
	}
	if err = e.P.Y.Z.Y.Unmarshal(m[11*numBytes:]); err != nil {
		return nil, err
	}
	montEncode(&e.P.X.X.X, &e.P.X.X.X)
	montEncode(&e.P.X.X.Y, &e.P.X.X.Y)
	montEncode(&e.P.X.Y.X, &e.P.X.Y.X)
	montEncode(&e.P.X.Y.Y, &e.P.X.Y.Y)
	montEncode(&e.P.X.Z.X, &e.P.X.Z.X)
	montEncode(&e.P.X.Z.Y, &e.P.X.Z.Y)
	montEncode(&e.P.Y.X.X, &e.P.Y.X.X)
	montEncode(&e.P.Y.X.Y, &e.P.Y.X.Y)
	montEncode(&e.P.Y.Y.X, &e.P.Y.Y.X)
	montEncode(&e.P.Y.Y.Y, &e.P.Y.Y.Y)
	montEncode(&e.P.Y.Z.X, &e.P.Y.Z.X)
	montEncode(&e.P.Y.Z.Y, &e.P.Y.Z.Y)

	return m[12*numBytes:], nil
}
