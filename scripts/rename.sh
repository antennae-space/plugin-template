#!/bin/sh
# Makes the template your plugin: ./scripts/rename.sh <id> "<Name>" <owner>/<repo>
#   id          lowercase letters, digits and dashes: the prefix of its methods, events,
#               commands and assistant tools (dashes become underscores in the tools)
#   Name        shown in the IDE
#   owner/repo  its GitHub repository: the address of its archives and its Go module
# Run it once, on a fresh copy; then build (npm run build) and test (go test ./...).
set -eu
cd "$(dirname "$0")/.."
if [ $# -ne 3 ]; then
  sed -n '2,7p' "$0" | sed 's/^# \{0,1\}//' >&2
  exit 2
fi
id=$1 name=$2 repo=$3
echo "$id" | grep -Eq '^[a-z][a-z0-9-]*$' || { echo "the id must be lowercase letters, digits and dashes: $id" >&2; exit 1; }
echo "$repo" | grep -Eq '^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$' || { echo "the repository is owner/repo: $repo" >&2; exit 1; }
case "$name" in *[\"\\/\&#]*) echo "the name cannot hold \" \\ / & #: $name" >&2; exit 1 ;; esac
grep -q '"id": "example"' plugin.json || { echo "already renamed" >&2; exit 1; }
tool=$(echo "$id" | tr - _)
code="plugin.json package.json go.mod backend/*.go web/src/*.ts web/src/*.tsx web/src/*.css i18n/*.json .github/workflows/*.yml"
# The repository and the module first, then the tools (example_summary), the id as a word
# (example.*, example-…, "example"), the name, and its translations in the catalogs.
sed -i \
  -e "s#antennae-space/plugin-template#$repo#g" \
  -e "s#\bexample_#${tool}_#g" \
  -e "s#\bexample\b#$id#g" \
  -e "s#\bExample\b#$name#g" \
  $code
sed -i -e "s#\"plugin-example\"#\"plugin-$id\"#" package-lock.json
sed -i -e "s#\bExemple\b#$name#g" i18n/fr.json
sed -i -e "s#\bEjemplo\b#$name#g" i18n/es.json
# The README: its title and the repository; the rest is yours to write.
sed -i -e "1s#.*#\# $name#" -e "s#antennae-space/plugin-template#$repo#g" -e "s#web/src/example.css#web/src/$id.css#" README.md
git mv web/src/example.css "web/src/$id.css" 2>/dev/null || mv web/src/example.css "web/src/$id.css"
sed -i "s#'./example.css'#'./$id.css'#" web/src/index.tsx
echo "renamed to $id ($name), $repo: check with git diff, then npm run build && go test ./..."
