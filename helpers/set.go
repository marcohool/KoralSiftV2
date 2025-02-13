package helpers

func CreateSet(values []string) map[string]bool {
	set := make(map[string]bool)

	for _, value := range values {
		set[value] = true
	}

	return set
}
