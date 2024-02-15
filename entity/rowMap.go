package entity

type RowMap map[string][]string

func (r *RowMap) RowsEqual(r2 *RowMap) bool {
	if !r.keysEqual(r2) {
		return false
	}

	for key1, values1 := range *r {
		values2, exists := (*r2)[key1]
		if !exists || !OrderedSlicesEqual(values1, values2) {
			return false
		}
	}
	return true
}

func (r *RowMap) keysEqual(r2 *RowMap) bool {
	if len(*r) != len(*r2) {
		return false
	}

	exists := make(map[string]struct{})
	for key := range *r {
		exists[key] = struct{}{}
	}

	for key := range *r2 {
		if _, ok := exists[key]; !ok {
			return false
		}
	}
	return true
}
