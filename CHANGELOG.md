## 0.7.0 (Unreleased)
## 0.6.0 (June 30, 2019)

NOTES:

* This release includes a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and there should not be any significant behavioural changes. ([#27](https://github.com/terraform-providers/terraform-provider-runscope/issues/27))

## 0.5.0 (October 04, 2018)
ENHANCEMENTS:

*  resource/runscope_step: New attributes `note` added ([#16](https://github.com/terraform-providers/terraform-provider-runscope/pull/16))

## 0.4.0 (September 22, 2018)
ENHANCEMENTS:

*  resource/runscope_environment: New attributes `webhooks` and `emails` added ([#13](https://github.com/terraform-providers/terraform-provider-runscope/pull/13))
## 0.3.0 (July 26, 2018)

FEATURES:

* **New Data Source:** `runscope_bucket` ([#12](https://github.com/terraform-providers/terraform-provider-runscope/issues/12))
* **New Data Source:** `runscope_buckets` ([#12](https://github.com/terraform-providers/terraform-provider-runscope/issues/12))

ENHANCEMENTS:

*  resource/runscope_bucket: Import support added ([#12](https://github.com/terraform-providers/terraform-provider-runscope/issues/12))

## 0.2.0 (July 03, 2018)

ENHANCEMENTS:

* resource/runscope_test: No longer forces a new resource when the `name` attribute changes.

## 0.1.0 (June 21, 2018)

Initial release.

