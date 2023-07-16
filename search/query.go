package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	meilisearch "github.com/meilisearch/meilisearch-go"
	es "github.com/opensearch-project/opensearch-go/v2"
)

type EsSearchHit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

type EsSearchHits struct {
	Total struct {
		Value    int
		Relation string
	} `json:"total"`
	MaxScore float64       `json:"max_score"`
	Hits     []EsSearchHit `json:"hits"`
}

type EsSearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	// Shards ???
	Hits EsSearchHits `json:"hits"`
}

type UserResult struct {
	Did    string `json:"did"`
	Handle string `json:"handle"`
}

type PostSearchResult struct {
	Tid  string     `json:"tid"`
	Cid  string     `json:"cid"`
	User UserResult `json:"user"`
	Post any        `json:"post"`
}

type MeiliPost struct {
	DocumentId  string     `json:"documentId"`
	Text        string     `json:"text"`
	CreatedAt   int64      `json:"createdAt"`
	User        string     `json:"user"`
}

func doSearchPosts(ctx context.Context, escli *es.Client, q string, offset int, size int) (*EsSearchResponse, error) {
	query := map[string]interface{}{
		"sort": map[string]any{
			"createdAt": map[string]any{
				"order": "desc",
			},
		},
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"text": map[string]any{
					"query":    q,
					"operator": "and",
				},
			},
		},
		"size": size,
		"from": offset,
	}

	return doSearch(ctx, escli, "posts", query)
}

func doSearchProfiles(ctx context.Context, escli *es.Client, q string) (*EsSearchResponse, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    q,
				"fields":   []string{"description", "displayName", "handle"},
				"operator": "or",
			},
		},
	}

	return doSearch(ctx, escli, "profiles", query)
}

func doSearch(ctx context.Context, escli *es.Client, index string, query interface{}) (*EsSearchResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := escli.Search(
		escli.Search.WithContext(ctx),
		escli.Search.WithIndex(index),
		escli.Search.WithBody(&buf),
		escli.Search.WithTrackTotalHits(true),
		escli.Search.WithSize(30),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var out EsSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding search response: %w", err)
	}

	return &out, nil
}

func doSearchPostsMeili(ctx context.Context, meilicli *meilisearch.Client, q string, offset int, size int) ([]MeiliPost, error) {
	query := &meilisearch.SearchRequest{
		Offset: int64(offset),
		Limit: int64(size),
		Sort: []string{
			"createdAt:desc",
		},
	}

	resp, err := doSearchMeili(ctx, meilicli, "posts", query, q)
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(resp.Hits)
	if err != nil {
		return nil, fmt.Errorf("failed to encode resp.Hits of Meilisearch: %w", err)
	}

	var out []MeiliPost
	if err := json.Unmarshal(encoded, &out); err != nil {
		return nil, fmt.Errorf("failed to decode json of resp.Hits of Meilisearch: %w", err)
	}

	return out, nil
}

func doSearchProfilesMeili(ctx context.Context, meilicli *meilisearch.Client, q string) ([]MeiliProfile, error) {
	query := &meilisearch.SearchRequest{}

	resp, err := doSearchMeili(ctx, meilicli, "profiles", query, q)
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(resp.Hits)
	if err != nil {
		return nil, fmt.Errorf("failed to encode resp.Hits of Meilisearch: %w", err)
	}

	var out []MeiliProfile
	if err := json.Unmarshal(encoded, &out); err != nil {
		return nil, fmt.Errorf("failed to decode json of resp.Hits of Meilisearch: %w", err)
	}

	return out, nil
}

func doSearchMeili(ctx context.Context, meilicli *meilisearch.Client, index string, query *meilisearch.SearchRequest, keyword string) (*meilisearch.SearchResponse, error) {
	resp, err := meilicli.Index(index).Search(keyword, query)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	return resp, nil
}