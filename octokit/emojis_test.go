package octokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootEmojisService_All(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/emojis", "emojis", nil)

	url, err := EmojisURL.Expand(nil)
	assert.NoError(t, err)

	emojis, result := client.Emojis(url).All()
	assert.False(t, result.HasError())
	var penguin = "https://github.global.ssl.fastly.net/images/icons/emoji/penguin.png?v5"
	var metal = "https://github.global.ssl.fastly.net/images/icons/emoji/metal.png?v5"
	assert.EqualValues(t, penguin, emojis["penguin"])
	assert.EqualValues(t, metal, emojis["metal"])
}
