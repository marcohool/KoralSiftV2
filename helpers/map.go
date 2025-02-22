package helpers

func CreateSliceFromMapKey(set map[string]struct{}) []string {
	slice := make([]string, 0, len(set))

	for key := range set {
		slice = append(slice, key)
	}

	return slice
}
