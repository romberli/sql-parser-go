package token

import (
	"github.com/romberli/go-util/constant"
)

func TypeExists(tokenTypeList []Type, t Type) bool {
	for _, tokenType := range tokenTypeList {
		if tokenType == t {
			return true
		}
	}

	return false
}

func HasIntersect(tokenTypeList [][]Type) bool {
	if len(tokenTypeList) <= 1 {
		return false
	}

	temp := tokenTypeList[constant.ZeroInt]
	for i := 1; i < len(tokenTypeList); i++ {
		if hasIntersect(temp, tokenTypeList[i]) {
			return true
		}
	}

	return false
}

func hasIntersect(first, second []Type) bool {
	for _, t := range first {
		if TypeExists(second, t) {
			return true
		}
	}

	return false
}
