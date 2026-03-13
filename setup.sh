#!/bin/bash

set -euo pipefail

repo=$(basename "$(git rev-parse --show-toplevel)")
module="github.com/liamcervante/${repo}"

sed -i '' "s|github.com/liamcervante/go-template|${module}|g" go.mod

read -rp "Is this an executable or a library? [exe/lib] " project_type

if [[ "${project_type}" == "lib" ]]; then
  rm -f main.go
  sed -i '' 's|go build -o bin/go-template main.go|go build ./...|' Makefile
fi

cat > .git/hooks/pre-commit << 'HOOK'
#!/bin/bash
make all || exit 1

if ! git diff --quiet; then
  echo "pre-commit: files were modified by make all, please stage the changes and commit again"
  exit 1
fi
HOOK
chmod +x .git/hooks/pre-commit

rm -- "$0"