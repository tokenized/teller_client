
# Configure Teller

## Get Initial Config File

We will give you the initial Teller configuration file with many values already specified.

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

## Create Entity Contract

An Entity Contract is needed for Teller to use as the Issuer for its Instrument Contracts.

Using other provided instructions create an Entity Contract on chain using your issuer details and get the smart contract agent address to use in this configuration.

Use the Entity Contract smart contract agent address as the `ENTITY_CONTRACT` value in the config file.

## Generate an Extended Key

A key is needed to derive keys for Teller to use for receiving bitcoin and administering contracts used to create instruments.

* `tokenized gen_extended_key`

Use the BIP32 extended key value as the `TELLER_KEY` value in the config file.

## SpyNode

Spynode is required as an interface to the blockchain. Most of the configuration will be supplied, but you will need to generate your own authentication key and set the latest block height before connecting.

### Generate SpyNode Key

A key is needed to authenticate connections with SpyNode.

* `tokenized gen_key`

Use this WIF key as the `client_key` value under the `SpyNode` section.

Also give this key to Tokenized to add to their authorized list of SpyNode clients.

### Update SpyNode Start Block Height

Set the `start_block_height` value under `SpyNode` to the block height of the block containing the Entity Contract formation transaction.

## Contract Operator

Contract Operator hosts your Teller's smart contract agents. We will give you the URL to access it, but you need to generate an authentication key.

### Generate Contract Operator Key

A key is needed to authenticate connections with Contract Operator.

* `tokenized gen_key`

Use this WIF key as the `client_key` value under the `Operator` section.

Also give this key to Tokenized to add to their authorized list of Operator clients.

## Peer Channels Key

Set the `channels_key` value under `Teller` to a randomly generate WIF key.

## Update AWS Region

Set the value of `AWS_REGION` under `AWS` to the appropriate region.

## Update Test

Set `IS_TEST` under `Bitcoin` to true in any environment except production.

## Other Values

`FEE_RATE` is the transaction fee rate in satoshis per byte used when building transactions.

`DUST_FEE_RATE` is the transaction fee rate used to calculate dust output amounts.
    0.0 means to use 1 satoshi for dust outputs. Otherwise the calculation below is used.
```
    dust := (output_serialize_size + 148) * 3) * dust_fee_rate
```

`FUNDING_AMOUNT` is the satoshies to provide to admin addresses when they initially need bitcoin or when they run out of bitcoin. It is best to use an amount that will last a while so there are not constant extra transactions to fund the admin addresses. Received bitcoin normally goes to "bitcoin addresses" then is only transfered to "admin addresses" as needed. All addresses are derived from the `TELLER_KEY`, except external addresses like smart contract agent addresses, of course.

## Sample Config

```
{
  "AWS": {
    "AWS_REGION": "<DEPENDENT ON ENVIRONMENT>"
  },
  "Bitcoin": {
    "BITCOIN_CHAIN": "mainnet",
    "IS_TEST": true,
    "FEE_RATE": 0.05,
    "DUST_FEE_RATE": 0.0
  },
  "Db": {
    "DB_DRIVER": "postgres",
    "DB_URL": "secretsmanager://<DEPENDENT ON ENVIRONMENT>"
  },
  "Storage": {
    "STORAGE_BUCKET": "<DEPENDENT ON ENVIRONMENT>",
    "STORAGE_ROOT": "/"
  },
  "SpyNode": {
    "server_address": "<DEPENDENT ON ENVIRONMENT>",
    "server_key": "<DEPENDENT ON ENVIRONMENT>",
    "client_key": "<GENERATED FOR ENVIRONMENT>",
    "start_block_height": <DEPENDENT ON ORIGINAL SETUP>,
    "connection_type": "full",
  },
  "Operator": {
    "url": "<DEPENDENT ON ENVIRONMENT>",
    "key": "<DEPENDENT ON ENVIRONMENT>",
    "client_key": "<GENERATED FOR ENVIRONMENT>"
  },
  "Teller": {
    "TELLER_KEY": "<GENERATED FOR ENVIRONMENT>",
    "FUNDING_AMOUNT": 25000,
    "ENTITY_CONTRACT": "<GENERATED FOR ENVIRONMENT>",
    "channels_key": "<GENERATED FOR ENVIRONMENT>",
    "peer_channel": "<DEPENDENT ON ENVIRONMENT>",
    "request_authorizations": [
      "<DEPENDENT ON ENVIRONMENT>"
    ]
  }
}
```
