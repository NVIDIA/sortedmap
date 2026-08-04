# Security Policy: sortedmap

## Reporting a Vulnerability

If you discover a potential security vulnerability, please **do not open a public
GitHub issue, discussion, or pull request.**

- **Web (preferred):** [NVIDIA Vulnerability Disclosure Program](https://www.nvidia.com/en-us/security/)
- **E-mail:** [psirt@nvidia.com](mailto:psirt@nvidia.com)
  - For secure communication, use the [NVIDIA public PGP key](https://www.nvidia.com/en-us/security/pgp-key).
- **GitHub:** Use this repository's **Security** tab and select **Report a vulnerability**.

Please include:

- Project name (`sortedmap`) and the affected branch, commit, or module version
- Affected implementation (LLRB vs B+Tree), API, and compare function when relevant
- Vulnerability type, reproduction steps, proof-of-concept if available, and impact assessment

NVIDIA's Product Security Incident Response Team (PSIRT) will acknowledge the report,
validate severity, coordinate remediation, and publish a security bulletin when
appropriate. See [PSIRT policies](https://www.nvidia.com/en-us/security/psirt-policies/).

## Supported Versions

`sortedmap` is developed on the `development` branch. Security fixes land there
unless a release branch is explicitly announced.

| Version or branch | Supported |
| --- | --- |
| `development` | Yes |
| Older tags / branches | No, unless explicitly stated |

## Security Architecture & Context

`sortedmap` is a Go library providing sorted maps as either an in-memory
left-leaning red-black tree or a pageable B+Tree, with pluggable key comparison
and optional dump/validate helpers. It has no network listeners or built-in
credential handling.

This software operates at the **library** level. Its primary security
responsibilities are preserving map ordering/integrity invariants for the chosen
tree implementation and safely interacting with any caller-supplied B+Tree page
storage callbacks.

**Repository Exposure Classification:** Public.
Basis: origin remote is the publicly accessible `NVIDIA/sortedmap` GitHub repository.

**Service Exposure Classification:** External / Regulated (high confidence).
Basis: externally distributed open-source Go library under the NVIDIA GitHub
organization.

Key security boundaries:

- Keys, values, and `Compare` functions are caller-controlled.
- B+Tree paging (when used) crosses into caller-provided page read/write
  backends; those backends define the persistence and trust model.
- `Dump` / `Validate` are diagnostic aids, not access-control mechanisms.

### Threat Model

1. **Ordering or lookup integrity failure:** Bugs in LLRB/B+Tree algorithms or
   compare handling return wrong values, miss keys, or break bisect indexes so
   callers make incorrect authorization or data decisions.
2. **Hostile or inconsistent compare functions:** Nondeterministic or
   non-transitive `Compare` implementations corrupt tree structure or enable
   logic bugs in consumers.
3. **B+Tree page-backend integrity:** A malicious or buggy page store returns
   inconsistent pages, causing data corruption, panics, or disclosure of stale
   page contents across logical keys.
4. **Resource exhaustion:** Pathological insert/delete patterns or unbounded
   caches cause memory growth or CPU spin in long-lived processes.

### Critical Security Assumptions

- Callers supply transitive, deterministic compare functions appropriate to
  their key types.
- Page backends used with the B+Tree are trusted for confidentiality and
  integrity of stored pages, or the application adds its own protection.
- `sortedmap` does not authenticate callers or isolate tenants.
- The Go runtime is a trusted substrate.

## Out of Scope

- Application misuse of sorted maps as an access-control list without additional
  checks.
- Performance variance or intentional stress tests in isolated harnesses.
- Issues solely in dependent applications that embed `sortedmap`.
