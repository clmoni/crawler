package linkengine

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestGetLinksSuccess(t *testing.T) {
	tables := []struct {
		htmlMarkUp    string
		host          string
		expectedLinks int
	}{
		{"<html><body><a href='test.com/test'/><body></html>", "test.com", 1},
		{"<html><body><a href='test.com/test'/><a href='test.com/test'/><body></html>", "test.com", 2},
		{"<html><body><a href='/'/><a href='/test'/><body></html>", "test.com", 2},
		{"<html><body><a href='not-test/'/><a href='not-test/test'/><body></html>", "test.com", 2},
		{"<html><body><a href='http://not-test/'/><a href='http://not-test/test'/><body></html>", "http://test.com", 2},
		{"<html><body><a href='https://not-test/'/><a href='https://not-test/test'/><body></html>", "http://test.com", 2},
		{"<html><body><a href='https://not-test/'/><a href='https://not-test/test'/><a href='https://test.com/test'/><body></html>", "https://test.com", 3},
		{"", "https://test.com", 0},
		{"#", "", 0},
	}

	for _, table := range tables {
		r := ioutil.NopCloser(bytes.NewReader([]byte(table.htmlMarkUp)))
		actualLinks := GetLinks(r, table.host)

		if len(actualLinks) != table.expectedLinks {
			t.Errorf("GetLinks returned %d links, expected %d", len(actualLinks), table.expectedLinks)
		}
	}
}
