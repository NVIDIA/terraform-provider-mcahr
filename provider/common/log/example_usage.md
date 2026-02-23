# Provider Logging Configuration

This document shows how to configure logging levels for different resources using environment variables.

## Master Switches (Choose One)

At least one master switch must be set for logs to appear:

```bash
# RECOMMENDED: Provider logs
export TF_LOG_PROVIDER=INFO

# OR: Both provider and terraform core logs
export TF_LOG=INFO

# OR: Terraform core logs
export TF_LOG_CORE=INFO
```

## Environment Variable Configuration

Set logging levels using environment variables:

```bash
# Master switch - REQUIRED (choose one from above)
export TF_LOG_PROVIDER=DEBUG

# Set global default for all resources
export TF_LOG_PROVIDER_ALL=ERROR

# Override specific resources
export TF_LOG_PROVIDER_RESOURCE_ACTION=DEBUG

```

## Complete Example

```bash
# Step 1: Set master switch (REQUIRED - choose one)
export TF_LOG_PROVIDER=DEBUG  # Recommended: provider logs only

# Override specific resources
export TF_LOG_PROVIDER_RESOURCE_ACTION=DEBUG
export TF_LOG_PROVIDER_RESOURCE_BOT=DEBUG

# Now run Terraform
terraform apply
```

## Precedence Order

The logging level is determined in this order (highest to lowest priority):

1. **Master Switch**: `TF_LOG_PROVIDER`, `TF_LOG`, or `TF_LOG_CORE` (acts as a filter)
2. **Resource-specific**: `TF_LOG_PROVIDER_RESOURCE_{RESOURCE}` (e.g., `TF_LOG_PROVIDER_RESOURCE_ACTION`)
3. **Global default**: `TF_LOG_PROVIDER_RESOURCE_ALL`
4. **Fallback**: `hclog.NoLevel` (inherits from parent)

## Valid Log Levels

- `TRACE` - Most verbose
- `DEBUG` - Debug information
- `INFO` - General information
- `WARN` - Warnings
- `ERROR` - Errors only (least verbose)

## Important Notes

- **At least one master switch must be set** for logs to appear:
  - `TF_LOG_PROVIDER` (recommended - provider logs only)
  - `TF_LOG` (both provider and terraform core logs)
  - `TF_LOG_CORE` (terraform core logs only)
- **Master switch must be MORE PERMISSIVE** than resource levels (or logs will be filtered out)
  - Example: If you want DEBUG logs, set `TF_LOG_PROVIDER=DEBUG` (not INFO or ERROR)
- Environment variables use uppercase resource names (e.g., `TF_LOG_PROVIDER_RESOURCE_ACTION`)
- Invalid log levels will be ignored and fall back to the next priority level
- All configuration is done via environment variables - no provider arguments needed
