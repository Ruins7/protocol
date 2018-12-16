#!/bin/bash

#
# Test creating a single ApplyValidator transaction in a 1-node chain, reset each time
#
CMD=$GOPATH/src/github.com/Oneledger/protocol/node/scripts
TEST=$GOPATH/src/github.com/Oneledger/protocol/node/tests


$TEST/register.sh
$TEST/testmint.sh

echo "================== Test dynamic validator ==================="

olclient testmint -c Emma --party Emma --amount 10 --currency VT
# Let the money get processed
sleep 3

echo "Emma sending validator token with 5 coins"
olclient applyvalidator -c Emma --id Emma --amount 5

sleep 4

olclient testmint -c David --party David --amount 10 --currency VT

sleep 4

olclient send -c David --party David --counterparty Bob --amount 6 --currency VT

echo "============================================================="
olclient list -c Emma
olclient list -c Bob

num=`pgrep -f "^olfullnode node -c *.*" | wc -l `

if [[ $num < 5 ]]; then
    echo "Validator test failed"
    exit  -1
fi

echo "Validator test succeed"

$CMD/stopOneLedger
