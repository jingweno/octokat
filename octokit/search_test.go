package octokit

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchService_Users(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/search/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("user_search.json"))
	})

	searchResults, result := client.Search().Users(nil, M{"query": "dhruvsinghal"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, searchResults.TotalCount, 2)
	assert.Equal(t, len(searchResults.Items), 2)
	assert.Equal(t, searchResults.Items[0].ID, 3338555)
	assert.Equal(t, searchResults.Items[0].Login, "dhruvsinghal")
	assert.Equal(t, searchResults.Items[1].ID, 9294878)
	assert.Equal(t, searchResults.Items[1].Login, "dhruvsinghal5")
}

func TestSearchService_Issues(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter,
		r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("issue_search.json"))
	})

	searchResults, result := client.Search().Issues(nil, M{"query": "color"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, searchResults.TotalCount, 180172)
	assert.Equal(t, searchResults.Items[0].Number, 1551)
	assert.Equal(t, searchResults.Items[0].Title, "Colorizer")
	assert.Equal(t, searchResults.Items[1].Number, 3402)
	assert.Equal(t, searchResults.Items[1].Title, "Colorizer")
}

func TestSearchService_Repositories(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter,
		r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("repository_search.json"))
	})

	searchResults, result := client.Search().Repositories(nil,
		M{"query": "asdfghjk"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, searchResults.TotalCount, 21)
	assert.Equal(t, len(searchResults.Items), 21)
	assert.Equal(t, searchResults.Items[0].ID, 8269299)
	assert.Equal(t, searchResults.Items[0].FullName, "ysai/asdfghjk")
	assert.Equal(t, searchResults.Items[1].ID, 8511889)
	assert.Equal(t, searchResults.Items[1].FullName, "ines949494/ikadasd")
}

func TestSearchService_Code(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/search/code", func(w http.ResponseWriter,
		r *http.Request) {
		testMethod(t, r, "GET")
		respondWithJSON(w, loadFixture("code_search.json"))
	})

	searchResults, result := client.Search().Code(nil, M{
		"query": "addClass in:file language:js repo:jquery/jquery"})

	assert.False(t, result.HasError())
	assert.False(t, searchResults.IncompleteResults)
	assert.Equal(t, searchResults.TotalCount, 7)
	assert.Equal(t, len(searchResults.Items), 7)
	assert.Equal(t, searchResults.Items[0].Name, "classes.js")
	assert.Equal(t, searchResults.Items[0].Path, "src/attributes/classes.js")
	assert.Equal(t, searchResults.Items[0].SHA,
		"f9dba94f7de43d6b6b7256e05e0d17c4741a4cde")
	assert.Equal(t, string(searchResults.Items[0].URL),
		"https://api.github.com/repositories/167174/contents/src/attributes/classes.js?ref=53aa87f3bf4284763405f3eb8affff296e55ba4f")
	assert.Equal(t, searchResults.Items[0].GitURL,
		"https://api.github.com/repositories/167174/git/blobs/f9dba94f7de43d6b6b7256e05e0d17c4741a4cde")
	assert.Equal(t, searchResults.Items[0].HTMLURL,
		"https://github.com/jquery/jquery/blob/53aa87f3bf4284763405f3eb8affff296e55ba4f/src/attributes/classes.js")
	assert.Equal(t, searchResults.Items[0].Repository.ID, 167174)
	assert.Equal(t, searchResults.Items[0].Repository.FullName,
		"jquery/jquery")
}
