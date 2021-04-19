package crawlerengine

import (
	"net/http"
)

// func TestGetLinkSuccess(t *testing.T) {
// 	html := "<html><body><a href='test.com/test'/><body></html>"
// 	r := ioutil.NopCloser(bytes.NewReader([]byte(html)))

// 	client := &MockClient{
// 		MockGet: func(*http.Request) (*http.Response, error) {
// 			return &http.Response{
// 				StatusCode: 200,
// 				Body:       r,
// 			}, nil
// 		},
// 	}
// 	queue := make(chan string)
// 	visited := make(map[string]bool)

// 	ce = NewCrawlerEngine(&queue, &client, &visited)

// }

type MockGetType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockGet MockGetType
}

func (m *MockClient) Get(req *http.Request) (*http.Response, error) {
	return m.MockGet(req)
}
