package PostEtagType

var msgSymbol = map[Enum]string{
	NONE:      "",
	RECOMMEND: "recommend",
	THEME:     "theme",
	TOP:       "top",
}

func GetSymbol(code Enum) string {
	return msgSymbol[code]
}
