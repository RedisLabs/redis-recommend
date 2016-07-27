package redrec

import (
	"fmt"
	"testing"
)

var rr *Redrec
var err error
var max = 100

func TestNew(t *testing.T) {
	rr, err = New("redis://localhost:6379")
	if err != nil {
		t.Error("Redis init Error", err)
	}
}

func TestRate(t *testing.T) {
	err := rr.Rate("item1", "user1", 0.4)
	if err != nil {
		t.Error("Rate Error", err)
	}

	err = rr.Rate("item2", "user1", 0.5)
	err = rr.Rate("item3", "user1", 0.6)

	err = rr.Rate("item1", "user2", 0.4)
	err = rr.Rate("item2", "user2", 0.5)
	err = rr.Rate("item3", "user2", 0.6)
	err = rr.Rate("item4", "user2", 0.7)
	err = rr.Rate("item5", "user2", 0.8)
	err = rr.Rate("item6", "user2", 0.9)

	err = rr.Rate("item5", "user3", 0.66)
	err = rr.Rate("item6", "user3", 0.66)
	err = rr.Rate("item7", "user3", 0.66)
	err = rr.Rate("item8", "user3", 0.66)

	items, err := rr.getUserItems("user2", max)
	if err != nil {
		t.Error("Rate Error", err)
	}
	if len(items) != 6 {
		t.Error("Rate Items len", len(items))
	}
}

func TestGetSimilarityCandidates(t *testing.T) {
	result, err := rr.getSimilarityCandidates("user1", max)
	if err != nil {
		t.Error("getSimilarityCandidates Error", err)
	}

	if len(result) != 2 {
		t.Error("getSimilarityCandidates result len", len(result))
	}

	result, err = rr.getSimilarityCandidates("user2", max)
	if err != nil {
		t.Error("getSimilarityCandidates Error", err)
	}

	if len(result) != 3 {
		t.Error("getSimilarityCandidates result len", len(result))
	}
}

func TestCalcSimilarity(t *testing.T) {
	result, err := rr.calcSimilarity("user1", "user2")
	if err != nil {
		t.Error("calcSimilarity Error", err)
	}

	if result != 0.0 {
		t.Error("calcSimilarity result", result)
	}
}

func TestGetSuggestCandidates(t *testing.T) {
	rr.BatchUpdateSimilarUsers(max)
	result, err := rr.getSuggestCandidates("user1", max)
	if err != nil {
		t.Error("getSuggestCandidates Error", err)
	}

	if len(result) != 3 {
		t.Error("getSuggestCandidates result len", len(result))
	}

	if result[0] != "item4" {
		t.Error("getSuggestCandidates result[0]", result)
	}
}

func TestGetUserSuggestions(t *testing.T) {
	rr.UpdateSuggestedItems("user1", max)
	result, err := rr.GetUserSuggestions("user1", max)
	if err != nil {
		t.Error("GetUserSuggestions Error", err)
	}

	fmt.Println("result: ", result)
	if result[0] != "item6" {
		t.Error("GetUserSuggestions", result)
	}
}
