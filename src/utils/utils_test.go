package utils

import (
	"fmt"
	"testing"
)

func TestRandomFloat32(t *testing.T) {
	for i:=0; i<1000; i++ {
		floatArray := RandomFloat32(10)
		for _, v := range floatArray {
			if v < -1 || v >= 1 {
				t.Errorf("RandomFloat32返回的值为(%f), 不在[-1, 1)之间", v)
			}
		}
		if len(floatArray) != 10 {
			t.Errorf("RandomFloat32返回的数组长度:(%d)与设置:(%d)不一致", len(floatArray),
				10)
		}
		valSet := make(map[float32]interface{})
		for _, v := range floatArray {
			valSet[v] = ""
		}
		if len(valSet) != len(floatArray) {
			t.Errorf("RandomFloat32返回的数组的值有重复")
		}
	}
}

func TestMeanArray(t *testing.T) {
	array1 := []float32 {1, 2, 3}
	array2 := []float32 {1, 2, 3}
	array3 := [][]float32 {array1, array2}
	meanArray := MeanArray(array3, 2, 3)
	if !(meanArray[0] == 1 && meanArray[1] == 2 && meanArray[2] == 3) {
		t.Errorf("MeanArray计算错误, 测试结果:[%f, %f, %f], 正确结果：[1, 2, 3]",
			meanArray[0], meanArray[1], meanArray[2])
	}
}

func TestArrayAdd(t *testing.T) {
	array1 := []float32 {1, 2, 3}
	array2 := []float32 {1, 2, 3}
	array3 := ArrayAdd(array1, array2)
	if !(array3[0] == 2 && array3[1] == 4 && array3[2] == 6) {
		t.Errorf("ArrayAdd计算错误，测试结果：[%f, %f, %f]，正确结果：[2, 4, 6]",
			array3[0], array3[1], array3[2])
	}
}

func TestArrayMultiply(t *testing.T) {
	array1 := []float32 {1, 2, 3}
	array2 := []float32 {1, 2, 3}
	res := ArrayMultiply(array1, array2)
	if res != 14 {
		t.Errorf("ArrayMulity计算结果错误，测试结果:%f，正确结果：14", res)
	}
}

func TestArrayMultiply1(t *testing.T) {
	c := 2
	array := []float32 {1, 2, 3}
	res := ArrayMultiply1(float32(c), array)
	if !(res[0] == 2 && res[1] == 4 && res[2] == 6) {
		t.Errorf("ArrayMultipy计算结果错误，测试结果：[%f, %f, %f]，正确结果：[2, 4, 5]",
			res[0], res[1], res[2])
	}
}

func abs(f float32) float32 {
	if f < 0 {
		f = -1 * f
	}
	return f
}

func TestSigmod(t *testing.T) {
	array1 := []float32 {2, 0, 1}
	array2 := []float32 {1, 2, 3}
	res := Sigmod(array1, array2)
	if abs(res - 0.993307) > 0.01 {
		t.Errorf("Sigmod计算结果错误，测试结果：%f，正确结果为：0.993307",
			res)
	}
}

func TestArraySqrt(t *testing.T) {
	array := []float32 {1, 1, 1}
	res := ArraySqrt(array)
	if abs(res - 1.73205) > 0.01 {
		t.Errorf("ArraySqrt计算结果错误，测试结果：%f，正确结果为：1.73205",
			res)
	}
}

func TestCosDistance(t *testing.T) {
	array1 := []float32 {1, 1, 1}
	array2 := []float32 {2, 2, 2}
	array3 := []float32 {-1, -1, -1}
	res1 := CosDistance(array1, array2)
	res2 := CosDistance(array2, array3)
	if abs(res1 - 1) > 0.01 {
		t.Errorf("CosDistance计算结果错误，测试结果为:%f，正确结果为：1",
			res1)
	}
	if abs(res2 + 1) > 0.01 {
		t.Errorf("CosDistance计算结果错误，测试结果为：%f，正确结果为：-1",
			res2)
	}
}

func ExampleVector2String() {
	array := []float32 {1, 2, 3}
	str := Vector2String(array)
	fmt.Println(*str)
	// Output:
	//  1.0000000000 2.0000000000 3.0000000000
}

func ExampleString2Vector() {
	str := "w 1 2 3"
	word, vector := String2Vector(str)
	fmt.Println(word)
	for _, v := range vector {
		fmt.Println(v)
	}
	// Output:
	// w
	// 1
	// 2
	// 3
}
