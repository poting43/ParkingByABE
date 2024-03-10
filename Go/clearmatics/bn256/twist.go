package bn256

import (
	"math/big"
)

// twistPoint implements the elliptic curve y²=x³+3/ξ over GF(p²). Points are
// kept in Jacobian form and T=z² when valid. The group G₂ is the set of
// n-torsion points of this curve over GF(p²) (where n = Order)
type twistPoint struct {
	X, Y, Z, T gfP2
}

var twistB = &gfP2{
	gfP{0x38e7ecccd1dcff67, 0x65f0b37d93ce0d3e, 0xd749d0dd22ac00aa, 0x0141b9ce4a688d4d},
	gfP{0x3bf938e377b802a8, 0x020b1b273633535d, 0x26b7edf049755260, 0x2514c6324384a86d},
}

// twistGen is the generator of group G₂.
var twistGen = &twistPoint{
	gfP2{
		gfP{0xafb4737da84c6140, 0x6043dd5a5802d8c4, 0x09e950fc52a02f86, 0x14fef0833aea7b6b},
		gfP{0x8e83b5d102bc2026, 0xdceb1935497b0172, 0xfbb8264797811adf, 0x19573841af96503b},
	},
	gfP2{
		gfP{0x64095b56c71856ee, 0xdc57f922327d3cbb, 0x55f935be33351076, 0x0da4a0e693fd6482},
		gfP{0x619dfa9d886be9f6, 0xfe7fd297f59e9b78, 0xff9e1a62231b7dfe, 0x28fd7eebae9e4206},
	},
	gfP2{*newGFp(0), *newGFp(1)},
	gfP2{*newGFp(0), *newGFp(1)},
}

func (c *twistPoint) String() string {
	c.MakeAffine()
	X, Y := gfP2Decode(&c.X), gfP2Decode(&c.Y)
	return "(" + X.String() + ", " + Y.String() + ")"
}

func (c *twistPoint) GetXX() string {
	c.MakeAffine()
	X, Y := gfP2Decode(&c.X), gfP2Decode(&c.Y)
	_ = Y
	return X.X.String()
}

func (c *twistPoint) GetXY() string {
	c.MakeAffine()
	X, Y := gfP2Decode(&c.X), gfP2Decode(&c.Y)
	_ = Y
	return X.Y.String()
}

func (c *twistPoint) GetYX() string {
	c.MakeAffine()
	X, Y := gfP2Decode(&c.X), gfP2Decode(&c.Y)
	_ = X
	return Y.X.String()
}

func (c *twistPoint) GetYY() string {
	c.MakeAffine()
	X, Y := gfP2Decode(&c.X), gfP2Decode(&c.Y)
	_ = X
	return Y.Y.String()
}

func (c *twistPoint) Set(a *twistPoint) {
	c.X.Set(&a.X)
	c.Y.Set(&a.Y)
	c.Z.Set(&a.Z)
	c.T.Set(&a.T)
}

// IsOnCurve returns true iff c is on the curve.
func (c *twistPoint) IsOnCurve() bool {
	c.MakeAffine()
	if c.IsInfinity() {
		return true
	}

	y2, x3 := &gfP2{}, &gfP2{}
	y2.Square(&c.Y)
	x3.Square(&c.X).Mul(x3, &c.X).Add(x3, twistB)

	if *y2 != *x3 {
		return false
	}
	cneg := &twistPoint{}
	cneg.Mul(c, Order)
	return cneg.Z.IsZero()
}

func (c *twistPoint) SetInfinity() {
	c.X.SetZero()
	c.Y.SetOne()
	c.Z.SetZero()
	c.T.SetZero()
}

func (c *twistPoint) IsInfinity() bool {
	return c.Z.IsZero()
}

