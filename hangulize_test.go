package hangulize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLang generates subtests for bundled lang specs.
func TestLang(t *testing.T) {
	for _, lang := range ListLangs() {
		spec, ok := LoadSpec(lang)

		assert.Truef(t, ok, `failed to load "%s" spec`, lang)

		for _, exm := range spec.Test {
			word := exm[0]
			expected := exm[1]

			t.Run(lang+"/"+word, func(t *testing.T) {
				assertHangulize(t, spec, expected, word)
			})
		}
	}
}

// -----------------------------------------------------------------------------
// Edge cases

func hangulize(spec *Spec, word string) string {
	h := NewHangulizer(spec)
	return h.Hangulize(word)
}

// TestSlash tests "/" in input word. The original Hangulize removes the "/" so
// the result was "글로르이아" instead of "글로르/이아".
func TestSlash(t *testing.T) {
	assert.Equal(t, "글로르/이아", Hangulize("ita", "glor/ia"))
	assert.Equal(t, "글로르{}이아", Hangulize("ita", "glor{}ia"))
}

func TestComma(t *testing.T) {
	assertHangulize(t, loadSpec("ita"), "글로르,이아", "glor,ia")
	assertHangulize(t, loadSpec("ita"), "콤,오", "com,o")
}

func TestPunctInVar(t *testing.T) {
	assertHangulize(t, loadSpec("nld"), "빔%", "Wim%")
	assertHangulize(t, loadSpec("cym"), "귀,림", "Gwi,lym")
	assertHangulize(t, loadSpec("wlm"), "카드,고데이", "Cad,Godeu")
}

func TestQuote(t *testing.T) {
	assert.Equal(t, "글로리아", Hangulize("ita", "glor'ia"))
	assert.Equal(t, "코모", Hangulize("ita", "com'o"))
}

func TestSpecials(t *testing.T) {
	assert.Equal(t, "<글로리아>", Hangulize("ita", "<gloria>"))
}

func TestHyphen(t *testing.T) {
	spec := mustParseSpec(`
	transcribe:
		"x" -> "-ㄱㅅ"
		"e-" -> "ㅣ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", hangulize(spec, "ex야!"))
}

func TestDifferentAges(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"x" -> "xx"

	transcribe:
		"xx" -> "-ㄱㅅ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", hangulize(spec, "ex야!"))
}

func TestKeepAndCleanup(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"𐌗"  -> "𐌗𐌗"
		"𐌄𐌗" -> "𐌊-"

	transcribe:
		"𐌊" -> "-ㄱ"
		"𐌗" -> "ㄱㅅ"
	`)
	// ㅋ𐌄 𐌗 !
	// ----│---------------------- rewrite
	//     ├─┐        𐌗->𐌗𐌗
	// ㅋ𐌄 𐌄 𐌗 !
	//   └┬┘
	//   ┌┴┐          𐌄𐌗->𐌊-
	// ㅋ𐌊 - 𐌗 !
	// --│------------------------ transcribe
	//   ├─┐          𐌊->ㄱ
	// ㅋ- ㄱ- 𐌗 !
	//         ├─┐    𐌗->-ㄱㅅ
	// ㅋ- ㄱ- ㄱㅅ!
	// ------│-------------------- cleanup
	//       x
	// ㅋ- ㄱㄱㅅ!
	// --├─┘┌┘┌┘------------------ jamo
	//   │ ┌┘┌┘
	// ㅋ윽그스!
	assert.Equal(t, "ㅋ윽그스!", hangulize(spec, "ㅋ𐌄𐌗!"))
}

func TestSpace(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"van " -> "van/"

	transcribe:
		"van"  -> "반"
		"gogh" -> "고흐"
	`)
	assert.Equal(t, "반고흐", hangulize(spec, "van gogh"))
}

func TestZeroWidthSpace(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"a b" -> "a{}b"
		"^b"  -> "v"

	transcribe:
		"a" -> "ㅇ"
		"b" -> "ㅂ"
		"v" -> "ㅍ"
		"c" -> "ㅊ"
	`)
	assert.Equal(t, "으프 츠", hangulize(spec, "a b c"))
}

