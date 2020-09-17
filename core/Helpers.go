package core

func checkIfItemWithIdExistsInArray(items []string, itemToFind string) (bool, int) {
	for i, item := range items {
		if item == itemToFind {
			return true, i
		}
	}

	return false, -1
}
