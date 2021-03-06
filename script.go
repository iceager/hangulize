package hangulize

import (
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// script represents a writing system.
type script interface {
	Is(rune) bool
	Normalize(rune) rune
	TransliteratePunct(rune) string
}

// scripts is the registry of Scripts by their name.
var scripts = map[string]script{
	// Latin is the default.
	"": &_Latin{},

	"cyrillic": &_Cyrillic{},
	"georgian": &_Georgian{},
	"greek":    &_Greek{},
	"kana":     &_Kana{},
	"latin":    &_Latin{},
	"pinyin":   &_Pinyin{},
}

// getScript chooses a script by the script name.
func getScript(name string) script {
	script, ok := scripts[name]
	if !ok {
		// Get the default.
		latin := scripts[""]
		return latin
	}
	return script
}

// -----------------------------------------------------------------------------

// _Latin represents the Latin or Roman script. Most langauges Hangulize
// supports use this script system. So it's the default script.
type _Latin struct{}

// Is checks whether the character is Latin or not.
func (_Latin) Is(ch rune) bool {
	return unicode.Is(unicode.Latin, ch)
}

// Normalize converts a Latin character into
// ISO basic Latin lower alphabet [a-z]:
//
//   Pokémon -> pokemon
//
func (_Latin) Normalize(ch rune) rune {
	props := norm.NFD.PropertiesString(string(ch))
	bin := props.Decomposition()
	if len(bin) != 0 {
		ch = rune(bin[0])
	}
	return unicode.ToLower(ch)
}

// TransliteratePunct does nothing.
func (_Latin) TransliteratePunct(punct rune) string {
	return string(punct)
}

// -----------------------------------------------------------------------------

// _Cyrillic represents the Cyrillic script.
//
//   вулкан
//
type _Cyrillic struct{}

// Is checks whether the character is Cyrillic or not.
func (_Cyrillic) Is(ch rune) bool {
	return unicode.Is(unicode.Cyrillic, ch)
}

// Normalize converts character into lower case.
func (_Cyrillic) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// TransliteratePunct does nothing.
func (_Cyrillic) TransliteratePunct(punct rune) string {
	return string(punct)
}

// -----------------------------------------------------------------------------

// _Georgian represents the Georgian script.
//
//   ასომთავრული
//
type _Georgian struct{}

// Is checks whether the character is Georgian or not.
func (_Georgian) Is(ch rune) bool {
	return unicode.Is(unicode.Georgian, ch)
}

// Normalize does nothing. Georgian is unicase, which means, there's only one
// case for each letter.
func (_Georgian) Normalize(ch rune) rune {
	return ch
}

// TransliteratePunct does nothing.
func (_Georgian) TransliteratePunct(punct rune) string {
	return string(punct)
}

// -----------------------------------------------------------------------------

// _Greek represents the Greek script.
//
//   ελληνικά
//
type _Greek struct{}

// Is checks whether the character is Greek or not.
func (_Greek) Is(ch rune) bool {
	return unicode.Is(unicode.Greek, ch)
}

// Normalize converts character into lower case.
func (_Greek) Normalize(ch rune) rune {
	return unicode.ToLower(ch)
}

// TransliteratePunct does nothing.
func (_Greek) TransliteratePunct(punct rune) string {
	return string(punct)
}

// -----------------------------------------------------------------------------

// _Kana represents the Kana script including Hiragana and Katakana.
//
//   ひらがな カタカナ
//
type _Kana struct{}

// Is checks whether the character is either Hiragana or Katakana.
func (_Kana) Is(ch rune) bool {
	return (ch == 'ー' ||
		unicode.Is(unicode.Hiragana, ch) ||
		unicode.Is(unicode.Katakana, ch))
}

// Normalize converts Hiragana to Katakana.
func (_Kana) Normalize(ch rune) rune {
	const (
		hiraganaMin = rune(0x3040)
		hiraganaMax = rune(0x309f)
	)

	if hiraganaMin <= ch && ch <= hiraganaMax {
		// hiragana to katakana
		return ch + 96
	}
	return ch
}

// TransliteratePunct converts a Japanese punctuation to fit in Korean.
func (_Kana) TransliteratePunct(punct rune) string {
	switch punct {
	case '。':
		return ". "
	case '、':
		return ", "
	case '：':
		return ": "
	case '！':
		return "! "
	case '？':
		return "? "
	case '〜':
		return "~"
	case '「':
		return " '"
	case '」':
		return "' "
	case '『':
		return " \""
	case '』':
		return "\" "
	}

	return string(punct)
}

// -----------------------------------------------------------------------------

// _Pinyin represents the Latin script for Chinese Pinyin.
type _Pinyin struct {
	_Latin
}

// Normalize converts a Latin character for Pinyin into ISO basic Latin lower
// alphabet [a-z]. Especially, it converts "ü" to "v":
//
//   lüè -> lve
//
func (s *_Pinyin) Normalize(ch rune) rune {
	switch ch {
	case 'ü', 'Ü':
		return 'v'
	}
	return s._Latin.Normalize(ch)
}
