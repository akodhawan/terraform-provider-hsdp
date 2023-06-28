---
subcategory: "Telemetry Data Repository (TDR)"
page_title: "HSDP: hsdp_tdr_contract"
description: |-
  Manages HSDP TDR Contracts
---

# hsdp_tdr_contract

Provides a resource to manage HSDP TDR [contract](https://www.hsdp.io/documentation/telemetry-data-repository/tdr-apr23/getting-started).

## Pre-requisite

TDR Onboarding is still via [SNOW request > Onboarding on TDR](https://www.hsdp.io/documentation/telemetry-data-repository/tdr-apr23/getting-started)

## Example Usage

The following example creates a TDR Contract

```hcl
# Create TDR Contract
resource "hsdp_tdr_contract" "tdr_contract_1" {
  tdr_endpoint = "${var.tdr_base_url}/store/tdr"
  description = "TDR Contract Example"
  dataType = {
    system = "TDRXYZSystem001"
    code =  "TDRXYZCode001"
  }
  sendNotifications:false
  organization: "HSDPSolutions"
  deletePolicy: {
    duration: 7
    unit: "DAY"
	}
  schema =  {
      "$schema": "http://json-schema.org/draft-04/schema#",
      "id": "http://jsonschema.net",
      "type": "object",
      "properties": {
          "hubId": {
              "id": "http://jsonschema.net/doc/string",
              "type": "string"
          },
          "deviceModel": {
              "id": "http://jsonschema.net/doc/boolean",
              "type": "string"
          }
      },
      "required": [
          "hubId",
          "deviceModel"
      ]
  }
}

```

## Argument Reference

The following arguments are supported:

* `tdr_endpoint` - (Optional) The TDR endpoint to be used (in case override between April 2023 vs December 2021 release is required). Default April 2023 release.
* `description` - (Optional) Description of the TDR Contract
* `organization` - (Required) The TDR Orgnization Namespace which is a textual representation of the namespace the DataItem belongs to (maxLength 255).
* `dataType` - (Required) consists of `system` and `code`
  * `system`: - (Required) URN identifying the system of the value (maxLength 255). 
  * `code`: - (Required) Value of the code within the system (maxLength 255).
* `sendNotifications` - (Optional) If set to `true`, uses the HSDP Notification Service for sending notifications when POST or DELETE operations are performed on DataItems for this Contract (boolean). Default: `false`
* `deletePolicy` - (Required) This policy specifies when the DataItem needs to be deleted.
  * `description`: - (Optional) Description of the deletion policy
  * `duration`: - (Required) Integer value determining duration (maximum 365). 
  * `unit`: - (Required) Define the unit of `duration` and is a enum with allowed values `[ DAY, MONTH, YEAR ]`.
* `schema`: (Required) The JSON schema describing how the data belonging to this Contract looks.

!> Post `deletePolicy` dataItems are marked as `TombStone` and are enventually deleted after grace period of `30 days`.

## Attributes Reference

No attributes are exported.

## Import

Importing existing contracts is currently not supported