---
page_title: 'NVIDIA Mission Control autonomous hardware recovery_report_template Resource - terraform-provider-mcahr'
subcategory: ''
description: |-
---

# NVIDIA Mission Control autonomous hardware recovery_report_template (Resource)


## Properties


- <b>name</b> (String) The name of the Report Template.
- <b>links</b> (List of Object) References to other related Report Templates. A list of objects with the following attributes:
    - <b>label</b> (String) A label for the link.
    - <b>report_template_name</b> (String) The name of the linked Report Template.
- <b>blocks</b> (String) A list of Report Template blocks in JSON format. Typically, this string should be set using the jsonencode function with a Terraform Input Variable. All it's properties must be present to avoid Terraform diffs. It has the following properties:
    - <i><b>title</b></i> (String) The name of the report block.
    - <i><b>resource_query</b></i> (String) Specifies which resources to include in the chart.
    - <i><b>group_by_tag</b></i> (String) The resource tag used to group resources in the chart.
    - <i><b>breakdown_by_tag</b></i> (String) The tag within each group used to further break down resources.
    - <i><b>breakdown_tags_values</b></i> (List of Object) Specifies which values of the breakdown tag to display in the chart. It's a list of objects with the following attributes:
        - <i><b>color</b></i> (String) The hexadecimal color code (`#RRGGBB`).
        - <i><b>values</b></i> (List of String) Tag values.
        - <i><b>label</b></i> (String) A label.
    - <i><b>include_other_breakdown_tag_values</b></i> (Boolean) When set to `true`, resources that do not have a value set for the breakdown tag are included in a separate `other` section of the specific row.
    - <i><b>view_mode</b></i> (String) Determines the display format for the bar charts, either as a `COUNT` (numerical count) or `PERCENTAGE` (percentage of the whole).
    - <i><b>resources_breakdown</b></i> (List of Object) Contains all data necessary for building the chart. It's a list of objects with the following attributes:
        - <i><b>group_by_value</b></i> (String) Existing tag value or `__no_value__`.
        - <i><b>breakdown_values</b></i> (List of Object) A list of objects with the following attributes:
            - <i><b>value</b></i> (String) Existing tag value or `__no_value__`.
            - <i><b>count</b></i> (Number)
    - <i><b>other_tags_to_export</b></i> (List of String) Additional tags (besides the group and breakdown tags) to include when exporting the Report Template.
    - <i><b>include_resources_without_group_tag</b></i> (Boolean) When set to `true`, resources without a group tag value are included in the chart in an another row labeled `other`.
    - <i><b>group_by_tag_order</b></i> (Object) Defines the display order for the values of the group by tag in the chart. Has the following attributes:
        - <i><b>type</b></i> (String) Can be one of the following: `DEFAULT`, `BY_TOTAL_ASC`, `BY_TOTAL_DESC`, `CUSTOM`.
        - <i><b>values</b></i> (List of String) If <b>type</b> is `CUSTOM`, this list defines the order of tags.



## Example


