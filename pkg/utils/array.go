package utils

// The function checks if a given string exists in a slice of strings and returns a boolean value.
func ValueExist(find string, in []string) bool {
	for _, v := range in {
		if v == find {
			return true
		}
	}
	return false
}

// The ReverseArray function takes a list of strings and returns a new list with the elements in
// reverse order.
func ReverseArray(list []string) []string {
	length := len(list)
	reversed := make([]string, length)
	for i, j := 0, length-1; i < length; i, j = i+1, j-1 {
		reversed[i] = list[j]
	}
	return reversed
}

func RemoveDuplicates(strList []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueList := []string{}
	for _, str := range strList {
		if !uniqueMap[str] {
			uniqueMap[str] = true
			uniqueList = append(uniqueList, str)
		}
	}

	return uniqueList
}
