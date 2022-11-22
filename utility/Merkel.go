package utility

func MerkelParent(hash1 []byte, hash2 []byte) []byte {
	new := append(hash1, hash2...)
	return Hash256(new)
}

func MerkelParentLevel(hashes [][]byte) [][]byte {

	numHashes := len(hashes)
	numParents := numHashes / 2

	// Account for odd numbers.
	if numHashes%2 == 1 {
		numParents++
	}

	results := make([][]byte, numParents)
	parentIndex := 0

	for i := 0; i < numHashes; i += 2 {

		// Get the left and right hashes
		left := hashes[i]
		right := hashes[i] // Default to odd number of hashes case (reuse the last one)

		// If there's room, grab the next one.
		if i+1 < numHashes {
			right = hashes[i+1]
		}

		// "Merkelize" these two and save the result.
		merkeled := MerkelParent(left, right)
		results[parentIndex] = merkeled
		parentIndex++
	}

	return results
}

func MerkelRoot(hashes [][]byte) []byte {

	for len(hashes) > 1 {
		hashes = MerkelParentLevel(hashes)
	}

	// Return the root and make it little-endian
	return hashes[0]
}
