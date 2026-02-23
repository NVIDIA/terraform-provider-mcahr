---
page_title: 'NVIDIA Mission Control autonomous hardware recovery_action Resource - terraform-provider-mcahr'
subcategory: ''
description: |-
---

# NVIDIA Mission Control autonomous hardware recovery_action (Resource)

Actions execute shell commands on associated Resources. Whenever an Alarm fires, the associated Bot triggers the corresponding Action, closing the basic auto-remediation loop of
NVIDIA Mission Control autonomous hardware recovery.

## Required Properties

Each Action has many properties that determine its behavior. The required properties are:

- name - The name of the Action.
- command - The shell command executed when the Action triggers.

-> Check out Action Properties for details on all available properties and how to use them.

## Usage

The following Action definition creates a `cpu_threshold_action` that compares host CPU usage against a `cpu_threshold` parameter value.

```terraform
resource "mcahr_action" "cpu_threshold_action" {
  name = "cpu_threshold_action"
  # Evaluates current CPU usage and compares it to a parameter value named $cpu_threshold
  command = "`if [ $[100-$(vmstat 1 2|tail -1|awk '{print $15}')] -gt $cpu_threshold ]; then exit 1; fi`"

  description    = "Check CPU usage"
  enabled        = true
  params         = ["cpu_threshold"]
  resource_query = "hosts"
  timeout        = 5000

  start_title_template    = "CPU threshold action started"
  complete_title_template = "CPU threshold action completed"
  error_title_template    = "CPU threshold action failed"

  start_short_template    = "CPU threshold action short started"
  complete_short_template = "CPU threshold action short completed"
  error_short_template    = "CPU threshold action short failed"
}

resource "mcahr_action" "cpu_threshold_action_escaped_command" {
  name = "cpu_threshold_action_escaped_command"
  # Command should be escaped quotes
  command     = "`echo \"test command escaped quotes\"`"
  description = "Check CPU usage with escaped command"
}
```

This Action can be executed via an Alarm's clear_query
/ fire_query, or directly via an Op command.

For example, the following Alarm fires and clears based on the result of the previously-generated `cpu_threshold_action`:

```terraform
resource "mcahr_alarm" "cpu_threshold_alarm" {
  fire_query = "cpu_threshold_action(cpu_threshold=75) == 1"
  name       = "cpu_threshold_alarm"

  clear_query        = "cpu_threshold_action(cpu_threshold=75) == 0"
  description        = "High CPU usage alarm"
  enabled            = true
  resource_query     = "hosts"
  check_interval_sec = 10

  fire_short_template    = "High CPU Alarm fired"
  resolve_short_template = "High CPU Alarm resolved"
}
```

You can also define Terraform Input Variables and use them within your Action definitions:

```terraform
variable "namespace" {
  type        = string
  description = "A namespace to isolate multiple instances of the module with different parameters."
}

variable "resource_query" {
  type        = string
  description = "The set of hosts/pods/containers monitored and affected by this module."
}

variable "jvm_process_regex" {
  type        = string
  description = "A regular expression to match and select the monitored Java processes."
}

variable "mem_threshold" {
  type        = number
  description = "The high-water-mark, in Mb, above which the JVM process stack-trace is dumped."
  default     = 2000
}

variable "check_interval" {
  type        = number
  description = "Frequency, in seconds, to check the memory usage."
  default     = 60
}

variable "script_path" {
  type        = string
  description = "Destination (on selected resources) for the check, and stack-dump scripts."
  default     = "/agent/scripts"
}

variable "s3_bucket" {
  type        = string
  description = "Destination in AWS S3 for stack-dump output files."
  default     = "shore-oppack-test"
}
```

```terraform
# Action to check the JVM heap usage on the selected resources and process.
resource "mcahr_action" "jvm_trace_check_heap" {
  name        = "${var.namespace}_jvm_check_heap"
  description = "Check heap utilization by process regex."
  # Parameters passed in: the regular expression to select process name.
  params = ["JVM_PROCESS_REGEX"]
  # Extract the heap used for the matching process and return 1 if above threshold.
  command = "`hm=$(jstat -gc $(jps | grep \"$${JVM_PROCESS_REGEX}\" | awk '{print $1}') | tail -n 1 | awk '{split($0,a,\" \"); sum=a[3]+a[4]+a[6]+a[8]; print sum/1024}'); hm=$${hm%.*}; if [ $hm -gt ${var.mem_threshold} ]; then echo \"heap memory $hm MB > threshold ${var.mem_threshold} MB\"; exit 1; fi`"

  # UI / CLI annotation informational messages:
  start_short_template    = "Checking JVM heap usage."
  error_short_template    = "Error checking JVM heap usage."
  complete_short_template = "Finished checking JVM heap usage."
  start_long_template     = "Checking JVM process ${var.jvm_process_regex} heap usage."
  error_long_template     = "Error checking JVM process ${var.jvm_process_regex} heap usage."
  complete_long_template  = "Finished checking JVM process ${var.jvm_process_regex} heap usage."

  enabled = true
}
```

-> See the NVIDIA Mission Control autonomous hardware recovery Actions Documentation for more info.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `command` (String) The command to execute for this action
- `name` (String) The name of the action

### Optional

- `allowed_entities` (List of String) Allowed entities
- `allowed_resources_query` (String) Query for allowed resources
- `communication_channel` (String) Communication channel
- `communication_workspace` (String) Communication workspace
- `complete_long_template` (String, Deprecated) **Deprecated** Long template for completion notifications (deprecated - server controlled)
- `complete_short_template` (String) Short template for completion notifications
- `complete_title_template` (String) Title template for completion notifications
- `description` (String) Description of the action
- `editors` (List of String) Editors of the action
- `enabled` (Boolean) Whether the action is enabled
- `error_long_template` (String, Deprecated) **Deprecated** Long template for error notifications (deprecated - server controlled)
- `error_short_template` (String) Short template for error notifications
- `error_title_template` (String) Title template for error notifications
- `file_deps` (List of String) File dependencies
- `params` (List of String) Action parameters
- `res_env_var` (String) Resource environment variable
- `resource_query` (String) Query to identify resources
- `resource_tags_to_export` (List of String) Resource tags to export
- `shell` (String) Shell to use for command execution
- `start_long_template` (String, Deprecated) **Deprecated** Long template for start notifications (deprecated - server controlled)
- `start_short_template` (String) Short template for start notifications
- `start_title_template` (String) Title template for start notifications
- `timeout` (Number) Timeout for action execution
