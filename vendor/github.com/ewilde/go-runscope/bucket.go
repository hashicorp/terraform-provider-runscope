/*
Package runscope implements a client library for the runscope api (https://www.runscope.com/docs/api)

*/
package runscope

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

const (
	// DefaultPageSize is the max number of items fetched in each request
	DefaultPageSize = 10
)

// Bucket resources are a simple way to organize your requests and tests. See https://www.runscope.com/docs/api/buckets and https://www.runscope.com/docs/buckets
type Bucket struct {
	Name           string `json:"name,omitempty"`
	Key            string `json:"key,omitempty"`
	Default        bool   `json:"default,omitempty"`
	AuthToken      string `json:"auth_token,omitempty"`
	TestsURL       string `json:"tests_url,omitempty" mapstructure:"tests_url"`
	CollectionsURL string `json:"collections_url,omitempty"`
	MessagesURL    string `json:"messages_url,omitempty"`
	TriggerURL     string `json:"trigger_url,omitempty"`
	VerifySsl      bool   `json:"verify_ssl,omitempty"`
	Team           *Team  `json:"team,omitempty"`
}

// CreateBucket creates a new bucket resource. See https://www.runscope.com/docs/api/buckets#bucket-create
func (client *Client) CreateBucket(bucket *Bucket) (*Bucket, error) {
	DebugF(1, "creating bucket %s", bucket.Name)
	data := url.Values{}
	data.Add("name", bucket.Name)
	data.Add("team_uuid", bucket.Team.ID)

	DebugF(2, "	request: POST %s %#v", "/buckets", data)

	req, err := client.newFormURLEncodedRequest("POST", "/buckets", data)
	if err != nil {
		return nil, err
	}

	DebugF(2, "%#v", req)
	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	DebugF(2, "	response: %d %s", resp.StatusCode, bodyString)

	if resp.StatusCode >= 300 {
		errorResp := new(errorResponse)
		if err = json.Unmarshal(bodyBytes, &errorResp); err != nil {
			return nil, fmt.Errorf("Error creating bucket: %s", bucket.Name)
		}

		return nil, fmt.Errorf("Error creating bucket: %s, status: %d reason: %q", bucket.Name,
			errorResp.Status, errorResp.ErrorMessage)

	}

	response := new(response)
	json.Unmarshal(bodyBytes, &response)
	return getBucketFromResponse(response.Data)
}

// ReadBucket list details about an existing bucket resource. See https://www.runscope.com/docs/api/buckets#bucket-list
func (client *Client) ReadBucket(key string) (*Bucket, error) {
	resource, err := client.readResource("bucket", key, fmt.Sprintf("/buckets/%s", key))
	if err != nil {
		return nil, err
	}

	bucket, err := getBucketFromResponse(resource.Data)
	return bucket, err
}

// DeleteBucket deletes a bucket by key. See https://www.runscope.com/docs/api/buckets#bucket-delete
func (client *Client) DeleteBucket(key string) error {
	return client.deleteResource("bucket", key, fmt.Sprintf("/buckets/%s", key))
}

// DeleteBuckets deletes all buckets matching the predicate
func (client *Client) DeleteBuckets(predicate func(bucket *Bucket) bool) error {

	buckets, err := client.ListBuckets()
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		if predicate(bucket) {
			client.DeleteBucket(bucket.Key)
		}
	}

	return nil
}

// ListBuckets lists all buckets for an account
func (client *Client) ListBuckets() ([]*Bucket, error) {
	resource, err := client.readResource("[]bucket", "", "/buckets")
	if err != nil {
		return nil, err
	}

	buckets, err := getBucketsFromResponse(resource.Data)
	return buckets, err
}

// ListTestsInput represents the input to ListTests func
type ListTestsInput struct {
	BucketKey string
	Count     int
	Offset    int
}

// ListTests lists some tests given ListTestsInput
func (client *Client) ListTests(input *ListTestsInput) ([]*Test, error) {
	count := input.Count
	if count == 0 {
		count = DefaultPageSize
	}

	resource, err := client.readResource("[]test", "",
		fmt.Sprintf("/buckets/%s/tests?count=%d&offset=%d", input.BucketKey, count, input.Offset))
	if err != nil {
		return nil, err
	}

	tests, err := getTestsFromResponse(resource.Data)
	return tests, err
}

// ListAllTests lists all tests for a bucket
func (client *Client) ListAllTests(input *ListTestsInput) ([]*Test, error) {
	var allTests []*Test
	cfg := &ListTestsInput{
		BucketKey: input.BucketKey,
		Count:     input.Count,
	}

	if cfg.Count == 0 {
		cfg.Count = DefaultPageSize
	}

	for cfg.Offset = 0; ; cfg.Offset += cfg.Count {
		tests, err := client.ListTests(cfg)
		if err != nil {
			return allTests, err
		}

		allTests = append(allTests, tests...)
		if len(tests) < cfg.Count {
			return allTests, nil
		}
	}
}

func (bucket *Bucket) String() string {
	value, err := json.Marshal(bucket)
	if err != nil {
		return ""
	}

	return string(value)
}

func getBucketsFromResponse(response interface{}) ([]*Bucket, error) {
	var buckets []*Bucket
	err := decode(&buckets, response)
	return buckets, err
}

func getBucketFromResponse(response interface{}) (*Bucket, error) {
	bucket := new(Bucket)
	err := decode(bucket, response)
	return bucket, err
}
