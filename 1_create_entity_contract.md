
# Create Entity Contract

## Get Tokenized CLI

Get the latest Tokenized specification repo and build the CLI (Command Line Interface) tool. Use the develop branch for now.

If you don't have the specification repo yet.
* Choose a directory on your computer, go into it, create a `specification` directory, and go into it.
* `git clone git@github.com:tokenized/specification.git`
* `git checkout develop`
* `make dist-cli`

If you already had the specification repo, then ensure you have the latest version:
* `git checkout develop`
* `git pull`
* `make dist-cli`

This will create an executable at `dist/tokenized`. Use this for all later references to CLI.

## Create an entity administration key using Tokenized CLI

* `tokenized gen_key`

This will output a new randomized WIF key. You can also create a WIF key using your own processes.

## Send some bitcoin to that key

Send some bitcoin, using any wallet, to the P2PKH `Address` output for the entity contract administrator key. Then get the outpoint for that UTXO (txid and output index), and the satoshi value in that UTXO.

## Create a JSON file containing a `ContractOffer`

This will specify the issuer information and jurisdiction.

Ensure the `ContractFee` field is specified by the smart contract agent host.

```
{
  "ContractType": 0,
  "GoverningLaw": "USA",
  "Jurisdiction": "USA",
  "ContractName": "Sample Business (Development - Test purposes only)",
  "Issuer": {
    "Name": "Sample Business Incorporated",
    "Type": "C",
    "CountryCode": "GBR",
    "DomainName": "sample_business.com"
  },
  "ContractFee": 2000
}
```

## Get Smart Contract Agent Address

Either host a smart contract agent and get an address from it, or get a smart contract agent address from Tokenized.

These addresses can expire, so make sure you are ready to complete this process and post the contract offer transaction before getting the address.

## Create an entity `ContractOffer` transaction

* `tokenized create_contract <admin_WIF> <outpoint> <outpoint_value> <contract_offer_json_file> <contract_address> <change_address>`

This will output information about the transaction created as well as the hex of the raw transaction.

## Post this transaction on chain

It is probably easiest to paste the hex into [WhatsOnChain Broadcast Tool](https://whatsonchain.com/broadcast).

## Ensure the request was accepted by the Smart Contract Agent

Look up the smart contract agent address in a block explorer and verify the transaction you broadcast is there and another transaction is there that spends the first output from that transaction and contains a `ContractFormation` action.

You can decode the transaction hex or just the output hex with `tokenized decode <hex>` to verify it is a `ContractFormation`.

If the response tx contains a `Rejection` action then it should specify a reason. If it failed or no response is seen then check with the smart contract agent host to determine what went wrong.
