#!/bin/bash
#
# RPC API Testing
#
# Checks transaction receipts generated from eth_getBlockReceipts
# by comparing the output with a reference API. The reference API
# endpoint is queried for receipts one transaction at a time.
#
# Usage: bash test_getBlockReceipts.sh <tested URL> <referential URL> [block_range] [end_block_id]
#        default block_range = 1000 blocks
#        default end_block_id = head block of the tested API
#
TESTED_API="$1"
REF_API="$2"
RANGE="$3"
HEAD="$4"

if [ -z "${RANGE}" ]; then
  RANGE=1000
fi

if [ -z "${HEAD}" ]; then
    HEAD=$(curl -s -X POST -H "Content-Type: application/json" --data '{"method":"eth_blockNumber","id":1,"jsonrpc":"2.0"}' ${TESTED_API} | jq -r ".result")
    HEAD=$((16#${HEAD#"0x"}))
    echo "Using end block #${HEAD}"
fi

while [ "${RANGE}" -ge 0 ]; do
    BLK=$((${HEAD}-${RANGE}))
    BLK16=$(printf '0x%x' ${BLK})

    BLK_HEAD=$(curl -s -X POST -H "Content-Type: application/json" --data '{"method":"eth_getBlockByNumber","params":["'${BLK16}'",false],"id":1,"jsonrpc":"2.0"}' ${REF_API})
    ROOT_HASH=$(echo "${BLK_HEAD}" | jq -r ".result.stateRoot")
    BLK_TXCOUNT=$(echo "${BLK_HEAD}" | jq -r ".result.transactions | length")

    echo "Testing block #${BLK}"

    MY_HEAD=$(curl -s -X POST -H "Content-Type: application/json" --data '{"method":"eth_getHeaderByNumber","params":["'${BLK16}'"],"id":1,"jsonrpc":"2.0"}' ${TESTED_API})
    MY_ROOT=$(echo "${MY_HEAD}" | jq -r ".result.stateRoot")

    if [ "${ROOT_HASH}" != "${MY_ROOT}" ]; then
        echo "Error: block hash not matched; ${ROOT_HASH} expected; ${MY_ROOT} received)"
        exit 1
    else
        echo "  block hash confirmed; ${ROOT_HASH}"
    fi

    BLK_RECEIPTS=$(curl -s -X POST -H "Content-Type: application/json" --data '{"method":"eth_getBlockReceipts","params":["'${BLK16}'"],"id":1,"jsonrpc":"2.0"}' ${TESTED_API})
    RCPT_COUNT=$(echo "${BLK_RECEIPTS}" | jq -r ".result | length")

    if [ ${BLK_TXCOUNT} -ne ${RCPT_COUNT} ]; then
        echo "Error: wrong number of receipts received; ${BLK_TXCOUNT} expected; ${RCPT_COUNT} received)"
        exit 1
    else
        echo "  transaction count confirmed; ${BLK_TXCOUNT} transaction(s) inside"
    fi

    INDEX=0
    while [ ${INDEX} -lt ${RCPT_COUNT} ]; do
        MY_RECEIPT=$(echo ${BLK_RECEIPTS} | jq --sort-keys ".result[${INDEX}]")
        TX_HASH=$(echo ${MY_RECEIPT} | jq -r ".transactionHash")

        TX_RECEIPT=$(curl -s -X POST -H "Content-Type: application/json" --data '{"method":"eth_getTransactionReceipt","params":["'${TX_HASH}'"],"id":1,"jsonrpc":"2.0"}' ${REF_API})
        TX_RECEIPT=$(echo ${TX_RECEIPT} | jq --sort-keys ".result")

        TEST=$(jq -r -n --argjson A "$MY_RECEIPT" --argjson B "$TX_RECEIPT" -f <(cat<<"EOF"
def walk(f):
  . as $in
  | if type == "object" then
      reduce keys[] as $key
        ( {}; . + { ($key):  ($in[$key] | walk(f)) } ) | f
  elif type == "array" then map( walk(f) ) | f
  else f
  end;

def normalize: walk(if type == "array" then sort else . end);

def equiv(x): normalize == (x | normalize);

if $A | equiv($B) then empty else "failed" end

EOF
)
        )
        if [ -z "${TEST}" ]; then
            echo "     #${INDEX}: ${TX_HASH}; receipt ok"
        else
            echo "Error: tx receipt check failed for #${INDEX}: ${TX_HASH}"

            echo "My Receipt"
            echo ${MY_RECEIPT} | jq

            echo "Reference Receipt"
            echo ${TX_RECEIPT} | jq

            exit 1
        fi

        INDEX=$((${INDEX}+1))
    done

    RANGE=$((${RANGE}-1))
done

exit 0
