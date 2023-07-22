package perlin

import (
	"math"
	"math/rand"
)

type Vec2 struct {
	X, Y float64
}

func DotProduct(lhs Vec2, rhs Vec2) float64 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y
}

func Fade(t float64) float64 {
	return ((6*t-15)*t + 10) * math.Pow(t, 3)
}

func Lerp(t float64, a1 float64, a2 float64) float64 {
	return a1 + t*(a2-a1)
}

func GetConstantVec2(v int) Vec2 {

	switch h := v & 3; h {
	case 0:
		return Vec2{1.0, 1.0}

	case 1:
		return Vec2{-1.0, 1.0}

	case 2:
		return Vec2{-1.0, -1.0}

	default:
		return Vec2{1.0, -1.0}
	}
}

func ShuffleArray(array_to_shuffle [256]int) [256]int {
	for e := len(array_to_shuffle) - 1; e > 0; e-- {
		var index int = int(math.Round(rand.Float64() * float64(e-1)))
		array_to_shuffle[e], array_to_shuffle[index] =
			array_to_shuffle[index], array_to_shuffle[e]
	}

	return array_to_shuffle
}

func GeneratePermutation() [512]int {
	var permutation1 [256]int

	for i := 0; i < 256; i++ {
		permutation1[i] = i
	}

	var permutation [512]int

	permutation1 = ShuffleArray(permutation1)

	for i := 0; i < 256; i++ {
		permutation[i] = permutation1[i]
	}

	for i := 0; i < 256; i++ {
		permutation[256+i] = permutation[i]
	}

	return permutation
}

func Generate2DNoise(x float64, y float64, permutation [512]int) float64 {
	var X int = int(math.Floor(x)) & 255
	var Y int = int(math.Floor(y)) & 255

	var xf float64 = x - math.Floor(x)
	var yf float64 = y - math.Floor(y)

	var top_right Vec2 = Vec2{xf - 1.0, yf - 1.0}
	var top_left Vec2 = Vec2{xf, yf - 1.0}
	var bottom_right Vec2 = Vec2{xf - 1.0, yf}
	var bottom_left Vec2 = Vec2{xf, yf}

	value_top_right := permutation[permutation[X+1]+Y+1]
	value_top_left := permutation[permutation[X]+Y+1]
	value_bottom_right := permutation[permutation[X+1]+Y]
	value_bottom_left := permutation[permutation[X]+Y]

	dot_top_right := DotProduct(top_right, GetConstantVec2(value_top_right))
	dot_top_left := DotProduct(top_left, GetConstantVec2(value_top_left))
	dot_bottom_right := DotProduct(bottom_right, GetConstantVec2(value_bottom_right))
	dot_bottom_left := DotProduct(bottom_left, GetConstantVec2(value_bottom_left))

	var u float64 = Fade(xf)
	var v float64 = Fade(yf)

	return Lerp(u,
		Lerp(v, dot_bottom_left, dot_top_left),
		Lerp(v, dot_bottom_right, dot_top_right))
}

func FractalBrownianMotion(x float64, y float64, num_octaves int, permutation [512]int) float64 {
	result := 0.0
	amplitude := 1.0
	frequency := 0.005

	for octave := 0; octave < int(num_octaves); octave++ {
		n := amplitude * Generate2DNoise(x*frequency, y*frequency, permutation)
		result += n

		amplitude *= 0.5
		frequency *= 2.0
	}

	return result
}
