
# Configure Teller Client

The command line client uses environment values to configure how it communicates with the Teller agent.

The easiest method is to build a file, like in the `config` directory, containing the export commands below and then run `source <filename>` to export those to the local environment, then run the command line tool in that environment.

## Auth Key

`export AUTH_KEY=<WIF KEY>`

Specify a WIF private key to use to authenticate. The public key for this private key must be in the Teller config as an authorized client.

## Teller Key

`export TELLER_KEY=<Compressed Public Key>`

This is the public key associated with the private key that teller is using that authenticates responses from the Teller agent.

## Teller Peer Channel

`export TELLER_PEER_CHANNEL=`

This is the peer channel through which Teller agent is sent requests.

The format is `https://<peer_channel_host_domain>/api/v1/channel/<channel_id>?token=<write_token>`.

## Response Peer Channel

`export RESPONSE_PEER_CHANNEL=`

This is the peer channel through which Teller sends responses back to this client.

The format is `https://<peer_channel_host_domain>/api/v1/channel/<channel_id>?token=<read_token>`.

## Response Read Token

`export RESPONSE_READ_TOKEN=`

This is the read token for response peer channel that is used to retrieve messages posted to the peer channel.
