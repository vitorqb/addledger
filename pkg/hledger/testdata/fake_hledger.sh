#!/bin/bash
# A fake hledger executable used during tests.

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

echo "ERROR: UNEXPECTED COMMAND" >&2
exit 1
