package entity

type Query []string

func (q *Query) DeleteEmpty() {
	var newQuery []string
	for _, value := range *q {
		if value != "" {
			newQuery = append(newQuery, value)
		}
	}
	*q = newQuery
}
