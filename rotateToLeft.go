func rotateArray(array []string, size, rotation int) string {
	var newArray []string
	for i := 0; i < rotation; i++ {
		newArray = array[1:size]
		newArray = append(newArray, array[0])
		array = newArray
	}
	return strings.Join(array, " ")
}
