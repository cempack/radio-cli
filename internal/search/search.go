package search

import (
	"strings"

	"github.com/cempack/radio-cli/internal/domain"
)

type SearchParams struct {
	Query      string
	Country    string
	Language   string
	Tags       []string
	MinBitrate int
	MaxBitrate int
	Codec      string
}

func FilterStations(stations []domain.Station, params SearchParams) []domain.Station {
	var result []domain.Station
	query := strings.ToLower(params.Query)
	for _, s := range stations {
		if query != "" {
			name := strings.ToLower(s.Name)
			if !strings.Contains(name, query) {
				allTags := strings.ToLower(strings.Join(s.Tags, " "))
				if !strings.Contains(allTags, query) {
					continue
				}
			}
		}
		if params.Country != "" && !strings.EqualFold(s.Country, params.Country) {
			continue
		}
		if params.Language != "" && !strings.EqualFold(s.Language, params.Language) {
			continue
		}
		if len(params.Tags) > 0 {
			matched := false
			stationTags := make(map[string]bool)
			for _, t := range s.Tags {
				stationTags[strings.ToLower(t)] = true
			}
			for _, pt := range params.Tags {
				if stationTags[strings.ToLower(pt)] {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		if params.MinBitrate > 0 && s.Bitrate < params.MinBitrate {
			continue
		}
		if params.MaxBitrate > 0 && s.Bitrate > params.MaxBitrate {
			continue
		}
		if params.Codec != "" && !strings.EqualFold(s.Codec, params.Codec) {
			continue
		}
		result = append(result, s)
	}
	return result
}

func RelatedScore(a, b domain.Station) int {
	if a.ID == b.ID {
		return -1
	}
	score := 0
	if strings.EqualFold(a.Country, b.Country) && a.Country != "" {
		score += 3
	}
	if strings.EqualFold(a.Language, b.Language) && a.Language != "" {
		score += 3
	}
	bTags := make(map[string]bool)
	for _, t := range b.Tags {
		bTags[strings.ToLower(t)] = true
	}
	for _, t := range a.Tags {
		if bTags[strings.ToLower(t)] {
			score += 2
		}
	}
	aTokens := strings.Fields(strings.ToLower(a.Name))
	bTokens := strings.Fields(strings.ToLower(b.Name))
	bTokenSet := make(map[string]bool)
	for _, t := range bTokens {
		bTokenSet[t] = true
	}
	for _, t := range aTokens {
		if len(t) > 2 && bTokenSet[t] {
			score += 2
		}
	}
	if strings.EqualFold(a.Codec, b.Codec) && a.Codec != "" {
		score += 1
	}
	diff := a.Bitrate - b.Bitrate
	if diff < 0 {
		diff = -diff
	}
	if a.Bitrate > 0 && b.Bitrate > 0 && diff <= 64 {
		score += 1
	}
	return score
}

func FindRelated(target domain.Station, stations []domain.Station, n int) []domain.Station {
	type scored struct {
		station domain.Station
		score   int
	}
	var candidates []scored
	for _, s := range stations {
		sc := RelatedScore(target, s)
		if sc > 0 {
			candidates = append(candidates, scored{s, sc})
		}
	}
	// simple insertion sort (n is small)
	for i := 1; i < len(candidates); i++ {
		for j := i; j > 0 && candidates[j].score > candidates[j-1].score; j-- {
			candidates[j], candidates[j-1] = candidates[j-1], candidates[j]
		}
	}
	result := make([]domain.Station, 0, n)
	for i, c := range candidates {
		if i >= n {
			break
		}
		result = append(result, c.station)
	}
	return result
}
