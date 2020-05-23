package tests

import (
	"log"
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
		t.Log(resp.Errors)
		t.Error("Response had errors")
	} else {
		if results.Thread.Name == "" {
			t.Error("Response was blank")
		}
	}
}

func TestThreadByName(t *testing.T) {
	variables := make(map[string]interface{})
	variables["name"] = "Power Low"

	var params = quest.Params{
		Query:     `query ThreadByName($name: String){ thread(name: $name) { name } }`,
		Variables: variables,
	}

	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		t.Error("Response had errors")
	} else {
		if results.Thread.Name != "Power Low" {
			t.Error("Response name was incorrect")
		}
	}
}

func TestThreadByTag(t *testing.T) {
	variables := make(map[string]interface{})
	variables["tags"] = "shield"

	var params = quest.Params{
		Query:     `query ThreadByTag($tags: [String]){ thread(tags: $tags) { name } }`,
		Variables: variables,
	}

	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		log.Fatalf("%v", resp.Errors)
		t.Error("Response had errors")
	} else {
		if results.Thread.Name != "You spin me right round" {
			log.Fatalf("%v", results.Thread.Name)
			t.Error("Response name was incorrect")
		}
	}
}

func TestEntry(t *testing.T) {
	var params = quest.Params{
		Query: `{ entry { name } }`,
	}
	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		t.Log(resp.Errors)
		t.Error("Response had errors")
	} else {
		if results.Entry.Name == "" {
			t.Error("Response was blank")
		}
	}
}

func TestEntryByName(t *testing.T) {
	variables := make(map[string]interface{})
	variables["name"] = "test entry"

	var params = quest.Params{
		Query:     `query EntryByName($name: String){ entry(name: $name) { name } }`,
		Variables: variables,
	}

	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		t.Error("Response had errors")
	} else {
		if results.Entry.Name != "test entry" {
			t.Errorf("Response name was incorrect: %s", results.Entry.Name)
		}
	}
}

func TestEntryByTag(t *testing.T) {
	variables := make(map[string]interface{})
	variables["tags"] = "good-vs-evil"

	var params = quest.Params{
		Query:     `query EntryByTag($tags: [String]){ entry(tags: $tags) { name } }`,
		Variables: variables,
	}

	results, resp, err := quest.Query(params)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Errors) > 0 {
		log.Fatalf("%v", resp.Errors)
		t.Error("Response had errors")
	} else {
		if results.Entry.Name != "test entry" {
			t.Error("Response name was incorrect")
		}
	}
}
