package util

func SplitSliceIntoSubslices(slice []string, numberOfSlices int) []interface{} {

	subSlices := make([]interface{}, 0)
	sliceLength := len(slice) / numberOfSlices //3/2 = 1
	sizeOfLastSlice := len(slice) % numberOfSlices
	lastSubsliceIndex := 0

	if numberOfSlices > len(slice) {
		for _, sliceElement := range slice {
			subSlice := make([]string, 0)
			subSlice = append(subSlice, sliceElement)
			subSlices = append(subSlices, subSlice)
		}
	} else {

		for lastSubsliceIndex = 0; lastSubsliceIndex < numberOfSlices; lastSubsliceIndex++ {
			subSlice := make([]string, 0)
			for i := 0; i < sliceLength; i++ {
				subSlice = append(subSlice, slice[lastSubsliceIndex*sliceLength+i])
			}
			subSlices = append(subSlices, subSlice)
		}

		remainingSlice := slice[lastSubsliceIndex:]

		if sizeOfLastSlice > 0 {
			for remainingSliceIndex, remainingSliceElement := range remainingSlice {
				subSlices[remainingSliceIndex] = append(subSlices[remainingSliceIndex].([]string), remainingSliceElement)
			}
		}
	}

	return subSlices
}

/*
//Old and wrong...
func SplitSliceIntoSubslices(slice []string, subSlicesMaximumSize int) []interface{} {

	subSlices := make([]interface{}, 0)
	numberOfSubSlices := len(slice) / subSlicesMaximumSize //3/2 = 1
	sizeOfLastSlice := len(slice) % subSlicesMaximumSize
	lastSubsliceIndex := 0

	for lastSubsliceIndex = 0; lastSubsliceIndex < numberOfSubSlices; lastSubsliceIndex++ {
		subSlice := make([]string, 0)
		for i2 := 0; i2 < subSlicesMaximumSize; i2++ {
			subSlice = append(subSlice, slice[lastSubsliceIndex*subSlicesMaximumSize+i2])
		}
		subSlices = append(subSlices, subSlice)
	}

	if sizeOfLastSlice != 0 {
		subSlice := make([]string, 0)
		for i1 := 0; i1 < sizeOfLastSlice; i1++ {
			subSlice = append(subSlice, slice[lastSubsliceIndex*subSlicesMaximumSize+i1])
		}
		subSlices = append(subSlices, subSlice)
	}

	return subSlices
}
*/
