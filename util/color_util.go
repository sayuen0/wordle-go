package util

type ColorScheme string

const (
	BgYellow      ColorScheme = "\033[1;43m%s\033[0m"
	BgGreen       ColorScheme = "\033[1;42m%s\033[0m"
	BgGray        ColorScheme = "\033[1;47m%s\033[0m"
	BgBrightRed   ColorScheme = "\033[1;101m%s\033[0m"
	BgBrightGreen ColorScheme = "\033[1;102m%s\033[0m"
	BgBrightBlue  ColorScheme = "\033[1;104m%s\033[0m"
)
