package radio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cempack/radio-cli/internal/domain"
)

const baseURL = "https://de1.api.radio-browser.info/json"

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

type apiStation struct {
	UUID            string `json:"stationuuid"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	URLResolved     string `json:"url_resolved"`
	Homepage        string `json:"homepage"`
	Favicon         string `json:"favicon"`
	Country         string `json:"country"`
	CountryCode     string `json:"countrycode"`
	State           string `json:"state"`
	Language        string `json:"language"`
	Tags            string `json:"tags"`
	Codec           string `json:"codec"`
	Bitrate         int    `json:"bitrate"`
	Votes           int    `json:"votes"`
	ClickCount      int    `json:"clickcount"`
	LastCheckOK     int    `json:"lastcheckok"`
	LastCheckTime   string `json:"lastchecktime"`
	LastCheckStatus string `json:"lastcheckoktime"`
}

func (a apiStation) toDomain() domain.Station {
	var tags []string
	for _, t := range strings.Split(a.Tags, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	var checkTime time.Time
	if a.LastCheckTime != "" {
		checkTime, _ = time.Parse("2006-01-02 15:04:05", a.LastCheckTime)
	}
	return domain.Station{
		ID:              a.UUID,
		Name:            a.Name,
		StreamURL:       a.URL,
		ResolvedURL:     a.URLResolved,
		HomepageURL:     a.Homepage,
		FaviconURL:      a.Favicon,
		Country:         a.Country,
		CountryCode:     a.CountryCode,
		State:           a.State,
		Language:        a.Language,
		Tags:            tags,
		Codec:           a.Codec,
		Bitrate:         a.Bitrate,
		Votes:           a.Votes,
		ClickCount:      a.ClickCount,
		LastCheckOK:     a.LastCheckOK == 1,
		LastCheckTime:   checkTime,
		LastCheckStatus: a.LastCheckStatus,
	}
}

func (c *Client) get(path string, params url.Values) ([]domain.Station, error) {
	u := baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "radiodrift/1.0")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}
	var apiStations []apiStation
	if err := json.NewDecoder(resp.Body).Decode(&apiStations); err != nil {
		return nil, err
	}
	stations := make([]domain.Station, 0, len(apiStations))
	for _, a := range apiStations {
		stations = append(stations, a.toDomain())
	}
	return stations, nil
}

type SearchParams struct {
	Name     string
	Country  string
	Language string
	Tags     string
	Limit    int
	Offset   int
	Order    string
	Reverse  bool
}

func (c *Client) Search(p SearchParams) ([]domain.Station, error) {
	params := url.Values{}
	if p.Name != "" {
		params.Set("name", p.Name)
	}
	if p.Country != "" {
		params.Set("country", p.Country)
	}
	if p.Language != "" {
		params.Set("language", p.Language)
	}
	if p.Tags != "" {
		params.Set("tagList", p.Tags)
	}
	limit := 50
	if p.Limit > 0 {
		limit = p.Limit
	}
	params.Set("limit", strconv.Itoa(limit))
	if p.Offset > 0 {
		params.Set("offset", strconv.Itoa(p.Offset))
	}
	order := "votes"
	if p.Order != "" {
		order = p.Order
	}
	params.Set("order", order)
	if p.Reverse {
		params.Set("reverse", "true")
	}
	return c.get("/stations/search", params)
}

func (c *Client) ByTag(tag string) ([]domain.Station, error) {
	return c.get("/stations/bytag/"+url.PathEscape(tag), nil)
}

func (c *Client) ByCountry(country string) ([]domain.Station, error) {
	return c.get("/stations/bycountry/"+url.PathEscape(country), nil)
}

func (c *Client) Random(limit int) ([]domain.Station, error) {
	if limit <= 0 {
		limit = 20
	}
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	return c.get("/stations/random", params)
}