func (c *twistPoint) Add(a, b *twistPoint) {
	// For additional comments, see the same function in curve.go.

	if a.IsInfinity() {
		c.Set(b)
		return
	}
	if b.IsInfinity() {
		c.Set(a)
		return
	}

	// See http://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian-0/addition/add-2007-bl.op3
	z12 := (&gfP2{}).Square(&a.Z)
	z22 := (&gfP2{}).Square(&b.Z)
	u1 := (&gfP2{}).Mul(&a.X, z22)
	u2 := (&gfP2{}).Mul(&b.X, z12)

	T := (&gfP2{}).Mul(&b.Z, z22)
	s1 := (&gfP2{}).Mul(&a.Y, T)

	T.Mul(&a.Z, z12)
	s2 := (&gfP2{}).Mul(&b.Y, T)

	h := (&gfP2{}).Sub(u2, u1)
	xEqual := h.IsZero()

	T.Add(h, h)
	i := (&gfP2{}).Square(T)
	j := (&gfP2{}).Mul(h, i)

	T.Sub(s2, s1)
	yEqual := T.IsZero()
	if xEqual && yEqual {
		c.Double(a)
		return
	}
	r := (&gfP2{}).Add(T, T)

	v := (&gfP2{}).Mul(u1, i)

	t4 := (&gfP2{}).Square(r)
	T.Add(v, v)
	t6 := (&gfP2{}).Sub(t4, j)
	c.X.Sub(t6, T)

	T.Sub(v, &c.X) // t7
	t4.Mul(s1, j)  // t8
	t6.Add(t4, t4) // t9
	t4.Mul(r, T)   // t10
	c.Y.Sub(t4, t6)

	T.Add(&a.Z, &b.Z) // t11
	t4.Square(T)      // t12
	T.Sub(t4, z12)    // t13
	t4.Sub(T, z22)    // t14
	c.Z.Mul(t4, h)
}

func (c *twistPoint) Double(a *twistPoint) {
	// See http://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian-0/doubling/dbl-2009-l.op3
	A := (&gfP2{}).Square(&a.X)
	B := (&gfP2{}).Square(&a.Y)
	C := (&gfP2{}).Square(B)

	T := (&gfP2{}).Add(&a.X, B)
	t2 := (&gfP2{}).Square(T)
	T.Sub(t2, A)
	t2.Sub(T, C)
	d := (&gfP2{}).Add(t2, t2)
	T.Add(A, A)
	e := (&gfP2{}).Add(T, A)
	f := (&gfP2{}).Square(e)

	T.Add(d, d)
	c.X.Sub(f, T)

	T.Add(C, C)
	t2.Add(T, T)
	T.Add(t2, t2)
	c.Y.Sub(d, &c.X)
	t2.Mul(e, &c.Y)
	c.Y.Sub(t2, T)

	T.Mul(&a.Y, &a.Z)
	c.Z.Add(T, T)
}

func (c *twistPoint) Mul(a *twistPoint, scalar *big.Int) {
	sum, T := &twistPoint{}, &twistPoint{}

	for i := scalar.BitLen(); i >= 0; i-- {
		T.Double(sum)
		if scalar.Bit(i) != 0 {
			sum.Add(T, a)
		} else {
			sum.Set(T)
		}
	}

	c.Set(sum)
}

func (c *twistPoint) MakeAffine() {
	if c.Z.IsOne() {
		return
	} else if c.Z.IsZero() {
		c.X.SetZero()
		c.Y.SetOne()
		c.T.SetZero()
		return
	}

	zInv := (&gfP2{}).Invert(&c.Z)
	T := (&gfP2{}).Mul(&c.Y, zInv)
	zInv2 := (&gfP2{}).Square(zInv)
	c.Y.Mul(T, zInv2)
	T.Mul(&c.X, zInv2)
	c.X.Set(T)
	c.Z.SetOne()
	c.T.SetOne()
}

func (c *twistPoint) Neg(a *twistPoint) {
	c.X.Set(&a.X)
	c.Y.Neg(&a.Y)
	c.Z.Set(&a.Z)
	c.T.SetZero()
}
