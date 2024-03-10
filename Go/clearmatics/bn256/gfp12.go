package bn256

// For details of the algorithms used, see "Multiplication and Squaring on
// Pairing-Friendly Fields, Devegili et al.
// http://eprint.iacr.org/2006/471.pdf.

import (
	"math/big"
)

// gfP12 implements the field of size p¹² as a quadratic extension of gfP6
// where ω²=τ.
type gfP12 struct {
	X, Y gfP6 // value is xω + Y
}

func (e *gfP12) String() string {
	return "(" + e.X.String() + "," + e.Y.String() + ")"
}

func (e *gfP12) Set(a *gfP12) *gfP12 {
	e.X.Set(&a.X)
	e.Y.Set(&a.Y)
	return e
}

func (e *gfP12) SetZero() *gfP12 {
	e.X.SetZero()
	e.Y.SetZero()
	return e
}

func (e *gfP12) SetOne() *gfP12 {
	e.X.SetZero()
	e.Y.SetOne()
	return e
}

func (e *gfP12) IsZero() bool {
	return e.X.IsZero() && e.Y.IsZero()
}

func (e *gfP12) IsOne() bool {
	return e.X.IsZero() && e.Y.IsOne()
}

func (e *gfP12) Conjugate(a *gfP12) *gfP12 {
	e.X.Neg(&a.X)
	e.Y.Set(&a.Y)
	return e
}

func (e *gfP12) Neg(a *gfP12) *gfP12 {
	e.X.Neg(&a.X)
	e.Y.Neg(&a.Y)
	return e
}

// Frobenius computes (xω+Y)^p = X^p ω·ξ^((p-1)/6) + Y^p
func (e *gfP12) Frobenius(a *gfP12) *gfP12 {
	e.X.Frobenius(&a.X)
	e.Y.Frobenius(&a.Y)
	e.X.MulScalar(&e.X, xiToPMinus1Over6)
	return e
}

// FrobeniusP2 computes (xω+Y)^p² = X^p² ω·ξ^((p²-1)/6) + Y^p²
func (e *gfP12) FrobeniusP2(a *gfP12) *gfP12 {
	e.X.FrobeniusP2(&a.X)
	e.X.MulGFP(&e.X, xiToPSquaredMinus1Over6)
	e.Y.FrobeniusP2(&a.Y)
	return e
}

func (e *gfP12) FrobeniusP4(a *gfP12) *gfP12 {
	e.X.FrobeniusP4(&a.X)
	e.X.MulGFP(&e.X, xiToPSquaredMinus1Over3)
	e.Y.FrobeniusP4(&a.Y)
	return e
}

func (e *gfP12) Add(a, b *gfP12) *gfP12 {
	e.X.Add(&a.X, &b.X)
	e.Y.Add(&a.Y, &b.Y)
	return e
}

func (e *gfP12) Sub(a, b *gfP12) *gfP12 {
	e.X.Sub(&a.X, &b.X)
	e.Y.Sub(&a.Y, &b.Y)
	return e
}

func (e *gfP12) Mul(a, b *gfP12) *gfP12 {
	tx := (&gfP6{}).Mul(&a.X, &b.Y)
	t := (&gfP6{}).Mul(&b.X, &a.Y)
	tx.Add(tx, t)

	ty := (&gfP6{}).Mul(&a.Y, &b.Y)
	t.Mul(&a.X, &b.X).MulTau(t)

	e.X.Set(tx)
	e.Y.Add(ty, t)
	return e
}

func (e *gfP12) MulScalar(a *gfP12, b *gfP6) *gfP12 {
	e.X.Mul(&e.X, b)
	e.Y.Mul(&e.Y, b)
	return e
}

func (c *gfP12) Exp(a *gfP12, power *big.Int) *gfP12 {
	sum := (&gfP12{}).SetOne()
	t := &gfP12{}

	for i := power.BitLen() - 1; i >= 0; i-- {
		t.Square(sum)
		if power.Bit(i) != 0 {
			sum.Mul(t, a)
		} else {
			sum.Set(t)
		}
	}

	c.Set(sum)
	return c
}

func (e *gfP12) Square(a *gfP12) *gfP12 {
	// Complex squaring algorithm
	v0 := (&gfP6{}).Mul(&a.X, &a.Y)

	t := (&gfP6{}).MulTau(&a.X)
	t.Add(&a.Y, t)
	ty := (&gfP6{}).Add(&a.X, &a.Y)
	ty.Mul(ty, t).Sub(ty, v0)
	t.MulTau(v0)
	ty.Sub(ty, t)

	e.X.Add(v0, v0)
	e.Y.Set(ty)
	return e
}

func (e *gfP12) Invert(a *gfP12) *gfP12 {
	// See "Implementing cryptographic pairings", M. Scott, section 3.2.
	// ftp://136.206.11.249/pub/crypto/pairings.pdf
	t1, t2 := &gfP6{}, &gfP6{}

	t1.Square(&a.X)
	t2.Square(&a.Y)
	t1.MulTau(t1).Sub(t2, t1)
	t2.Invert(t1)

	e.X.Neg(&a.X)
	e.Y.Set(&a.Y)
	e.MulScalar(e, t2)
	return e
}
