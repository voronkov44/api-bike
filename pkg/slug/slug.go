package slug

import (
	"regexp"
	"strconv"
	"strings"
)

var nonAlNum = regexp.MustCompile(`[^a-z0-9\-]+`)
var multiDash = regexp.MustCompile(`\-{2,}`)

var ru = []rune("абвгдеёжзийклмнопрстуфхцчшщъыьэюя")
var en = []string{
	"a", "b", "v", "g", "d", "e", "e", "zh", "z", "i", "y", "k", "l", "m", "n", "o", "p", "r", "s", "t", "u", "f", "h",
	"ts", "ch", "sh", "sch", "", "y", "", "e", "yu", "ya",
}
var ruMap map[rune]string

func init() {
	ruMap = make(map[rune]string, len(ru))
	for i, r := range ru {
		ruMap[r] = en[i]
	}
}

// Slugify: "Пицца Маргарита" -> "pizza-margarita"

func Slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	var b strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == ' ' || r == '-' {
			b.WriteRune(r)
			continue
		}
		// ru translit
		if tr, ok := ruMap[r]; ok {
			b.WriteString(tr)
			continue
		}
		// space-ish for others
		b.WriteRune(' ')
	}
	// spaces -> dash
	out := strings.ReplaceAll(b.String(), " ", "-")
	out = nonAlNum.ReplaceAllString(out, "")
	out = multiDash.ReplaceAllString(out, "-")
	out = strings.Trim(out, "-")
	if out == "" {
		out = "item"
	}
	return out
}

// WithSuffix("pizza-margarita", 2) -> "pizza-margarita-2"

func WithSuffix(base string, n int) string {
	return base + "-" + strconv.Itoa(n)
}
