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
$bin host new -i "${HOSTIDENT}"

if ! [ -f "${HOSTIDENT}.pub" ]; then
  echo "No host public key!"
  exit 1
fi
if ! [ -f "${HOSTIDENT}.sec" ]; then
  echo "No host secret key!"
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
    --outfile "${ident}.out"

  if ! [ -f "${ident}" ]; then
    echo "Identity file ${ident} not found."
    exit 1
  fi
  if ! [ -f "${ident}.out" ]; then
    echo "Payload file ${ident}.out not found."
    exit 1
  fi
}
mktestuser "james" "James Jameson" "123 Some Street\nABC 123 Some town\nEngland"
mktestuser "hans" "Hans Hansen" "Einestraße 23\n12345 Einestadt\nDeutschland"
mktestuser "gigi" "Giacomo Gianluca" "Via Esempio 1\nLorem Citta\nItalia"
mktestuser "testy" "Testy McTestface" "Tester road\nLoch Ness\nScottland"


echo "[Step 3] Generate gift exchange"
$bin host run \
  --secret "${HOSTIDENT}.sec" \
  --outdir "${RESULTSDIR}" \
  "james.out" \
  "hans.out" \
  "gigi.out" \
  "testy.out"
if ! [ -d "${RESULTSDIR}" ]; then
  echo "Results directory ${RESULTSDIR} not found"
  exit 1
fi
for user in "James Jameson" "Hans Hansen" "Giacomo Gianluca" "Testy McTestface"; do
  if ! [ -f "${RESULTSDIR}/${user}" ]; then
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
readuser "james" "results/James Jameson"
readuser "hans" "results/Hans Hansen"
readuser "gigi" "results/Giacomo Gianluca"
readuser "testy" "results/Testy McTestface"


popd
