[![Build Status](https://travis-ci.org/ewilde/terraform-provider-runscope.svg?branch=master)](https://travis-ci.org/ewilde/terraform-provider-runscope)

Terraform Runscope Provider
===========================

This repository contains a plugin form of the Runscope provider that was proposed
and submitted in [Terraform PR #14221][1].

The Runscope provider is used to create and manage Runscope tests using
the official [Runscope API][2]

## Installing

See the [Plugin Basics][4] page of the Terraform docs to see how to plunk this
into your config. Check the [releases page][5] of this repo to get releases for
Linux, OS X, and Windows.

## Releasing
Releases are automatically setup to go out from the master branch after a build is made on master with a tag.

To perform a release simply create a tag:
` git tag -a v0.0.2 -m "Release message"`

Then push your tag:
`git push origin v0.0.2`


That's it, the build will now run and create a new release on [github](https://github.com/form3tech/ewilde/terraform-provider-runscope) :

## Usage

The following section details the use of the provider and its resources.

These docs are derived from the middleman templates that were created for the
PR itself, and can be found in their original form [here][5].

### Example Usage

The below example is an end-to-end demonstration of the setup of a basic
runscope test:


### Creating a test with a step
```hcl
resource "runscope_step" "main_page" {
  bucket_id      = "${runscope_bucket.bucket.id}"
  test_id        = "${runscope_test.test.id}"
  step_type      = "request"
  url            = "http://example.com"
  method         = "GET"
  variables      = [
  	{
  	   name     = "httpStatus"
  	   source   = "response_status"
  	},
  	{
  	   name     = "httpContentEncoding"
  	   source   = "response_header"
  	   property = "Content-Encoding"
  	},
  ]
  assertions     = [
  	{
  	   source     = "response_status"
       comparison = "equal_number"
       value      = "200"
  	},
  	{
  	   source     = "response_json"
       comparison = "equal"
       value      = "c5baeb4a-2379-478a-9cda-1b671de77cf9",
       property   = "data.id"
  	},
  ],
  headers        = [
  	{
  		header = "Accept-Encoding",
  		value  = "application/json"
  	},
  	{
  		header = "Accept-Encoding",
  		value  = "application/xml"
  	},
  	{
  		header = "Authorization",
  		value  = "Bearer bb74fe7b-b9f2-48bd-9445-bdc60e1edc6a",
	}
  ]
}

resource "runscope_test" "test" {
  bucket_id   = "${runscope_bucket.bucket.id}"
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_bucket" "bucket" {
  name      = "terraform-provider-test"
  team_uuid = "dfb75aac-eeb3-4451-8675-3a37ab421e4f"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_id` - (Required) The id of the bucket to associate this step with.
* `test_id` - (Required) The id of the test to associate this step with.
* `step_type` - (Required) The type of step.
 * [request](#request-steps)
 * pause
 * condition
 * ghost
 * subtest

### Request steps
When creating a `request` type of step the additional arguments also apply:

* `method` - (Required) The HTTP method for this request step.
* `variables` - (Optional) A list of variables to extract out of the HTTP response from this request. Variables documented below.
* `assertions` - (Optional) A list of assertions to apply to the HTTP response from this request. Assertions documented below.
* `headers` - (Optional) A list of headers to apply to the request. Headers documented below.
* `body` - (Optional) A string to use as the body of the request.

Variables (`variables`) supports the following:

* `name` - (Required) Name of the variable to define.
* `property` - (Required) The name of the source property. i.e. header name or json path
* `source` - (Required) The variable source, for list of allowed values see: https://www.runscope.com/docs/api/steps#assertions

Assertions (`assertions`) supports the following:

* `source` - (Required) The assertion source, for list of allowed values see: https://www.runscope.com/docs/api/steps#assertions
* `property` - (Optional) The name of the source property. i.e. header name or json path
* `comparison` - (Required) The assertion comparison to make i.e. `equals`, for list of allowed values see: https://www.runscope.com/docs/api/steps#assertions
* `value` - (Optional) The value the `comparison` will use

**Example Assertions**

Status Code == 200

```json
"assertions": [
    {
        "source": "response_status",
        "comparison": "equal_number",
        "value": 200
    }
]
```

JSON element 'address' contains the text "avenue"


```json
"assertions": [
    {
        "source": "response_json",
        "property": "address",
        "comparison": "contains",
        "value": "avenue"
    }
]
```

Response Time is faster than 1 second.


```json
"assertions": [
    {
        "source": "response_time",
        "comparison": "is_less_than",
        "value": 1000
    }
]
```

The `headers` list supports the following:

* `header` - (Required) The name of the header
* `value` - (Required) The name header value

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the step.


[1]: https://github.com/hashicorp/terraform/pull/14221
[2]: https://www.runscope.com/docs/api
[3]: https://www.terraform.io/docs/plugins/basics.html
[4]: https://github.com/ewilde/terraform-provider-runscope/releases
[5]: website/source/docs/providers/runscope

## Developing
### Running the integration tests
`make TF_ACC=1 RUNSCOPE_TEAM_ID=xxx RUNSCOPE_ACCESS_TOKEN=xxx RUNSCOPE_INTEGRATION_DESC="test integration"`