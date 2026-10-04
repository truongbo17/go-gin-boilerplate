#!/usr/bin/env bash
set -euo pipefail

if git grep -n -i -E '9pay|platform-email|service_email|gitlab\.9pay' -- ':!scripts/check-public-base.sh'; then
  echo 'Company-specific text found' >&2
  exit 1
fi
if git grep -n -E -- '-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----' -- ':!scripts/check-public-base.sh'; then
  echo 'Private key found' >&2
  exit 1
fi
for path in internal/app/core/{aws,campaign,contact,sms,tracking} internal/app/v1/{campaign,contact,sms,tracking}; do
  if test -e "$path"; then
    echo "Product-specific path remains: $path" >&2
    exit 1
  fi
done
