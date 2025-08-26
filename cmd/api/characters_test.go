package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestCreateCharacterHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	payload := map[string]any{
		"name":        "Monkey D. Luffy",
		"age":         19,
		"description": "Captain of the Straw Hat Pirates",
		"origin":      "East Blue",
		"race":        "human",
		"episode":     1,
		"time_skip":   "pre",
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/characters", bytes.NewReader(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	equal(t, res.StatusCode, http.StatusCreated)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response map[string]any
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(response)

	character, ok := response["character"]
	if !ok {
		t.Fatal("expected response to contain 'character' field")
	}

	characterMap := character.(map[string]any)
	equal(t, characterMap["name"], "Monkey D. Luffy")
	equal(t, int(characterMap["age"].(float64)), 19)
	equal(t, characterMap["race"], "Human")

}
