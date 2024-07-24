package zmath

import "time"

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func Sum[T Numeric](nums ...T) T {
	var sum T
	for _, v := range nums {
		sum += v
	}
	return sum
}

func Avg[T Numeric](nums ...T) float64 {
	return float64(Sum(nums...)) / float64(len(nums))
}

func Min[T Numeric](nums ...T) T {
	if len(nums) == 0 {
		return 0
	}
	min := nums[0]
	for _, v := range nums {
		if v < min {
			min = v
		}
	}
	return min
}

func Max[T Numeric](nums ...T) T {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for _, v := range nums {
		if v > max {
			max = v
		}
	}
	return max
}

func SumDuration(durations ...time.Duration) time.Duration {
	return time.Duration(Sum(durations...))
}

func AvgDuration(durations ...time.Duration) time.Duration {
	return time.Duration(Avg(durations...))
}

func MinDuration(durations ...time.Duration) time.Duration {
	return time.Duration(Min(durations...))
}

func MaxDuration(durations ...time.Duration) time.Duration {
	return time.Duration(Max(durations...))
}
