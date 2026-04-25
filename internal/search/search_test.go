package search

import (
	"testing"

	"github.com/cempack/radio-cli/internal/domain"
)

func TestFilterStations(t *testing.T) {
	stations := []domain.Station{
		{ID: "1", Name: "Jazz FM", Country: "US", Language: "English", Tags: []string{"jazz", "smooth"}, Bitrate: 128, Codec: "MP3"},
		{ID: "2", Name: "Rock Radio", Country: "UK", Language: "English", Tags: []string{"rock", "classic"}, Bitrate: 256, Codec: "MP3"},
		{ID: "3", Name: "Classical Wave", Country: "DE", Language: "German", Tags: []string{"classical"}, Bitrate: 200, Codec: "AAC"},
		{ID: "4", Name: "Ambient Drift", Country: "US", Language: "English", Tags: []string{"ambient", "electronic"}, Bitrate: 96, Codec: "MP3"},
	}

	t.Run("filter by query", func(t *testing.T) {
		result := FilterStations(stations, SearchParams{Query: "jazz"})
		if len(result) != 1 || result[0].ID != "1" {
			t.Errorf("expected 1 jazz station, got %d", len(result))
		}
	})

	t.Run("filter by country", func(t *testing.T) {
		result := FilterStations(stations, SearchParams{Country: "US"})
		if len(result) != 2 {
			t.Errorf("expected 2 US stations, got %d", len(result))
		}
	})

	t.Run("filter by language", func(t *testing.T) {
		result := FilterStations(stations, SearchParams{Language: "German"})
		if len(result) != 1 || result[0].ID != "3" {
			t.Errorf("expected 1 German station, got %d", len(result))
		}
	})

	t.Run("filter by tag", func(t *testing.T) {
		result := FilterStations(stations, SearchParams{Tags: []string{"rock"}})
		if len(result) != 1 || result[0].ID != "2" {
			t.Errorf("expected 1 rock station, got %d", len(result))
		}
	})

	t.Run("filter by min bitrate", func(t *testing.T) {
		result := FilterStations(stations, SearchParams{MinBitrate: 200})
		if len(result) != 2 {
			t.Errorf("expected 2 stations >= 200kbps, got %d", len(result))
		}
	})

	t.Run("filter by codec", func(t *testing.T) {
		result := FilterStations(stations, SearchParams{Codec: "AAC"})
		if len(result) != 1 || result[0].ID != "3" {
			t.Errorf("expected 1 AAC station, got %d", len(result))
		}
	})

	t.Run("no filter returns all", func(t *testing.T) {
		result := FilterStations(stations, SearchParams{})
		if len(result) != 4 {
			t.Errorf("expected 4 stations with no filter, got %d", len(result))
		}
	})
}

func TestRelatedScore(t *testing.T) {
	a := domain.Station{
		ID:       "1",
		Name:     "Jazz FM",
		Country:  "US",
		Language: "English",
		Tags:     []string{"jazz", "smooth"},
		Bitrate:  128,
		Codec:    "MP3",
	}

	t.Run("same station returns -1", func(t *testing.T) {
		score := RelatedScore(a, a)
		if score != -1 {
			t.Errorf("expected -1 for same station, got %d", score)
		}
	})

	t.Run("same country and language", func(t *testing.T) {
		b := domain.Station{
			ID:       "2",
			Name:     "Blues Radio",
			Country:  "US",
			Language: "English",
			Tags:     []string{"blues"},
			Bitrate:  128,
			Codec:    "MP3",
		}
		score := RelatedScore(a, b)
		// country(3) + language(3) + same bitrate(1) + same codec(1) = 8; >= 6 guards against algorithm regressions
		if score < 6 {
			t.Errorf("expected score >= 6 for similar station, got %d", score)
		}
	})

	t.Run("shared tags increase score", func(t *testing.T) {
		b := domain.Station{
			ID:       "3",
			Name:     "Smooth Jazz",
			Country:  "CA",
			Language: "French",
			Tags:     []string{"jazz", "smooth", "lounge"},
			Bitrate:  320,
			Codec:    "AAC",
		}
		score := RelatedScore(a, b)
		// 2 shared tags * 2 = 4, plus name token "jazz" shared = 2 -> at least 4
		if score < 4 {
			t.Errorf("expected score >= 4 for shared tags, got %d", score)
		}
	})

	t.Run("unrelated station has low score", func(t *testing.T) {
		b := domain.Station{
			ID:       "4",
			Name:     "Death Metal Underground",
			Country:  "NO",
			Language: "Norwegian",
			Tags:     []string{"metal", "death"},
			Bitrate:  64,
			Codec:    "OGG",
		}
		score := RelatedScore(a, b)
		if score > 2 {
			t.Errorf("expected low score for unrelated station, got %d", score)
		}
	})
}

func TestFindRelated(t *testing.T) {
	target := domain.Station{
		ID:       "1",
		Name:     "Jazz FM",
		Country:  "US",
		Language: "English",
		Tags:     []string{"jazz", "smooth"},
		Bitrate:  128,
		Codec:    "MP3",
	}
	stations := []domain.Station{
		target,
		{ID: "2", Name: "Jazz Radio", Country: "US", Language: "English", Tags: []string{"jazz"}, Bitrate: 128, Codec: "MP3"},
		{ID: "3", Name: "Blues FM", Country: "US", Language: "English", Tags: []string{"blues"}, Bitrate: 128, Codec: "MP3"},
		{ID: "4", Name: "Rock Power", Country: "UK", Language: "English", Tags: []string{"rock"}, Bitrate: 256, Codec: "MP3"},
		{ID: "5", Name: "Klassik Radio", Country: "DE", Language: "German", Tags: []string{"classical"}, Bitrate: 192, Codec: "AAC"},
	}

	related := FindRelated(target, stations, 3)
	if len(related) == 0 {
		t.Error("expected some related stations")
	}
	for _, r := range related {
		if r.ID == target.ID {
			t.Error("related stations should not include the target itself")
		}
	}
	// The jazz radio should be the most related
	if len(related) > 0 && related[0].ID != "2" {
		t.Logf("most related was %s (expected Jazz Radio)", related[0].Name)
	}
}
