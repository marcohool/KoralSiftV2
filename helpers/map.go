package helpers

func CreateSliceFromMap(set map[string]struct{}) []string {
	slice := make([]string, 0, len(set))

	for key := range set {
		slice = append(slice, key)
	}

	return slice
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