```terraform
resource "mcahr_report_template" "full_report_template" {
  name = "full_report_template"
  blocks_list = [
    {
      title            = "Block Name"
      resource_query   = "host"
      group_by_tag     = "tag_0"
      breakdown_by_tag = "tag_1"
      breakdown_tags_values = [
        {
          color  = "#AAAAAA"
          values = ["passed", "skipped"]
          label  = "label_0"
        }
      ]
      view_mode                           = "PERCENTAGE"
      include_other_breakdown_tag_values  = true
      other_tags_to_export                = ["other_tag_1", "other_tag_2"]
      include_resources_without_group_tag = false
      group_by_tag_order = {
        type   = "DEFAULT"
        values = []
      }
      resources_breakdown = [
        {
          group_by_value = "tag_0"
          breakdown_values = [
            {
              value = "value"
              count = 1
            }
          ]
        }
      ]
    }
  ]
  depends_on = [
    mcahr_report_template.minimal_report_template
  ]
  links_list = [
    {
      label                = "minimal-report"
      report_template_name = "minimal_report_template"
    }
  ]
}


resource "mcahr_report_template" "minimal_report_template" {
  name = "minimal_report_template"
  blocks_list = [
    {
      title            = "Block Name"
      resource_query   = "host"
      group_by_tag     = "tag_0"
      breakdown_by_tag = "tag_1"
      breakdown_tags_values = [
        {
          color  = "#AAAAAA"
          values = ["passed", "skipped"]
          label  = "label_0"
        }
      ]
      view_mode                           = "COUNT"
      include_other_breakdown_tag_values  = true
      other_tags_to_export                = ["other_tag_1", "other_tag_2"]
      include_resources_without_group_tag = false
      group_by_tag_order = {
        type   = "DEFAULT"
        values = []
      }
      resources_breakdown = [
        {
          group_by_value = "tag_0"
          breakdown_values = [
            {
              value = "value"
              count = 1
            }
          ]
        }
      ]
    }
  ]
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the report template

### Optional

- `blocks` (String, Deprecated) The JSON encoded blocks of the report template.
- `blocks_list` (Attributes List) The blocks of the report template as a native Terraform list. Provides better plan changes and drift detection than the deprecated `blocks` JSON string. Cannot be used together with `blocks`. (see [below for nested schema](#nestedatt--blocks_list))
- `links` (String, Deprecated) The JSON encoded links of a report template with other report templates.
- `links_list` (Attributes List) The links of a report template with other report templates as a native Terraform list. Provides better plan changes and drift detection than the deprecated `links` JSON string. Cannot be used together with `links`. (see [below for nested schema](#nestedatt--links_list))

### Read-Only

- `blocks_full` (String) Complete blocks configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.
- `links_full` (String) Complete links configuration returned by the API, including server-added fields. Shows diffs when external drift is detected and when configuration changes.

<a id="nestedatt--blocks_list"></a>
### Nested Schema for `blocks_list`

Required:

- `breakdown_by_tag` (String) The tag to break down resources by.
- `breakdown_tags_values` (Attributes List) Breakdown tag value configurations. (see [below for nested schema](#nestedatt--blocks_list--breakdown_tags_values))
- `group_by_tag` (String) The tag to group resources by.
- `resource_query` (String) The resource query for the block.
- `title` (String) The title of the block.

Optional:

- `group_by_tag_order` (Attributes) The ordering configuration for group-by tags. (see [below for nested schema](#nestedatt--blocks_list--group_by_tag_order))
- `include_other_breakdown_tag_values` (Boolean) Whether to include other breakdown tag values.
- `include_resources_without_group_tag` (Boolean) Whether to include resources without the group tag.
- `other_tags_to_export` (List of String) Additional tags to export.
- `resources_breakdown` (Attributes List) Resources breakdown configurations. (see [below for nested schema](#nestedatt--blocks_list--resources_breakdown))
- `view_mode` (String) The view mode (COUNT or PERCENTAGE).

<a id="nestedatt--blocks_list--breakdown_tags_values"></a>
### Nested Schema for `blocks_list.breakdown_tags_values`

Required:

- `color` (String) The color of the breakdown value.
- `values` (List of String) The values in the breakdown group.

Optional:

- `label` (String) The label for the breakdown value.


<a id="nestedatt--blocks_list--group_by_tag_order"></a>
### Nested Schema for `blocks_list.group_by_tag_order`

Optional:

- `type` (String) The ordering type (DEFAULT, BY_TOTAL_ASC, BY_TOTAL_DESC, CUSTOM).
- `values` (List of String) Custom ordering values.


<a id="nestedatt--blocks_list--resources_breakdown"></a>
### Nested Schema for `blocks_list.resources_breakdown`

Required:

- `breakdown_values` (Attributes List) The breakdown values. (see [below for nested schema](#nestedatt--blocks_list--resources_breakdown--breakdown_values))
- `group_by_value` (String) The group-by value.

<a id="nestedatt--blocks_list--resources_breakdown--breakdown_values"></a>
### Nested Schema for `blocks_list.resources_breakdown.breakdown_values`

Required:

- `count` (Number) The count for the breakdown value.
- `value` (String) The breakdown value.




<a id="nestedatt--links_list"></a>
### Nested Schema for `links_list`

Required:

- `label` (String) The label for the link.
- `report_template_name` (String) The name of the linked report template.
