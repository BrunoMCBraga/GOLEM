package util

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
