package furigana

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hangulize/hangulize"
)

func TestPronouncer(t *testing.T) {
	assert.Implements(t, (*hangulize.Pronouncer)(nil), &P)
}

func TestKana(t *testing.T) {
	assert.Equal(t, "ない", P.Pronounce("ない"))
	assert.Equal(t, "ゲーム", P.Pronounce("ゲーム"))
}

func TestPersons(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", P.Pronounce("新海誠"))
	assert.Equal(t, "マツモト ユキヒロ", P.Pronounce("松本行弘"))
}

func TestPlaces(t *testing.T) {
	assert.Equal(t, "トウキョウ", P.Pronounce("東京"))
}

func TestDeduplicateSpaces(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", P.Pronounce("新海 誠"))
	assert.Equal(t, "マツモト ユキヒロ", P.Pronounce("松本 行弘"))
}