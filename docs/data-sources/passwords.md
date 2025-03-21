---
page_title: "securden_passwords Data Source - terraform-provider-securden"
subcategory: ""
description: |-
  Securden Accounts Passwords
---

# securden_passwords (Data Source)

Multiple Accounts Passwords


## Schema

### Required

- `account_ids` (List of String) List of account IDs needs to be fetched

### Read-Only

- `passwords` (Map of String) Multiple accounts passwords with account ID as key and value will be account password
