package util

func GetFileTypeByFilename(filename string) string {
	dot := false
	fileType := ""
	for _, symbol := range filename {

		if dot || symbol == '.' {
			dot = true
			fileType += string(symbol)
		}
	}

	return fileType
}
