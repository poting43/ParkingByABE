package bn256

import (
	"math/big"
)

// curvePoint implements the elliptic curve y²=x³+3. Points are kept in Jacobian
// form and T=z² when valid. G₁ is the set of points of this curve on GF(p).
type curvePoint struct {
	X, Y, Z, T gfP
}

var curveB = newGFp(3)

// curveGen is the generator of G₁.
var curveGen = &curvePoint{
	X: *newGFp(1),
	Y: *newGFp(2),
	Z: *newGFp(1),
	T: *newGFp(1),
}

func (c *curvePoint) String() string {
	c.MakeAffine()
	X, Y := &gfP{}, &gfP{}
	montDecode(X, &c.X)
	montDecode(Y, &c.Y)
	return "(" + X.String() + ", " + Y.String() + ")"
}

func (c *curvePoint) GetX() string {
	c.MakeAffine()
	X, Y := &gfP{}, &gfP{}
	montDecode(X, &c.X)
	montDecode(Y, &c.Y)
	return X.String()
}

func (c *curvePoint) GetY() string {
	c.MakeAffine()
	X, Y := &gfP{}, &gfP{}
	montDecode(X, &c.X)
	montDecode(Y, &c.Y)
	return Y.String()
}

func (c *curvePoint) Set(a *curvePoint) {
	c.X.Set(&a.X)
	c.Y.Set(&a.Y)
	c.Z.Set(&a.Z)
	c.T.Set(&a.T)
}

// IsOnCurve returns true iff c is on the curve.
func (c *curvePoint) IsOnCurve() bool {
	c.MakeAffine()
	if c.IsInfinity() {
		return true
	}

	y2, x3 := &gfP{}, &gfP{}
	gfpMul(y2, &c.Y, &c.Y)
	gfpMul(x3, &c.X, &c.X)
	gfpMul(x3, x3, &c.X)
	gfpAdd(x3, x3, curveB)

	return *y2 == *x3
}

func (c *curvePoint) SetInfinity() {
	c.X = gfP{0}
	c.Y = *newGFp(1)
	c.Z = gfP{0}
	c.T = gfP{0}
}

func (c *curvePoint) IsInfinity() bool {
	return c.Z == gfP{0}
}

func (c *curvePoint) Add(a, b *curvePoint) {
	if a.IsInfinity() {
		c.Set(b)
		return
	}
	if b.IsInfinity() {
		c.Set(a)
		return
	}

	// See http://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian-0/addition/add-2007-bl.op3

	// Normalize the points by replacing a = [x1:y1:z1] and b = [x2:y2:z2]
	// by [u1:s1:z1·z2] and [u2:s2:z1·z2]
	// where u1 = x1·z2², s1 = y1·z2³ and u1 = x2·z1², s2 = y2·z1³
	z12, z22 := &gfP{}, &gfP{}
	gfpMul(z12, &a.Z, &a.Z)
	gfpMul(z22, &b.Z, &b.Z)

	u1, u2 := &gfP{}, &gfP{}
	gfpMul(u1, &a.X, z22)
	gfpMul(u2, &b.X, z12)

	T, s1 := &gfP{}, &gfP{}
	gfpMul(T, &b.Z, z22)
	gfpMul(s1, &a.Y, T)

	s2 := &gfP{}
	gfpMul(T, &a.Z, z12)
	gfpMul(s2, &b.Y, T)

	// Compute X = (2h)²(s²-u1-u2)
	// where s = (s2-s1)/(u2-u1) is the slope of the line through
	// (u1,s1) and (u2,s2). The extra factor 2h = 2(u2-u1) comes from the value of Z below.
	// This is also:
	// 4(s2-s1)² - 4h²(u1+u2) = 4(s2-s1)² - 4h³ - 4h²(2u1)
	//                        = r² - j - 2v
	// with the notations below.
	h := &gfP{}
	gfpSub(h, u2, u1)
	xEqual := *h == gfP{0}

	gfpAdd(T, h, h)
	// i = 4h²
	i := &gfP{}
	gfpMul(i, T, T)
	// j = 4h³
	j := &gfP{}
	gfpMul(j, h, i)

	gfpSub(T, s2, s1)
	yEqual := *T == gfP{0}
	if xEqual && yEqual {
		c.Double(a)
		return
	}
	r := &gfP{}
	gfpAdd(r, T, T)

	v := &gfP{}
	gfpMul(v, u1, i)

	// t4 = 4(s2-s1)²
	t4, t6 := &gfP{}, &gfP{}
	gfpMul(t4, r, r)
	gfpAdd(T, v, v)
	gfpSub(t6, t4, j)

	gfpSub(&c.X, t6, T)

	// Set Y = -(2h)³(s1 + s*(X/4h²-u1))
	// This is also
	// Y = - 2·s1·j - (s2-s1)(2x - 2i·u1) = r(v-X) - 2·s1·j
	gfpSub(T, v, &c.X) // t7
	gfpMul(t4, s1, j)  // t8
	gfpAdd(t6, t4, t4) // t9
	gfpMul(t4, r, T)   // t10
	gfpSub(&c.Y, t4, t6)

	// Set Z = 2(u2-u1)·z1·z2 = 2h·z1·z2
	gfpAdd(T, &a.Z, &b.Z) // t11
	gfpMul(t4, T, T)      // t12
	gfpSub(T, t4, z12)    // t13
	gfpSub(t4, T, z22)    // t14
	gfpMul(&c.Z, t4, h)
}

