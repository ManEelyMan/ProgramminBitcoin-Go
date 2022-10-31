package main

import (
	"bitcoin-go/utility"
	"fmt"
)

func main() {

	// arr := [...]fieldElement{
	// 	NewFieldElement(1, 11),
	// 	NewFieldElement(2, 11),
	// 	NewFieldElement(3, 11),
	// 	NewFieldElement(4, 11),
	// 	NewFieldElement(5, 11),
	// 	NewFieldElement(6, 11),
	// 	NewFieldElement(7, 11),
	// 	NewFieldElement(8, 11),
	// 	NewFieldElement(9, 11),
	// 	NewFieldElement(10, 11),
	// }

	// for i := 0; i < len(arr); i++ {
	// 	(arr[i].Power(9)).Print()
	// }

	// var a ecc.FieldElement = ecc.NewFieldElement(7, 19)
	// var b ecc.FieldElement = ecc.NewFieldElement(5, 19)
	//var c fieldElement = NewFieldElement(3, 31)

	// res := a.Multiply(&b)
	// res.Print()

	// res = res.Multiply(&c)

	// res := a.Divide(&b)
	// res.Print()

	//var pt = ecc.NewPoint(2, 4, 5, 7)
	// var pt2 = ecc.NewPoint2(-1, -1, 5, 7)
	// var pt3 = ecc.NewPoint2(18, 77, 5, 7)
	//var pt4 = ecc.NewPoint(5, 7, 5, 7)

	// pt2.Equals(&pt3)

	//mod := big.NewInt(10)

	// fmt.Println("Mod 10:")
	// for i := -10; i < 10; i++ {

	// 	//tmp := big.NewInt((int64)(i))
	// 	//tmp3 := tmp.Mod(tmp, mod)

	// 	//fmt.Printf("\t Big Int\t%+v\n", tmp3)
	// 	fmt.Printf("\t Reg Int\t%+v\n", i%10)
	// }

	// blah := ecc.NewSecp256k1Point(big.NewInt(3), big.NewInt(-7), big.NewInt(5), big.NewInt(7))
	// blah = ecc.NewSecp256k1Point(big.NewInt(18), big.NewInt(77), big.NewInt(5), big.NewInt(7))

	// fmt.Printf("Blah: %+v\n", blah)

	i := utility.HexStringToBigInt("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d")
	bytes := i.Bytes()
	s := utility.EncodeBase58(bytes)

	fmt.Printf("%v", s)

	fmt.Println("Hello, 世界")
}
