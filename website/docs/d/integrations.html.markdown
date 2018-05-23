---
layout: "runscope"
page_title: "Runscope: runscope_integrations"
sidebar_current: "docs-runscope-datasource-integrations"
description: |-
  Get information about runscope integrations enabled on for your team.
---

# runscope\_integration

Use this data source to list all of your [integrations](https://www.runscope.com/docs/api/integrations)
that you can use with other runscope resources.

## Example Usage

```hcl
data "runscope_integrations" "slack" {
	team_uuid = "d26553c0-3537-40a8-9d3c-64b0453262a9"
	filter = {
		name = "type"
		values = ["slack"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter to reduce the list of integrations returned.

Variables (`filter`) supports the following:

* `name` - The name of the field to filter on, currently either: `id`, `type` or `description`.
* `values` - The list of values to match against

## Attributes Reference
The following attributes are exported:

* `id` - The unique identifier of the found integration.
