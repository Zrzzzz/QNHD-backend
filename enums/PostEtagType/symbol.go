package PostEtagType

var msgSymbol = map[Enum]string{
	NONE:      "",
	RECOMMEND: "recommend",
	THEME:     "theme",
	TOP:       "top",
}

func (code Enum) GetSymbol() string {
	return msgSymbol[code]
}
