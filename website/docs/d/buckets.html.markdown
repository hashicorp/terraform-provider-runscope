---
layout: "runscope"
page_title: "Runscope: runscope_buckets"
sidebar_current: "docs-runscope-datasource-buckets"
description: |-
  Get information about runscope buckets.
---

# runscope\_buckets

Use this data source to get information about matching [buckets](https://www.runscope.com/docs/api/buckets)
that you can use with other runscope resources.

## Example Usage

```hcl
data "runscope_buckets" "buckets" {
	filter = [
		{
			name = "name"
			values = ["test-bucket"]
		}
	]
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter to reduce the list of buckets returned.

Variables (`filter`) supports the following:

* `name` - The name of the field to filter on, currently either: `key`, `name`.
* `values` - The list of values to match against

## Attributes Reference

The following attributes are exported:

* `keys` - A list of the keys of matching buckets.