func TestVarToVar(t *testing.T) {
	spec := mustParseSpec(`
	vars:
		"abc" = "a", "b", "c"
		"def" = "d", "e", "f"
		"ghi" = "g", "h", "i"

	rewrite:
		"<abc><abc>" -> "<def><ghi>"

	transcribe:
		"a" -> "a"
		"b" -> "b"
		"c" -> "c"
		"d" -> "d"
		"e" -> "e"
		"f" -> "f"
		"g" -> "g"
		"h" -> "h"
		"i" -> "i"
	`)
	assert.Equal(t, "dg", hangulize(spec, "aa"))
	assert.Equal(t, "ei", hangulize(spec, "bc"))
}

func TestSingleDollar(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"$" -> "END"
	transcribe:
		"123" -> "123"
		"END" -> "END"
	`)
	assertHangulize(t, spec, "123END", "123")
}

// -----------------------------------------------------------------------------
// Language-specific edge cases

func TestJpnIgnoresRoman(t *testing.T) {
	assertHangulize(t, loadSpec("jpn"), "abc아", "abcあ")
}

func TestJpnKatakanaLongVowel(t *testing.T) {
	// http://www.textfugu.com/season-3/learn-katakana/2-3/
	jpn := loadSpec("jpn")

	assertHangulize(t, jpn, "위", "ウィ")
	assertHangulize(t, jpn, "웨", "ウェ")
	assertHangulize(t, jpn, "워", "ウォ")

	assertHangulize(t, jpn, "바", "ヴァ")
	assertHangulize(t, jpn, "비", "ヴィ")
	assertHangulize(t, jpn, "브", "ヴ")
	assertHangulize(t, jpn, "베", "ヴェ")
	assertHangulize(t, jpn, "보", "ヴォ")

	assertHangulize(t, jpn, "셰", "シェ")
	assertHangulize(t, jpn, "제", "ジェ")

	assertHangulize(t, jpn, "파", "ファ")
	assertHangulize(t, jpn, "피", "フィ")
	assertHangulize(t, jpn, "페", "フェ")
	assertHangulize(t, jpn, "포", "フォ")

	assertHangulize(t, jpn, "바", "ブァ")
	assertHangulize(t, jpn, "비", "ブィ")
	assertHangulize(t, jpn, "베", "ブェ")
	assertHangulize(t, jpn, "보", "ブォ")

	assertHangulize(t, jpn, "디", "ディ")
	assertHangulize(t, jpn, "도", "ドゥ")
	assertHangulize(t, jpn, "제", "ヂェ")

	assertHangulize(t, jpn, "디티", "ティティ")
	assertHangulize(t, jpn, "도토", "トゥトゥ")
	assertHangulize(t, jpn, "제체", "チェチェ")
}

type stubFurigana struct{}

func (p *stubFurigana) ID() string {
	return "furigana"
}

func (p *stubFurigana) Pronounce(word string) string {
	return "スタブ"
}

func TestInstancePronouncers(t *testing.T) {
	spec, _ := LoadSpec("jpn")
	h := NewHangulizer(spec)
	h.UsePronouncer(&stubFurigana{})
	assert.Equal(t, "스타부", h.Hangulize("1234"))
}

// -----------------------------------------------------------------------------
// Benchmarks

func BenchmarkGloria(b *testing.B) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Hangulize("GLORIA")
	}
}

func BenchmarkGloriaTrace(b *testing.B) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.HangulizeTrace("GLORIA")
	}
}

// -----------------------------------------------------------------------------
// Examples

func Example() {
	// Person names from http://iceager.egloos.com/2610028
	fmt.Println(Hangulize("ron", "Cătălin Moroşanu"))
	fmt.Println(Hangulize("nld", "Jerrel Venetiaan"))
	fmt.Println(Hangulize("por", "Vítor Constâncio"))
	// Output:
	// 커털린 모로샤누
	// 예럴 페네티안
	// 비토르 콘스탄시우
}

func ExampleHangulize_gloria() {
	fmt.Println(Hangulize("ita", "gloria"))
	// Output: 글로리아
}

func ExampleHangulize_nietzsche() {
	fmt.Println(Hangulize("deu", "Friedrich Wilhelm Nietzsche"))
	// Output: 프리드리히 빌헬름 니체
}

func ExampleHangulize_shinkaiMakoto() {
	// import "github.com/hangulize/hangulize/pronounce/furigana"
	// UsePronouncer(&furigana.P)

	fmt.Println(Hangulize("jpn", "新海誠"))
	// Output: 신카이 마코토
}

func ExampleNewHangulizer() {
	spec, _ := LoadSpec("nld")
	h := NewHangulizer(spec)

	fmt.Println(h.Hangulize("Vincent van Gogh"))
	// Output: 빈센트 반고흐
}
