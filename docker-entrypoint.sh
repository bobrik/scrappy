#!/bin/sh -e

FLAVOR="${1}"

shift

if [ -z "${FLAVOR}" ]; then
  echo "usage:" 1>&2
  echo " slaves <flags>" 1>&2
  echo " roles  <flags>" 1>&2
  exit 1
fi

exec "/usr/bin/scrappy-${FLAVOR}" "$@"
