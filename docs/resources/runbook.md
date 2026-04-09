---
page_title: "NVIDIA Mission Control autonomous hardware recovery_runbook Resource - terraform-provider-mcahr"
subcategory: ""
description: |- Runbooks replace static runbooks by capturing interactive debug and remediation sessions in a convenient UI.
---

# NVIDIA Mission Control autonomous hardware recovery_runbook (Resource)

Through NVIDIA Mission Control autonomous hardware recovery's web-based UI, Runbooks automatically capture an entire debug and remediation session -- which can optionally be associated with a specific Alarm -- and then be shared with other team members to streamline future incident response. With Runbooks you can:

- Create a series of interactive Op statement cells allowing you to execute Op commands within your browser -- all without installing or configuring the local CLI.
- Define and use dynamic parameters across Runbook Op cells.
- Memorialize Runbooks with historical snapshots.
- Add Markdown-based notes to inform operators how to use the Runbook.
- Associate existing Alarms and Runbooks, allowing on-call members to click through to an interactive debugging and remediation Runbook directly from the triggered Alarm UI.

## Required Properties

Each Runbook uses a variety of properties to determine its behavior. The required properties when creating a Runbook are:

- `name`: string - A unique symbol name for the Runbook object.
- `cells`: list(object) - A list of cells represented by JSON objects. Cells may either be Op statement cells or Markdown cells.

### Download a Runbook as a Terraform resource

You can download an entire Runbook directly as a Terraform resource. This will allow you to just plug in the TF code into your infrastructure and deploy the runbook immediately.

1. Click the **Actions** button on the right side of the active Runbook panel.
2. Select the **Download Runbook as Terraform** button to download the full configuration of the current Runbook as a Terraform resource.

## Defining a Runbook using the legacy `data` property

You can also export the Runbook's configuration as a JSON file and then freely modify, share, and upload this Runbook at any time.

Note: this way of defining is **deprecated**. Please refer to the above instructions using the new format.

The following example creates a Runbook named `my_runbook`.

1. Download a Runbook as JSON.
2. Only keep the `cells`, `params`, `external_params`, and `enabled` fields fron the JSON file. Note: `externalParams` needs to be renamed to `external_params` in the JSON file.
3. Save the Runbook JSON to local path within your Terraform project.
4. Define a new `NVIDIA Mission Control autonomous hardware recovery_runbook` Terraform resource in your Terraform configuration that points the `data` property to the correct local module path.

   ```terraform
resource "mcahr_runbook" "data_runbook" {
  name        = "data_runbook"
  description = "A sample runbook defined using the data field, which loads the runbook configuration from a separate JSON file."
  data        = file("${path.module}/data.json")
}


resource "mcahr_runbook" "full_runbook" {
  cells_list = [
    {
      md = "CREATE"
    },
    {
      op          = "action success = `echo SUCCESS`"
      description = "Creates an action that echoes SUCCESS"
    },
    {
      op = "enable success"
    },
    {
      op          = "success"
      enabled     = false
      description = "Runs the success action. This cell is disabled."
    },
    {
      md = "CLEANUP"
    },
    {
      op = "delete success"
    }
  ]
  params_list = [
    {
      name  = "param_1"
      value = "<default_value>"
    },
    {
      name     = "param_2"
      value    = "<default_value>"
      required = false
      export   = true
    },
    {
      name   = "param_3"
      value  = "<default_value>"
      export = true
    },
    {
      name     = "param_4"
      required = false
    },
    {
      name        = "param_5"
      required    = false
      description = "Param #5 description"
    }
  ]
  external_params_list = [
    {
      name      = "external_param_1"
      source    = "alertmanager"
      json_path = "$.<path>"
      export    = true
      value     = "<default_value>"
    },
    {
      name      = "external_param_2"
      source    = "alertmanager"
      json_path = "$.<path>"
      value     = "<default_value>"
    },
    {
      name      = "external_param_3"
      source    = "alertmanager"
      json_path = "$.<path>"
      export    = true
    },
    {
      name      = "external_param_4"
      source    = "alertmanager"
      json_path = "$.<path>"
    },
    {
      name        = "external_param_5"
      source      = "alertmanager"
      json_path   = "$.<path>"
      description = "External parameter #5 description"
    }
  ]
  name                                  = "full_runbook"
  description                           = "A sample runbook."
  timeout_ms                            = 5000
  allowed_entities                      = ["<user_1>", "<user_2>"]
  approvers                             = ["<user_2>", "<user_3>"]
  editors                               = ["<user_2>", "<user_4>"]
  is_run_output_persisted               = true
  allowed_resources_query               = "host"
  communication_workspace               = "<workspace_name>"
  communication_channel                 = "<channel_name>"
  labels                                = ["label1", "label2"]
  communication_cud_notifications       = true
  communication_approval_notifications  = false
  communication_execution_notifications = true
  filter_resource_to_action             = true
  enabled                               = true
  secret_names                          = ["secret_1", "secret_2"]
  category                              = "general"
  params_groups = {
    "required" = []
    "optional" = ["param_2", "param_4"]
    "exported" = ["param_3"]
    "external" = ["external_param_2", "external_param_3"]
  }
}


resource "mcahr_runbook" "minimal_runbook" {
  name = "minimal_runbook"
}
```

