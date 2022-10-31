package ecc

import (
	"bitcoin-go/utility"
	"math/big"
)

type Point struct {
	x *big.Int
	y *big.Int
	a *big.Int
	b *big.Int
}

func NewSecp256k1Point(x *big.Int, y *big.Int) Point {
	return newSecp256k1PointWithCustomAAndB(x, y, A, B)
}

// User beware. This will create a point with the sufficient values, but many of the point operations still assume prime = P.
func newSecp256k1PointWithCustomAAndB(x *big.Int, y *big.Int, a *big.Int, b *big.Int) Point {

	if x != nil && y != nil {

		y_squared := ModPowInt(y, 2)
		sum := ModAdd(ModAdd(ModPowInt(x, 3), ModMul(a, x)), b)

		if y_squared.Cmp(sum) != 0 {
			panic("Bad point parameters.")
		}
	}

	return Point{x: x, y: y, a: a, b: b}
}

func (p *Point) Equals(p2 *Point) bool {
	if p2 == nil {
		return false
	}

	if p.x == nil && p2.x != nil {
		return false
	}

	if p.y == nil && p2.x != nil {
		return false
	}

	if p.x != nil && p.y != nil {
		if p.x.Cmp(p2.x) != 0 || p.y.Cmp(p2.y) != 0 {
			return false
		}
	}

	if p.a.Cmp(p2.a) == 0 && p.b.Cmp(p2.b) == 0 {
		return true
	}

	return false
}

func (p *Point) NotEquals(p2 *Point) bool {
	return !p.Equals(p2)
}

func (p *Point) Add(p2 *Point) Point {

	if p.a.Cmp(p2.a) != 0 || p.b.Cmp(p2.b) != 0 {
		panic("Points aren't on the same curve.")
	}

	// Case 0.0: self is the point at infinity, return other
	if p.x == nil {
		return *p2
	}

	// Case 0.1: other is the point at infinity, return self
	if p2.x == nil {
		return *p
	}

	// Case 1: self.x == other.x, self.y != other.y
	// Result is point at infinity
	if p.x.Cmp(p2.x) == 0 && p.y.Cmp(p2.y) != 0 {
		return Point{nil, nil, p.a, p.b}
	}

	// Case 2: self.x ≠ other.x
	// Formula (x3,y3)==(x1,y1)+(x2,y2)
	// s=(y2-y1)/(x2-x1)
	// x3=s**2-x1-x2
	// y3=s*(x1-x3)-y1
	if p.x.Cmp(p2.x) != 0 {
		s := ModDiv(ModSub(p2.y, p.y), ModSub(p2.x, p.x))
		x := ModSub(ModSub(ModPowInt(s, 2), p.x), p2.x)
		y := ModSub(ModMul(s, ModSub(p.x, x)), p.y)

		return Point{x, y, p.a, p.b}
	}

	// Case 4: if we are tangent to the vertical line,
	// we return the point at infinity
	// note instead of figuring out what 0 is for each type
	// we just use 0 * self.x
	mult_x_by_zero := ModMulInt(p.x, 0)
	if p.Equals(p2) && p.y.Cmp(mult_x_by_zero) == 0 {
		return Point{nil, nil, p.a, p.b}
	}

	// Case 3: self == other
	// Formula (x3,y3)=(x1,y1)+(x1,y1)
	// s=(3*x1**2+a)/(2*y1)
	// x3=s**2-2*x1
	// y3=s*(x1-x3)-y1
	if p.Equals(p2) {
		left_side := ModAdd(ModMulInt(ModPowInt(p.x, 2), 3), p.a)
		right_side := ModMulInt(p.y, 2)
		s := ModDiv(left_side, right_side)
		x := ModSub(ModPowInt(s, 2), ModMulInt(p.x, 2))
		y := ModSub(ModMul(s, ModSub(p.x, x)), p.y)
		return Point{x, y, p.a, p.b}
	}

	panic("We shouldn't get here!")
}

