package tests

import (
	"flag"
	"testing"

	quest "github.com/suhay/quest-log/query"
)

var flagPath = flag.String("path", "", "Path to local threads.")

func TestThread(t *testing.T) {
	var params = quest.Params{
		Query: `{ thread { name } }`,
	}
	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		t.Log(resp.Errors)
		t.Error("Response had errors")
	} else {
		if results.Thread.Name == "" {
			t.Error("Response was blank")
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
		t.Error("Response had errors")
	} else {
		if results.GetThread.Name != "Power Low" {
			t.Error("Response name was incorrect")
		}
	}
}
