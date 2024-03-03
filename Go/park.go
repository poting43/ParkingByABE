package main

import (
	"crypto/rand"
	// "encoding/hex"
	"crypto/sha256"
	"fmt"
	"strconv"
	// "encoding/binary"
	"clearmatics/bn256"
	"math/big"
	"reflect"
	// "encoding/binary"
)

// convert the string, which represents the hex value, into a bitInt
func Hex2Dec(input string) *big.Int {
	sum := big.NewInt(0)
	temp, _ := strconv.ParseInt(string(input[0]), 16, 32)
	sum = sum.Add(sum, big.NewInt(temp))
	for i:=1; i < len(input); i++ {
		sum = sum.Mul(sum, big.NewInt(16))
		temp, _ = strconv.ParseInt(string(input[i]), 16, 32)
		sum = sum.Add(sum, big.NewInt(temp))
	}
	return sum
}

// convert bigInt into 256-bit []int
func Big2Bits(input *big.Int) []int{
	temp := new(big.Int)
	temp = temp.Set(input)
	arr := make([]int, 256)
	i := 255
	for  temp.Cmp(big.NewInt(0)) == 1 {
		if( big.NewInt(0).Cmp(new(big.Int).Mod(temp, big.NewInt(2))) == 0 ){
			arr[i] = 0
			temp.Div(temp, big.NewInt(2))
		} else {
			arr[i] = 1
			temp.Div(temp, big.NewInt(2))
		}
		i = i - 1
	}
	return arr
}

// convert 256-bit []int into bigInt
func Bits2Big(input []int) *big.Int{
	sum := big.NewInt(0)
	for i:=0; i < 256; i++ {
		sum.Mul(sum, big.NewInt(2))
		if( input[i] == 1 ){
			sum.Add(sum, big.NewInt(1))
		}
	}
	return sum
}

// xor operation (run through 256 bits)
func Xor256Bits(a []int, b[]int) []int{
	arr := make([]int, 256)
	for i:=0; i < 256; i++ {
		arr[i] = a[i]^b[i];
	}
	return arr
}

// Big H function: G1 -> {0,1}*
func BigH(data *bn256.G1) *big.Int{
	byte32_H1 := sha256.Sum256([]byte(data.String()))
	byte_H1 := byte32_H1[:]
	bigInt_H1 := new(big.Int).SetBytes(byte_H1)
	return bigInt_H1
}

// Little h function:  Z_q^* X G1 X Z_q^* X Z_q^* -> Z_q^*
func LittleH(a *big.Int, b *bn256.G1, c *big.Int, d *big.Int) *big.Int{
	str := string(a.Bytes()) + b.String() + string(c.Bytes()) + string(d.Bytes())
	byte32_H2 := sha256.Sum256([]byte(str))
	byte_H2 := byte32_H2[:]
	bigInt_H2 := new(big.Int).SetBytes(byte_H2)
	return bigInt_H2
}



func main() {

  // SETUP
  r_3, _ := rand.Int(rand.Reader, bn256.Order)
  z := new(bn256.G1).ScalarBaseMult(r_3)
  
  _ = r_3
  _ = z
  
  ID, _ := rand.Int(rand.Reader, bn256.Order)
  _ = ID
  
  // PIDGEN
  x_i, _ := rand.Int(rand.Reader, bn256.Order)
  p_1 := new(bn256.G1).ScalarBaseMult(x_i)
  _ = x_i
  _ = p_1
  p_2 := Xor256Bits(Big2Bits(ID),Big2Bits(BigH(new(bn256.G1).ScalarMult(z,x_i))))
  
  _ = p_2  

	// SIGN
  T_i, _ := rand.Int(rand.Reader, bn256.Order)
  M := new(big.Int).SetBytes([]byte("Parking information"))
  _ = T_i
  _ = M
  
	s_i := new(bn256.G2).ScalarBaseMult(new(big.Int).Add(x_i,new(big.Int).Mul(r_3,LittleH(M,p_1,Bits2Big(p_2),T_i))))
  _ = s_i
  
  // VERIFY
	p1 := new(bn256.G1).ScalarBaseMult(big.NewInt(1))
  p2 := new(bn256.G2).ScalarBaseMult(big.NewInt(1))
  addhash := new(bn256.G1).ScalarMult(z,LittleH(M,p_1,Bits2Big(p_2),T_i))
	_ = p1
	_ = p2
  _ = addhash
  
  e1 := bn256.Pair(p1,s_i)
	e2 := bn256.Pair(new(bn256.G1).Add(p_1,addhash), p2)
	_ = e1
	_ = e2
  
 
	fmt.Println("\nVerificaton via bilinear pairing:", e1.String() == e2.String())
 
  // Trace
  NID := Xor256Bits(p_2,Big2Bits(BigH(new(bn256.G1).ScalarMult(p_1,r_3))))
  _ = NID
  fmt.Println("\nTrace : ", reflect.DeepEqual(Big2Bits(ID),  NID))
  
  fmt.Println("\ndata_ours = [", Hex2Dec(p_1.P.GetX()),",")
  fmt.Print(Hex2Dec(p_1.P.GetY()),",")
	
	fmt.Print(Hex2Dec(addhash.P.GetX()),",")
	fmt.Print(Hex2Dec(addhash.P.GetY()),",")

	fmt.Print(Hex2Dec(s_i.P.GetXX()),",")
	fmt.Print(Hex2Dec(s_i.P.GetXY()),",")
	fmt.Print(Hex2Dec(s_i.P.GetYX()),",")
	fmt.Print(Hex2Dec(s_i.P.GetYY()),"]")









}
