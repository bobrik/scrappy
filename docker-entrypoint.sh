#!/bin/sh -e

FLAVOR="${1}"

shift

if [ -z "${FLAVOR}" ]; then
  echo "usage:" 1>&2
  echo " agents <flags>" 1>&2
  echo " roles  <flags>" 1>&2
  exit 1
fi

if [ "${FLAVOR}" = "slaves" ]; then
  echo "The 'slaves' command has been replaced by the 'agents' command." 1>&2
  exit 1
fi

exec "/usr/bin/scrappy-${FLAVOR}" "$@"
