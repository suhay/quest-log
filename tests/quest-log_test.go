package tests

import (
	"testing"

	quest "github.com/suhay/quest-log/query"
)

func TestThread(t *testing.T) {
	var params = quest.Params{
		Query: `{ thread { name } }`,
	}
	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		t.Errorf("Response had errors")
	} else {
		if results.Thread.Name == "" {
			t.Errorf("Response was blank")
		}
	}
}

func TestGetThread(t *testing.T) {
	var params = quest.Params{
		Query: `{ GetThread(name: "Power Low") { name } }`,
	}

	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		t.Errorf("Response had errors")
	} else {
		if results.GetThread.Name != "Power Low" {
			t.Errorf("Response name was incorrect")
		}
	}
}
