package utils

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func RandomFloat32(n int) []float32 {
	res := make([]float32, n, n)
	for i:=0; i<n; i++ {
		res[i] = 2 * rand.Float32() - 1
	}
	return res
}

func MeanArray(arrayLst [][]float32, row int, col int) ([]float32){
	res := make([]float32, col, col)
	for i:=0; i<col; i++ {
		res[i] = 0
		for j:=0; j<row; j++ {
			res[i] += arrayLst[j][i]
		}
		res[i] /= float32(row)
	}
	return res
}

func ArrayAdd(array1 []float32, array2 []float32) []float32 {
	n := len(array1)
	res := make([]float32, n, n)
	for i:=0; i<n; i++ {
		res[i] = array1[i] + array2[i]
	}
	return res
}

func ArrayMultiply(array1 []float32, array2 []float32) float32 {
	n := len(array1)
	res := float32(0)
	for i:=0; i<n; i++ {
		res += array1[i] * array2[i]
	}
	return res
}

func ArrayMultiply1(c float32, array []float32) []float32 {
	res := make([]float32, len(array), len(array))
	for i:=0; i<len(array); i++ {
		res[i] = c * array[i]
	}
	return res
}

func Sigmod(vector []float32, theta []float32) float32 {
	r := ArrayMultiply(vector, theta)
	return 1.0 / (1 + (float32)(math.Exp(float64(-1.0*r))))
}

func ArraySqrt(array []float32) float32 {
	res := float32(0)
	for _, v := range array {
		res += v * v
	}
	return float32(math.Sqrt(float64(res)))
}

func CosDistance(array1 []float32, array2 []float32) float32 {
	son := ArrayMultiply(array1, array2)
	res := son / (ArraySqrt(array1) * ArraySqrt(array2))
	return res
}

func Vector2String(vector []float32) *string {
	str := " "
	for _, v := range vector {
		str += strconv.FormatFloat(float64(v), 'f', 10, 32)
		str += " "
	}
	return &str
}

func String2Vector(line string) (string, []float32) {
	line = strings.Trim(line, " ")
	items := strings.Split(line, " ")
	word := items[0]
	vector := make([]float32, len(items)-1, len(items)-1)
	for i:=0; i<len(items)-1; i++ {
		val, err := strconv.ParseFloat(items[i+1], 32)
		if err != nil {
			fmt.Println(items[i+1], "转换为float32失败")
			os.Exit(-1)
		}
		vector[i] = float32(val)
	}
	return word, vector
}