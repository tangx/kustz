package kustz

import (
	"math"
	"strconv"
)

func (kz *Config) CommonLabels() map[string]string {
	return map[string]string{
		"app": kz.Name,
	}
}

func StringToInt32(val string) int32 {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0
	}
	return int32(i)
}
