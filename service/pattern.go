package service

import (
	"alpha-executor/entity"
	"fmt"
	"regexp"
)

var patterns = map[string]string{}

func MatchQuery(query []string) error {
	for index := 0; index < len(query); index++ {
		var line = index + 1
		var matches = false
		for _, pattern := range patterns {
			result, _ := regexp.MatchString(pattern, query[index])
			if result {
				matches = true
				break
			}
		}

		if matches == false {
			return &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   fmt.Sprint("No matching in line ", line),
			}
		}
	}

	return nil
}
