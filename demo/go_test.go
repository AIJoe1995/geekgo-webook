package demo

import (
	"fmt"
	"math"
	"testing"
)

func TestE(t *testing.T) {
	fmt.Println(findTheArrayConcVal([]int{7,52,2,4})
	)
}

func findTheArrayConcVal(nums []int) int64 {
	i, j := 0, 0
	ans := 0
	for i < j {
		ans += nums[i]*int(math.Pow(10, float64(getNumDigits(nums[j])))) + nums[j]
		i += 1
		j -= 1
	}
	if i == j {
		ans += nums[i]
	}
	return int64(ans)

}

func getNumDigits(num int) int {
	cnt := 0
	for num != 0 {
		cnt += 1
		num = num / 10
	}
	return cnt
}
