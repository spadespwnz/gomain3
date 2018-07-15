package models

type (
	Response struct {
		Message string `json:"message"`
		Error   int    `json"error"`
	}
	Lists struct {
		Lists []ListMeta `json:"lists"`
	}
	ListMeta struct {
		Title string `json:"title"`
		Id    string `json:"id"`
	}
	JPWord struct {
		Type    string `json:"Type"`
		State   string `json:"state"`
		Romaji  string `json:"romaji"`
		Kana    string `json:"kana"`
		Kanji   string `json:"kanji"`
		Meaning string `json:"meaning"`
	}
	JPConj struct {
		Type        string `json:"Type"`
		State       string `json:"state"`
		Romaji      string `json:"romaji"`
		Kana        string `json:"kana"`
		Kanji       string `json:"kanji"`
		Meaning     string `json:"meaning"`
		ConjRomaji  string `json:"conj_romaji"`
		ConjKana    string `json:"conj_kana"`
		ConjMeaning string `json:"conj_meaning"`
	}
)
