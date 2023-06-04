#!/bin/bash
# A fake hledger executable used during tests.

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# accounts are all accounts returned in success scenarios.
function accounts() {
    cat <<EOF
assets:bank:current:bnext
assets:bank:savings:itau
assets:cash
assets:other
expenses:bank-fees
expenses:trips-and-travels
expenses:unknown
expenses:urban-transportation:public
expenses:urban-transportation:taxi-uber-others
initial-balance
liabilities:credit-cards:amex
liabilities:other
revenues:earned-interests
revenues:salary
EOF
}

function transactions() {
    cat ${SCRIPT_DIR}/transactions.json
}

# Case 1 - `accounts` w/out file
if [[ "$1" == "accounts" ]] && [[ "$#" == "1" ]]
then
    accounts
    exit 0
fi

# Case 2 - `accounts` w/ file
if [[ "$1" == "accounts" ]] && [[ "$2" == "--file=foo" ]] && [[ "$#" == "2" ]]
then
    accounts
    exit 0
fi

# Case 3 - `print` w/ file
if [[ "$1" == "--file=foo" ]] && [[ "$2" == "print" ]] && [[ "$3" == "--output-format=json" ]] && [[ "$#" == "3" ]]
then
    transactions
    exit 0
fi

echo "ERROR: UNEXPECTED COMMAND" >&2
exit 1
