#!/usr/bin/env bash
set -euo pipefail

if git grep -n -E -- '-----BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY-----' -- ':!scripts/check-source.sh'; then
  echo 'Private key found' >&2
  exit 1
fi
