package octokit

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer/hypermedia"
	"testing"
)

func TestNewResult_Pageable(t *testing.T) {
	resp := &Response{Relations: hypermedia.Relations{"next": hypermedia.Hyperlink("/path")}}
	result := newResult(resp, nil)

	assert.Equal(t, "/path", string(*result.NextPage))
	assert.T(t, result.PrevPage == nil)
	assert.T(t, result.LastPage == nil)
	assert.T(t, result.FirstPage == nil)
}
