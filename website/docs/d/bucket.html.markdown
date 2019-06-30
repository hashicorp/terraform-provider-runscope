---
layout: "runscope"
page_title: "Runscope: runscope_bucket"
sidebar_current: "docs-runscope-datasource-bucket"
description: |-
  Get information about a single runscope bucket.
---

# runscope\_bucket

Use this data source to get information about a specific [bucket](https://www.runscope.com/docs/api/buckets)
that you can use with other runscope resources.

## Example Usage

```hcl
data "runscope_bucket" "website" {
  key = "t2f4bkvnggct"
}

resource "runscope_environment" "environment" {
  bucket_id = "${runscope_bucket.website.id}"
  name      = "test-environment"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The unique key of the bucket.

## Attributes Reference

The following attributes are exported:

* `id` - The unique key of the found bucket.
* `team_uuid` - The team unique identifier that owns the bucket.
* `name` - Type name of the bucket.
