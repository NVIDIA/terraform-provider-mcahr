# Terraform Provider NVIDIA Mission Control autonomous hardware recovery (terraform-provider-mcahr)

The Terraform NVIDIA Mission Control autonomous hardware recovery provider allows you to manage NVIDIA Mission Control autonomous hardware recovery resources using [OpenTofu](https://opentofu.org/).

## Requirements

*   [OpenTofu](https://opentofu.org/docs/cli/install/) version 1.12.0 or later.
*   [Go](https://golang.org/doc/install) version 1.19 or later (to build the provider plugin).

## Example Usage

For detailed examples, please refer to the content in the [`examples/`](./examples/) directory and the documentation for each resource.

## Resources and Data Sources

Full documentation for available resources and data sources can be found in the [`docs/`](./docs/) directory of this repository.

## Security considerations

* **Token delivery:** prefer the `MCAHR_TOKEN` environment variable (read during `Configure`, bypasses state and plan output). Never hardcoded in a provider block or a .tfvars file.
* **Relayed commands:** `mcahr_action`, `mcahr_bot`, and `mcahr_runbook` send their `command` to the NVIDIA Mission Control autonomous hardware recovery backend for execution on platform-managed agents — vet and pin any imported modules that produce them.

## License

This provider is distributed under the Apache 2.0 License. See the [LICENSE](./LICENSE) file for more information.
