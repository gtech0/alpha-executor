package service

import (
	"alpha-executor/entity"
	"fmt"
	"regexp"
)

var patterns = map[string]string{
	"RANGE":     "RANGE\\s+\\S+\\s+\\S+\\s*",
	"GET":       "GET\\s+[^\\s,:;](\\s+\\(\\d+\\)){0,1}\\s+\\(\\s*[^\\s,]+\\s*(\\,\\s*[^\\s,]+\\s*)*\\)\\s*:\\s*.+",
	"HOLD":      "HOLD\\s+[^\\s,:;]\\s+\\(\\s*[^\\s,]+\\s*(\\,\\s*[^\\s,]+\\s*)*\\)\\s*:\\s*.+",
	"UPDATE":    "UPDATE\\s+\\S+\\s*",
	"DELETE":    "DELETE\\s+\\S+\\s*",
	"RELEASE":   "RELEASE\\s+\\S+\\s*",
	"PUT":       "PUT\\s+\\S+\\s+\\(\\S+\\)\\s*",
	"ASSIGMENT": "\\S+\\s*=\\s*\\S+\\s*",
}

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
				Message:   fmt.Sprintf("No matching in line %d", line),
			}
		}
	}

	return nil
}
