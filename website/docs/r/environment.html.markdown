---
layout: "runscope"
page_title: "Runscope: runscope_environment"
sidebar_current: "docs-runscope-resource-environment"
description: |-
  Provides a Runscope environment resource.
---

# runscope\_environment

An [environment](https://www.runscope.com/docs/api/environments) resource.
An [environment](https://www.runscope.com/docs/api-testing/environments)
is is a group of configuration settings (initial variables, locations,
notifications, integrations, etc.) used when running a test.
Every test has at least one environment, but you can create additional
environments as needed. For common settings (base URLs, API keys)
that you'd like to use across all tests within a bucket,
use a [Shared Environment](https://www.runscope.com/docs/api-testing/environments#shared).

### Creating a shared environment

> Note: to create a shared environment you do not include a `test_id`

```hcl
resource "runscope_environment" "environment" {
  bucket_id    = "${runscope_bucket.bucket.id}"
  name         = "shared-environment"

  integrations = [
    "${data.runscope_integration.pagerduty.id}"
  ]

  initial_variables {
    var1 = "true",
    var2 = "value2"
  }
}

data "runscope_integration" "pagerduty" {
  team_uuid = "%s"
  type = "pagerduty"
}
```
### Creating a test environment

> Note: to create an environment specific to a test include the associated `test_id`

```hcl
resource "runscope_environment" "environment" {
  bucket_id    = "${runscope_bucket.bucket.id}"
  test_id      = "${runscope_test.api.id}
  name         = "test-environment"

  integrations = [ 
    "${data.runscope_integration.pagerduty.id}"
  ]

  initial_variables {
    var1 = "true",
    var2 = "value2"
  }
}

data "runscope_integration" "pagerduty" {
  team_uuid = "194204f3-19a3-4ef7-a492-b14a277025da"
  type = "pagerduty"
}

# Add a test to a bucket
resource "runscope_test" "api" {
  name         = "api-test"
  description  = "checks the api is up and running"
  bucket_id    = "${runscope_bucket.main}"
}

# Create a bucket
resource "runscope_bucket" "main" {
  name         = "terraform-ftw"
  team_uuid    = "870ed937-bc6e-4d8b-a9a5-d7f9f2412fa3"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_id` - (Required) The id of the bucket to associate this environment with.
* `test_id` - (Optional) The id of the test to associate this environment with.
If given, creates a test specific environment, otherwise creates a shared environment.
* `name` - (Required) The name of environment.
* `script` - (Optional) The [script](https://www.runscope.com/docs/api-testing/scripts#initial-script)
to to run to setup the environment
* `preserve_cookies` - (Optional) If this is set to true, tests using this enviornment will manage cookies between steps.
* `initial_variables` - (Optional) Map of keys and values being used for variables when the test begins.
* `integrations` - (Optional) A list of integration ids to enable for test runs using this environment.
* `regions` - (Optional) A list of [Runscope regions](https://www.runscope.com/docs/regions) to execute test runs in when using this environment.
* `remote_agents` - (Optional) A list of [Remote Agents](https://www.runscope.com/docs/api/agents) to execute test runs in when using this environment.
Remote Agents documented below.

Remote Agents (`remote_agents`) supports the following:

* `name` - (Required) The name of the remote agent
* `uuid` - (Required) The uuid of the remote agent

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the environment.
