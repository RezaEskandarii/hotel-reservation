package utils

import (
	"strconv"
	"strings"
)

// ConvertToUint gives interface as a input and converts interface to uint.
func ConvertToUint(input interface{}) (uint64, error) {

	if input == nil {
		return 0, nil
	}

	switch input.(type) {
	case string:
		returnValue, err := strconv.ParseUint(input.(string), 10, 64)
		if err != nil {
			return 0, err
		}

		return returnValue, nil
	}

	return 0, nil
}

// GenerateCacheKey returns string to use as a cache key.
func GenerateCacheKey(keys ...string) string {
	strBuilder := strings.Builder{}
	for _, str := range keys {
		strBuilder.Write([]byte(str))
	}
	return GenerateSHA256(strBuilder.String())
}
