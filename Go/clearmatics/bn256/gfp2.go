package bn256

// For details of the algorithms used, see "Multiplication and Squaring on
// Pairing-Friendly Fields, Devegili et al.
// http://eprint.iacr.org/2006/471.pdf.

// gfP2 implements a field of size p² as a quadratic extension of the base field
// where i²=-1.
type gfP2 struct {
	X, Y gfP // value is xi+Y.
}

func gfP2Decode(in *gfP2) *gfP2 {
	out := &gfP2{}
	montDecode(&out.X, &in.X)
	montDecode(&out.Y, &in.Y)
	return out
}

func (e *gfP2) String() string {
	return "(" + e.X.String() + ", " + e.Y.String() + ")"
}

func (e *gfP2) Set(a *gfP2) *gfP2 {
	e.X.Set(&a.X)
	e.Y.Set(&a.Y)
	return e
}

func (e *gfP2) SetZero() *gfP2 {
	e.X = gfP{0}
	e.Y = gfP{0}
	return e
}

func (e *gfP2) SetOne() *gfP2 {
	e.X = gfP{0}
	e.Y = *newGFp(1)
	return e
}

func (e *gfP2) IsZero() bool {
	zero := gfP{0}
	return e.X == zero && e.Y == zero
}

func (e *gfP2) IsOne() bool {
	zero, one := gfP{0}, *newGFp(1)
	return e.X == zero && e.Y == one
}

func (e *gfP2) Conjugate(a *gfP2) *gfP2 {
	e.Y.Set(&a.Y)
	gfpNeg(&e.X, &a.X)
	return e
}

func (e *gfP2) Neg(a *gfP2) *gfP2 {
	gfpNeg(&e.X, &a.X)
	gfpNeg(&e.Y, &a.Y)
	return e
}

func (e *gfP2) Add(a, b *gfP2) *gfP2 {
	gfpAdd(&e.X, &a.X, &b.X)
	gfpAdd(&e.Y, &a.Y, &b.Y)
	return e
}

func (e *gfP2) Sub(a, b *gfP2) *gfP2 {
	gfpSub(&e.X, &a.X, &b.X)
	gfpSub(&e.Y, &a.Y, &b.Y)
	return e
}

// See "Multiplication and Squaring in Pairing-Friendly Fields",
// http://eprint.iacr.org/2006/471.pdf
func (e *gfP2) Mul(a, b *gfP2) *gfP2 {
	tx, t := &gfP{}, &gfP{}
	gfpMul(tx, &a.X, &b.Y)
	gfpMul(t, &b.X, &a.Y)
	gfpAdd(tx, tx, t)

	ty := &gfP{}
	gfpMul(ty, &a.Y, &b.Y)
	gfpMul(t, &a.X, &b.X)
	gfpSub(ty, ty, t)

	e.X.Set(tx)
	e.Y.Set(ty)
	return e
}

func (e *gfP2) MulScalar(a *gfP2, b *gfP) *gfP2 {
	gfpMul(&e.X, &a.X, b)
	gfpMul(&e.Y, &a.Y, b)
	return e
}

// MulXi sets e=ξa where ξ=i+9 and then returns e.
func (e *gfP2) MulXi(a *gfP2) *gfP2 {
	// (xi+Y)(i+9) = (9x+Y)i+(9y-X)
	tx := &gfP{}
	gfpAdd(tx, &a.X, &a.X)
	gfpAdd(tx, tx, tx)
	gfpAdd(tx, tx, tx)
	gfpAdd(tx, tx, &a.X)

	gfpAdd(tx, tx, &a.Y)

	ty := &gfP{}
	gfpAdd(ty, &a.Y, &a.Y)
	gfpAdd(ty, ty, ty)
	gfpAdd(ty, ty, ty)
	gfpAdd(ty, ty, &a.Y)

	gfpSub(ty, ty, &a.X)

	e.X.Set(tx)
	e.Y.Set(ty)
	return e
}

func (e *gfP2) Square(a *gfP2) *gfP2 {
	// Complex squaring algorithm:
	// (xi+Y)² = (X+Y)(Y-X) + 2*i*X*Y
	tx, ty := &gfP{}, &gfP{}
	gfpSub(tx, &a.Y, &a.X)
	gfpAdd(ty, &a.X, &a.Y)
	gfpMul(ty, tx, ty)

	gfpMul(tx, &a.X, &a.Y)
	gfpAdd(tx, tx, tx)

	e.X.Set(tx)
	e.Y.Set(ty)
	return e
}

func (e *gfP2) Invert(a *gfP2) *gfP2 {
	// See "Implementing cryptographic pairings", M. Scott, section 3.2.
	// ftp://136.206.11.249/pub/crypto/pairings.pdf
	t1, t2 := &gfP{}, &gfP{}
	gfpMul(t1, &a.X, &a.X)
	gfpMul(t2, &a.Y, &a.Y)
	gfpAdd(t1, t1, t2)

	inv := &gfP{}
	inv.Invert(t1)

	gfpNeg(t1, &a.X)

	gfpMul(&e.X, t1, inv)
	gfpMul(&e.Y, &a.Y, inv)
	return e
}
