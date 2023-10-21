// utils.go

package utils

import (
	"spysat/logger"

	"github.com/gin-gonic/gin"
)

func Error(err error, c *gin.Context, statusCode int) {
	logger.Error("", err.Error())
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Keys(m map[string]string) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func RemoveDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