func (p *Point) ScalarMultiply(coefficient *big.Int) Point {

	coef := big.NewInt(0)
	coef = coef.Mod(coefficient, N)
	current := p.Clone()
	result := Point{nil, nil, p.a, p.b}

	for coef.Cmp(BigZero) != 0 {

		and := new(big.Int)
		and.And(coef, BigOne)

		if and.Cmp(BigZero) != 0 {
			result = result.Add(&current)
		}

		current = current.Add(&current)
		coef = coef.Rsh(coef, 1)
	}

	return result
}

func (p *Point) Verify(hash *big.Int, sig Signature) bool {
	s_inv := ModPowPrime(sig.S, ModSubInt(N, 2), N)
	u := ModMulPrime(hash, s_inv, N)
	v := ModMulPrime(sig.R, s_inv, N)

	sub_total_uG := G.ScalarMultiply(u)
	sub_total_vSelf := p.ScalarMultiply(v)
	total := sub_total_uG.Add(&sub_total_vSelf)
	return total.x.Cmp(sig.R) == 0
}

func (p *Point) ToSEC(compressed bool) []byte {
	if compressed {
		bytes := make([]byte, 33)
		isEven := IsEven(p.y)

		bytes[0] = utility.IIF(isEven, (byte)(0x02), (byte)(0x03)).(byte)

		fillBufferWithIntBytes(bytes[1:], p.x, false)
		return bytes
	} else {
		bytes := make([]byte, 65)
		bytes[0] = 0x04
		fillBufferWithIntBytes(bytes[1:], p.x, false)
		fillBufferWithIntBytes(bytes[33:], p.y, false)
		return bytes
	}
}

func NewPointFromSEC(buffer []byte) Point {

	// TODO: Modify to return (Point, error) or (Point, ok) for failure cases.
	if buffer[0] == 0x04 {
		x := new(big.Int)
		y := new(big.Int)
		x.SetBytes(buffer[1:33])
		y.SetBytes(buffer[33:65])
		return NewSecp256k1Point(x, y)
	} else {
		x := new(big.Int)
		x.SetBytes(buffer[1:33])
		yEven := (buffer[0] == 0x02)

		alpha := ModAdd(ModPowInt(x, 3), B)
		beta := ModSqrt(alpha)
		alt_beta := ModSub(P, beta)
		even_beta := beta
		odd_beta := alt_beta

		if !IsEven(beta) {
			even_beta = alt_beta
			odd_beta = beta
		}

		return NewSecp256k1Point(x, utility.IIF(yEven, even_beta, odd_beta).(*big.Int))
	}
}

func (p *Point) Hash160(compressed bool) []byte {
	return utility.Hash160(p.ToSEC(compressed))
}

func (p *Point) Address(compressed bool, testnet bool) string {
	h160 := p.Hash160(compressed)
	concat := make([]byte, len(h160)+1)
	concat[0] = utility.IIF(testnet, (byte)(0x6f), (byte)(0x00)).(byte)
	copy(concat[1:], h160)
	return utility.EncodeBase58Checksum(concat)
}

func (p *Point) Clone() Point {
	clone := Point{x: new(big.Int), y: new(big.Int), a: new(big.Int), b: new(big.Int)}
	clone.x.Set(p.x)
	clone.y.Set(p.y)
	clone.a.Set(p.a)
	clone.b.Set(p.b)
	return clone
}

func fillBufferWithIntBytes(b []byte, i *big.Int, littleEndian bool) {
	intBytes := i.Bytes()

	if littleEndian {
		for i, c := 0, len(intBytes)-1; c >= 0; i, c = i+1, c-1 {
			b[i] = intBytes[c]
		}
	} else {
		// In the case there aren't a full 32 bytes, skip the pad bytes.
		pad := 32 - len(intBytes)
		for c := 0; c < len(intBytes); c++ {
			b[c+pad] = intBytes[c]
		}
	}
}