-> See the NVIDIA Mission Control autonomous hardware recovery Runbooks Documentation for more info on creating and using Runbooks.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `allowed_entities` (List of String) The list of users who can run a runbook. Any user can run if left empty.
- `allowed_resources_query` (String) The list of resources on which a runbook can run. No restriction, if left empty.
- `approvers` (List of String) List of users who can approve runbook execution
- `category` (String) Specifies the category for this runbook. To use categories, make sure your platform administrator has enabled the `ENABLE_RUNBOOK_CATEGORIES` setting. Once enabled, you can organize your runbooks by assigning them a category.
- `cells` (String, Deprecated) The data cells inside a runbook. Defined as a list of JSON objects encoded in base64. These may be either Markdown or Op commands. Shows diffs only when configuration changes.
- `cells_list` (Attributes List) The data cells inside a runbook as a native Terraform list. Each cell is either a Markdown cell (set `md`) or an Op command cell (set `op`). Only one of `op` or `md` must be set per cell. Provides better plan changes and drift detection than the deprecated `cells` JSON string. Cannot be used together with `cells`. (see [below for nested schema](#nestedatt--cells_list))
- `communication_approval_notifications` (Boolean) Enables slack notifications for approval operations
- `communication_channel` (String) A string value denoting the slack channel where notifications related to the runbook should be sent to
- `communication_cud_notifications` (Boolean) Enables slack notifications for create/update/delete operations
- `communication_execution_notifications` (Boolean) Enables slack notifications for runbook executions
- `communication_workspace` (String) A string value denoting the slack workspace where notifications related to the runbook should be sent to
- `data` (String) JSON-encoded string containing the runbook data. This can be loaded from a file using the `file()` function, e.g., `data = file("${path.module}/runbook.json")`. Unlike other JSON fields (params, cells, external_params), this field only stores what the user sets and does not have a corresponding _full attribute.
- `description` (String) Description of the runbook
- `editors` (List of String) List of users who can edit the runbook (with configure permission). Empty maps to all users.
- `enabled` (Boolean) Whether the runbook is enabled
- `external_params` (String, Deprecated) Runbook parameters defined via JSON path used to extract the parameter's value from an external payload, encoded as JSON.
- `external_params_list` (Attributes List) Runbook parameters defined via JSON path used to extract the parameter's value from an external payload, as a native Terraform list. Provides better plan changes and drift detection than the deprecated `external_params` JSON string. Cannot be used together with `external_params`. (see [below for nested schema](#nestedatt--external_params_list))
- `filter_resource_to_action` (Boolean) Determines whether parameters containing resources are exported to actions
- `is_run_output_persisted` (Boolean) A boolean value denoting whether or not cell outputs should be persisted when running a runbook
- `labels` (List of String) A list of strings by which runbooks can be grouped
- `name` (String) The name of the runbook
- `params` (String, Deprecated) Named variables to pass to a runbook, encoded as JSON. Shows diffs only when configuration changes.
- `params_groups` (Attributes) Categorized parameter lists. Defaults to null if not specified. (see [below for nested schema](#nestedatt--params_groups))
- `params_list` (Attributes List) Named variables to pass to a runbook as a native Terraform list. Provides better plan changes and drift detection than the deprecated `params` JSON string. Cannot be used together with `params`. (see [below for nested schema](#nestedatt--params_list))
- `secret_names` (List of String) A list of strings that contains the name of the secrets that are used in the runbook.
- `timeout_ms` (Number) Maximum time to wait for runbook execution, in milliseconds

### Read-Only

- `cells_full` (String) Complete cells configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.
- `external_params_full` (String) Complete external parameter configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.
- `params_full` (String) Complete parameter configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.

<a id="nestedatt--cells_list"></a>
### Nested Schema for `cells_list`

Optional:

- `description` (String) A description for the cell.
- `enabled` (Boolean) Whether the cell is enabled.
- `md` (String) The Markdown content for this cell. Mutually exclusive with `op`.
- `name` (String) The name of the cell.
- `op` (String) The Op command content for this cell. Mutually exclusive with `md`.
- `secret_aware` (Boolean) Whether the cell is secret-aware.


<a id="nestedatt--external_params_list"></a>
### Nested Schema for `external_params_list`

Required:

- `name` (String) The name of the external parameter.
- `source` (String) The source of the external parameter.

Optional:

- `description` (String) A description for the external parameter.
- `export` (Boolean) Whether the parameter value is exported to downstream actions.
- `json_path` (String) The JSON path used to extract the parameter value from the external payload.
- `value` (String) The default value of the external parameter.


<a id="nestedatt--params_groups"></a>
### Nested Schema for `params_groups`

Optional:

- `exported` (List of String) List of exported parameter names.
- `external` (List of String) List of external parameter names.
- `optional` (List of String) List of optional parameter names.
- `required` (List of String) List of required parameter names.


<a id="nestedatt--params_list"></a>
### Nested Schema for `params_list`

Required:

- `name` (String) The name of the parameter.

Optional:

- `description` (String) A description for the parameter.
- `export` (Boolean) Whether the parameter value is exported to downstream actions.
- `required` (Boolean) Whether the parameter is required.
- `value` (String) The default value of the parameter.
