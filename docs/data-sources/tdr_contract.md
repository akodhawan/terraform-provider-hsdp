---
subcategory: "Telemetry Data Repository (TDR)"
---

# hsdp_tdr_contract

Retrieve HSDP TDR [contract](https://www.hsdp.io/documentation/telemetry-data-repository/tdr-apr23/getting-started).


## Example Usage

The following example creates a TDR Contract

```hcl
# Retrieve TDR Contract
data "hsdp_tdr_contracts" "tdr_contracts" {
  organization_namespace = "HSDPSolutions"
}
```

```hcl
# Retrieve TDR Contract
data "hsdp_tdr_contracts" "tdr_contracts_datatypefilter" {
  tdr_endpoint = "${var.tdr_base_url}/store/tdr"
  organization_namespace = "HSDPSolutions"
  dataType = "TDRXYZSystem001|TDRXYZCode001"
}

output "contracts" {
  value = data.hsdp_tdr_contract.tdr_contracts_datatypefilter.entry
}
```

## Argument Reference

The following arguments are supported:

* `tdr_endpoint` - (Optional) The TDR endpoint to be used (in case override between April 2023 vs December 2021 release is required)
* `organization_namespace` - (Required) The TDR Orgnization Namespace which is a textual representation of the namespace the DataItem belongs to (maxLength 255).
* `dataType` - (Optional) consists of `system` and `code`. The dataType uses a token format [system]|[code] where the value of [code] matches the Coding.code and the value of [system] matches the system property of the Coding.
* `_count`: (Optional) The maximum amount of resources to return (in a single page). The default and maximum value is 100.

!> When providing multiple search parameters, please be aware that they will be used in an ‘AND’ style query. Support for ‘OR’ type queries is not present in this version of the API.

## Attributes Reference

The following attributes are exported as a bundle response:

* `resourceType` - The type of the resource - supports at this time only for resource type `Bundle`.
* `type` -  Indicates the purpose of this bundle & intended use - supports at this time only for resource type `searchset`.
* `total` - This is the number of matches for the search returned in this bundle. The total number of matching Contracts could be greater than the number in this bundle. Use the Next link to get the next set of matches.
* `_startAt` - The initial _startAt offset for this bundle.
* `link` - Links to previous and next pages.
  * `relation` - Description of the type of link - `next`
  * `url` - A Uniform Resource Identifier Reference (RFC 3986 )
* `entry` - Bundle entry array of TDR Contracts
  * `fullUrl` - Absolute URL for retrieving the resource.
  * `resource` - TDR Contract resource.
    * `id` - Identifier of the contract in token format datatype [system]|[code].
    * `description` -  Description of the TDR Contract
    * `organization` - The TDR Orgnization or Namespace which is a textual representation of the TDR organization the DataItem belongs to .
    * `dataType` -  consists of `system` and `code`
      * `description`: - Description of the TDR Data Type
      * `system`: -  URN identifying the system of the value . 
      * `code`: -  Value of the code within the system.
    * `sendNotifications` - If set to `true`, uses the HSDP Notification Service for sending notifications when POST or DELETE operations are performed on DataItems for this Contract (boolean). Default: `false`
    * `deletePolicy` -  This policy specifies when the DataItem needs to be deleted.
      * `description`: - Description of the deletion policy
      * `duration`: -  Integer value determining duration (maximum 365). 
      * `unit`: -  Define the unit of `duration` and is a enum with allowed values `[ DAY, MONTH, YEAR ]`.
    * `schema`:  The JSON schema describing how the data belonging to this Contract looks.




## Import

Importing existing contracts is currently not supported