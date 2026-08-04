# Security Policy: terraform-provider-mcahr

## Reporting a Vulnerability

If you discover a potential security vulnerability in this project, please do not open a public GitHub issue.

Report it through NVIDIA Product Security:

* Preferred: [NVIDIA Product Security](https://www.nvidia.com/en-us/security/)
* Email: [psirt@nvidia.com](mailto:psirt@nvidia.com)
* PGP key for secure email: [NVIDIA public PGP key](https://www.nvidia.com/en-us/security/pgp-key)
* GitHub: the **Security** tab > **Report a vulnerability** on this repository

Please include the following where possible:

* Provider version and OpenTofu or Terraform version
* Type of vulnerability
* Step-by-step reproduction instructions
* Proof-of-concept Terraform configuration, if applicable
* Expected impact and affected resources

Detailed reports help NVIDIA evaluate and address issues more quickly.

This repository is generated from an upstream source maintained by NVIDIA. Fixes are applied upstream and republished here, so please report issues through the channels above rather than as pull requests.

## Security Architecture and Context

`terraform-provider-mcahr` is an OpenTofu and Terraform provider plugin that manages resources on the NVIDIA Mission Control autonomous hardware recovery platform through its REST API. These resources may include actions, bots, runbooks, alarms, integrations, notebooks, dashboards, file uploads, principals, and secrets.

The provider authenticates using a customer-specific bearer token and a customer-specific API endpoint URL supplied at provider configuration time. It relays configuration defined by operators into API calls to the backend. In particular, command-bearing resources can ultimately influence actions executed by the platform, which makes the integrity of Terraform configurations and imported modules an important security boundary.

This software operates at the provider plugin layer. Its primary security responsibilities are:

* protecting the platform bearer credential and the third-party integration credentials that pass through it from unintended disclosure
* correctly scoping authenticated API requests to the configured backend endpoint
* preserving operator intent when relaying configuration to the backend

Repository exposure classification: Public. Basis: publicly readable repository; this document is written for public consumption.

Service exposure classification: External / Regulated (high confidence). Basis: publicly distributed provider for external customer use; handles customer API tokens, third-party integration credentials, and secret references.

Key security boundaries:

* Trusted:
  * operator-reviewed Terraform and OpenTofu configurations
  * Terraform or OpenTofu execution environment, including its local filesystem
  * platform API server and authenticated backend services
  * operator-controlled environment variables and `.tfvars` inputs
* Untrusted:
  * imported third-party or community Terraform modules
  * public module registries and other external dependency sources
  * untrusted CI inputs or environment-variable injection sources
  * any endpoint URL or token value derived from attacker-controlled input

Primary data flows:

* inbound at provider configuration time: API endpoint URL and bearer credential, from the provider block or the `MCAHR_URL` and `MCAHR_TOKEN` environment variables
* outbound during resource CRUD operations: authenticated REST API requests
* file content: `mcahr_file` reads from a local path or a remote URL and uploads to object storage using presigned URLs issued by the backend

## Threat Model

The following scenarios represent the primary security concerns for this project.

### 1. Credential exposure through Terraform configuration, plan artifacts, or state

If a bearer token is supplied in provider configuration rather than through environment-based secret handling, it may be captured in local files, CI logs, or Terraform state storage. Even when output redaction is enabled, downstream state handling remains an important risk.

### 2. Credential exposure through provider debug logging

With debug logging enabled — via `MCAHR_DEBUG`, `TF_LOG`, or `TF_LOG_PROVIDER` — the provider records HTTP request and response bodies. Authorization headers are masked, and fields whose names identify credential material — API keys, client secrets, certificates, passwords, tokens, secret values and presigned storage URLs — are masked inside those bodies as well; a body that is not JSON is dropped rather than logged. The masking is name-based, so a credential returned under an unrecognised field name would still be recorded. Treat provider debug logs as secret material and avoid enabling them in shared CI environments.

### 3. Malicious or unsafe command content from imported modules

Some platform resources may carry command content that is relayed to the backend for execution. A malicious or compromised Terraform module could introduce unsafe command strings that an operator did not fully review before apply.

### 4. Token exfiltration through an attacker-controlled endpoint URL

If the configured provider URL is set from an untrusted source, the provider may send an otherwise valid bearer token to an attacker-controlled destination. The provider does not restrict the endpoint to a particular host or scheme.

### 5. Supply-chain compromise of published provider binaries or release automation

A compromise affecting release automation, signing material, or artifact publication could result in a trojanized provider binary being distributed to users. Because provider plugins execute in the Terraform or OpenTofu workflow, a compromised release could access sensitive runtime context.

### 6. Secret exposure through state-backed managed resources

If resource values representing secrets are persisted in Terraform state or other shared backend artifacts, users with read access to those systems may gain access to sensitive data.

Credential attributes on the integration resource — API keys, client secrets, certificates, service-account credentials and SMTP passwords — are declared `Sensitive`, so Terraform redacts them in plan and apply output and in CI job logs. They are still written to state unredacted, and the provider encrypts nothing at rest, so state remains secret-bearing regardless.

### 7. Remote file fetch initiated by the provider

An `input_file` value beginning with `https://` is fetched by the provider from the machine running Terraform, and its content is then uploaded to platform storage. In CI this can reach internal services not otherwise exposed, so an unreviewed module that sets `input_file` can both probe internal endpoints and relay their content outward. The platform bearer credential is not attached to this request.

Plaintext `http://` sources are rejected, so the fetched content cannot be chosen by a network-position attacker without also breaking TLS. There is no destination allowlist, so treat `input_file` values from modules you did not write as untrusted.

### 8. Network-path interception against the configured API endpoint

The provider depends on standard TLS protections for communication with the configured backend endpoint. If operators direct the provider to an untrusted endpoint or execute it in a compromised network environment, authenticated traffic may be exposed.

### 9. Destructive rollback after a failed create

When a create request is accepted by the backend but a later step of the same apply fails — response conversion, plan-consistency checking, or post-processing — the provider attempts to delete the resource again so the apply does not leave a partially created object behind.

That rollback establishes only that the create request succeeded, not that it brought the object into existence. If a resource of that name already existed and the platform treated the create as an update, the rollback deletes an object your configuration did not create. Operators managing resources whose names may already exist on the platform should review state carefully after a failed apply.

## Critical Security Assumptions

This project assumes the following.

### 1. Sensitive credentials are handled as secrets

Operators are expected to provide bearer credentials through secure secret-handling mechanisms such as the `MCAHR_TOKEN` environment variable or protected CI secret stores, rather than hardcoding them into version-controlled Terraform files.

### 2. Imported modules are operator-reviewed before use

The provider does not independently validate the safety or intent of command-bearing configuration supplied by external modules. Operators are responsible for reviewing imported modules before applying them.

### 3. The configured backend URL is legitimate and operator-controlled

The provider assumes the configured API endpoint belongs to the intended backend and is not attacker-controlled. It does not pin certificates or enforce that the URL uses HTTPS.

### 4. The backend performs all authentication and authorization

The provider makes no access-control decisions of its own. It forwards a bearer token and relies on the platform API server to reject any operation the caller is not entitled to perform.

### 5. Terraform state storage is access-controlled

The provider does not control where Terraform or OpenTofu stores state. Operators are responsible for securing local and remote state backends against unauthorized read access.

### 6. The execution host is trusted

The credential is read from the environment or configuration, file content transits local temporary files, and logs are written to an operator-specified path — all in plaintext, protected only by local filesystem permissions.

### 7. Backend-issued values are authentic

Presigned storage URLs returned by the API are assumed genuine and short-lived, and the version the backend reports for itself is assumed truthful. That version determines which attributes the provider sends, so an incorrect value silently changes the request payload.

### 8. Release artifacts are obtained from trusted distribution paths

Users are expected to install provider binaries from trusted release channels and to follow their normal artifact verification practices before production use.
