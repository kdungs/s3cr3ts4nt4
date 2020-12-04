#!/bin/zsh

set -e

HOSTIDENT="testhost"
RESULTSDIR="results"

here="$(cd "$(dirname "${0}")" && pwd)"
bin="${here}/../s3cr3ts4nt4"

if ! [ -f "${bin}" ]; then
  echo "Binary ${bin} not found."
  exit 1
fi

dir=$(mktemp -d)
echo $dir
pushd $dir


echo "[Step 1] Generate host key"
$bin host -i "${HOSTIDENT}" new

if ! [ -f "${HOSTIDENT}.pub" ]; then
  echo "No host public key!"
  exit 1
fi
if ! [ -f "${HOSTIDENT}.id" ]; then
  echo "No host identity!"
  exit 1
fi


echo "[Step 2] Generate user payloads"
function mktestuser {
  ident="${1}"
  name="${2}"
  address="${3}"

  $bin participate \
    --hostkey "${HOSTIDENT}.pub" \
    --identity "${ident}" \
    --name "${name}" \
    --address "${address}" \

  if ! [ -f "${ident}.id" ]; then
    echo "Identity file ${ident}.id not found."
    exit 1
  fi
  if ! [ -f "${name}.in" ]; then
    echo "Payload file ${name}.in not found."
    exit 1
  fi
}
mktestuser "james" "James Jameson" '123 Some Street
ABC 123 Some town
England'
mktestuser "hans" "Hans Hansen" 'Einestraße 23
12345 Einestadt
Deutschland'
mktestuser "gigi" "Giacomo Gianluca" '
Via Esempio 1
Lorem Citta
Italia'
mktestuser "testy" "Testy McTestface" 'Tester road
Loch Ness
Scottland'


echo "[Step 3] Generate gift exchange"
$bin host --identity "${HOSTIDENT}" run \
  --outdir "${RESULTSDIR}" \
  "James Jameson.in" \
  "Hans Hansen.in" \
  "Giacomo Gianluca.in" \
  "Testy McTestface.in"
if ! [ -d "${RESULTSDIR}" ]; then
  echo "Results directory ${RESULTSDIR} not found"
  exit 1
fi
for user in "James Jameson" "Hans Hansen" "Giacomo Gianluca" "Testy McTestface"; do
  if ! [ -f "${RESULTSDIR}/${user}.out" ]; then
    echo "No recipient file for ${user} found"
    exit 1
  fi
done


echo "[Step 4] Verify that users can read their files"
function readuser {
  user="${1}"
  file="${2}"
  echo "Decrypting ${file} for ${user}"
  $bin decrypt \
    --identity "${user}" \
    "${file}"
}
readuser "james" "results/James Jameson.out"
readuser "hans" "results/Hans Hansen.out"
readuser "gigi" "results/Giacomo Gianluca.out"
readuser "testy" "results/Testy McTestface.out"


popd
