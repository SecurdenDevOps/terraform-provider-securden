# Securden Terraform Provider

The `securden` provider allows Terraform to interact with Securden's API, enabling you to retrieve account details such as passwords and ports based on an `account_id`, or a combination of `account_name` and `account_title`.

## Changelog

### v1.0.0
- **New Feature**: Added support for bulk password retrieval using the `securden_passwords` data source.
  - Fetch passwords for multiple accounts in a single API request.
  - Reduces the time required for retrieving passwords compared to individual account fetches.
  - Allows specifying a list of `account_ids` to retrieve their corresponding passwords in bulk.

### v0.1.1
- Initial release with the following features:
  - Fetch account details using `account_id`, or a combination of `account_name` and `account_title`.
  - Supports key-value retrieval, including passwords, ports, private keys and additional fields.
