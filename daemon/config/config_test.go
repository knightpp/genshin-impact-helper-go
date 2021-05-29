package config_test

import (
	"bytes"
	"fmt"
	"helper/daemon/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigOne(t *testing.T) {
	assert := assert.New(t)
	configBytes := []byte(`
	[[Account]]
	Name = "main"
	Cookie = "cookiestr"
	`)
	c, err := config.ReadConfig(bytes.NewReader(configBytes))
	if assert.NoError(err) {
		assert.Equal(1, len(c.Account))
		first := c.Account[0]
		assert.Equal("main", first.Name)
		assert.Equal("cookiestr", first.Cookie)
	}
}

func TestReadConfigMultiple(t *testing.T) {
	assert := assert.New(t)
	configBytes := []byte(`
	[[Account]]
	Name = "1"
	Cookie = "cookiestr1"
	[[Account]]
	Name = "2"
	Cookie = "cookiestr2"
	[[Account]]
	Name = "3"
	Cookie = "cookiestr3"
	`)
	c, err := config.ReadConfig(bytes.NewReader(configBytes))
	if assert.NoError(err) {
		assert.Equal(3, len(c.Account))
		for i, acc := range c.Account {
			assert.Equal(fmt.Sprint(i+1), acc.Name)
			assert.Equal(fmt.Sprintf("cookiestr%d", i+1), acc.Cookie)
		}
	}
}
