package runscope

import (
	"errors"
	"fmt"
)

// TestStep represents each step that makes up part of the test. See https://www.runscope.com/docs/api/steps
type TestStep struct {
	URL           string                 `json:"url,omitempty"`
	Variables     []*Variable            `json:"variables,omitempty"`
	Args          map[string]interface{} `json:"args,omitempty"`
	StepType      string                 `json:"step_type,omitempty"`
	Auth          map[string]string      `json:"auth,omitempty"`
	ID            string                 `json:"id,omitempty"`
	Body          string                 `json:"body,omitempty"`
	Note          string                 `json:"note,omitempty"`
	Headers       map[string][]string    `json:"headers,omitempty"`
	RequestID     string                 `json:"request_id,omitempty"`
	Assertions    []*Assertion           `json:"assertions,omitempty"`
	Scripts       []string               `json:"scripts,omitempty"`
	BeforeScripts []string               `json:"before_scripts,omitempty"`
	Method        string                 `json:"method,omitempty"`
}

// NewTestStep creates a new test step struct
func NewTestStep() *TestStep {
	return &TestStep{}
}

// CreateTestStep creates a new runscope test step. See https://www.runscope.com/docs/api/steps#add
func (client *Client) CreateTestStep(testStep *TestStep, bucketKey string, testID string) (*TestStep, error) {
	if error := testStep.validate(); error != nil {
		return nil, error
	}

	client.Lock()
	defer client.Unlock()
	newResource, error := client.createResource(testStep, "test step", testStep.ID,
		fmt.Sprintf("/buckets/%s/tests/%s/steps", bucketKey, testID))
	if error != nil {
		return nil, error
	}

	steps := newResource.Data.([]interface{})
	step := steps[len(steps)-1].(map[string]interface{})
	newTestStep, error := getTestStepFromResponse(step)
	if error != nil {
		return nil, error
	}

	return newTestStep, nil
}

// ReadTestStep list details about an existing test step. https://www.runscope.com/docs/api/steps#detail
func (client *Client) ReadTestStep(testStep *TestStep, bucketKey string, testID string) (*TestStep, error) {
	resource, error := client.readResource("test step", testStep.ID,
		fmt.Sprintf("/buckets/%s/tests/%s/steps/%s", bucketKey, testID, testStep.ID))
	if error != nil {
		return nil, error
	}

	readTestStep, error := getTestStepFromResponse(resource.Data)
	if error != nil {
		return nil, error
	}

	return readTestStep, nil
}

// UpdateTestStep updates an existing test step. https://www.runscope.com/docs/api/steps#modify
func (client *Client) UpdateTestStep(testStep *TestStep, bucketKey string, testID string) (*TestStep, error) {
	resource, error := client.updateResource(testStep, "test step", testStep.ID,
		fmt.Sprintf("/buckets/%s/tests/%s/steps/%s", bucketKey, testID, testStep.ID))
	if error != nil {
		return nil, error
	}

	readTestStep, error := getTestStepFromResponse(resource.Data)
	if error != nil {
		return nil, error
	}

	return readTestStep, nil
}

// DeleteTestStep delete an existing test step. https://www.runscope.com/docs/api/steps#delete
func (client *Client) DeleteTestStep(testStep *TestStep, bucketKey string, testID string) error {
	return client.deleteResource("test step", testStep.ID,
		fmt.Sprintf("/buckets/%s/tests/%s/steps/%s", bucketKey, testID, testStep.ID))
}

func getTestStepFromResponse(response interface{}) (*TestStep, error) {
	testStep := new(TestStep)
	err := decode(testStep, response)
	return testStep, err
}

func (step *TestStep) validate() error {
	if step.StepType == "request" {
		if err := step.validateRequestType(); err != nil {
			return err
		}
	}

	return nil
}

func (step *TestStep) validateRequestType() error {
	if step.Method == "" {
		return errors.New("A request test step must specify 'Method' property")
	}

	if step.Method == "GET" && step.Body != "" {
		return errors.New("A request test step that specifies a 'GET' method can not include a body property")
	}

	return nil
}
