# Security Policy

## gengo Is Not Cryptographically Secure

**gengo is a test-data generator, not a security tool.** It is built on
[`math/rand/v2`](https://pkg.go.dev/math/rand/v2), a fast pseudo-random number
generator whose output is **predictable** and **must not** be relied upon for
any security-sensitive purpose.

Do **not** use gengo to generate:

- Passwords or passphrases
- API tokens, session IDs, or CSRF tokens
- Encryption keys, salts, nonces, or initialization vectors
- Password-reset tokens or any other secret

For all of the above, use Go's standard
[`crypto/rand`](https://pkg.go.dev/crypto/rand) package, which provides a
cryptographically secure random source.

## Supported Versions

Security fixes are applied to the latest released version only. Please upgrade to
the most recent release before reporting an issue.

| Version  | Supported          |
| -------- | ------------------ |
| 0.0.27   | :white_check_mark: |
| < 0.0.27 | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in gengo, please report it privately so
it can be addressed before public disclosure.

Use GitHub's private vulnerability reporting:

1. Open the repository's
   **[Security](https://github.com/FlavioCFOliveira/gengo/security)** tab.
2. Click **"Report a vulnerability"** to start a private security advisory.
3. Describe the issue, including steps to reproduce and the affected version(s).

Please **do not** open a public issue for security reports.

Reports will be acknowledged within a reasonable timeframe, and you will be kept
informed of the progress toward a fix and release.

> **Note:** Because gengo is intentionally non-cryptographic (see above), the
> predictability of its `math/rand/v2`-based output is **by design** and is not
> considered a vulnerability. Reports should concern genuine security defects —
> for example, memory-safety issues or behavior that could compromise a
> consuming application.