func (c *curvePoint) Double(a *curvePoint) {
	// See http://hyperelliptic.org/EFD/g1p/auto-code/shortw/jacobian-0/doubling/dbl-2009-l.op3
	A, B, C := &gfP{}, &gfP{}, &gfP{}
	gfpMul(A, &a.X, &a.X)
	gfpMul(B, &a.Y, &a.Y)
	gfpMul(C, B, B)

	T, t2 := &gfP{}, &gfP{}
	gfpAdd(T, &a.X, B)
	gfpMul(t2, T, T)
	gfpSub(T, t2, A)
	gfpSub(t2, T, C)

	d, e, f := &gfP{}, &gfP{}, &gfP{}
	gfpAdd(d, t2, t2)
	gfpAdd(T, A, A)
	gfpAdd(e, T, A)
	gfpMul(f, e, e)

	gfpAdd(T, d, d)
	gfpSub(&c.X, f, T)

	gfpAdd(T, C, C)
	gfpAdd(t2, T, T)
	gfpAdd(T, t2, t2)
	gfpSub(&c.Y, d, &c.X)
	gfpMul(t2, e, &c.Y)
	gfpSub(&c.Y, t2, T)

	gfpMul(T, &a.Y, &a.Z)
	gfpAdd(&c.Z, T, T)
}

func (c *curvePoint) Mul(a *curvePoint, scalar *big.Int) {
	precomp := [1 << 2]*curvePoint{nil, {}, {}, {}}
	precomp[1].Set(a)
	precomp[2].Set(a)
	gfpMul(&precomp[2].X, &precomp[2].X, xiTo2PSquaredMinus2Over3)
	precomp[3].Add(precomp[1], precomp[2])

	multiScalar := curveLattice.Multi(scalar)

	sum := &curvePoint{}
	sum.SetInfinity()
	T := &curvePoint{}

	for i := len(multiScalar) - 1; i >= 0; i-- {
		T.Double(sum)
		if multiScalar[i] == 0 {
			sum.Set(T)
		} else {
			sum.Add(T, precomp[multiScalar[i]])
		}
	}
	c.Set(sum)
}

func (c *curvePoint) MakeAffine() {
	if c.Z == *newGFp(1) {
		return
	} else if c.Z == *newGFp(0) {
		c.X = gfP{0}
		c.Y = *newGFp(1)
		c.T = gfP{0}
		return
	}

	zInv := &gfP{}
	zInv.Invert(&c.Z)

	T, zInv2 := &gfP{}, &gfP{}
	gfpMul(T, &c.Y, zInv)
	gfpMul(zInv2, zInv, zInv)

	gfpMul(&c.X, &c.X, zInv2)
	gfpMul(&c.Y, T, zInv2)

	c.Z = *newGFp(1)
	c.T = *newGFp(1)
}

func (c *curvePoint) Neg(a *curvePoint) {
	c.X.Set(&a.X)
	gfpNeg(&c.Y, &a.Y)
	c.Z.Set(&a.Z)
	c.T = gfP{0}
}
