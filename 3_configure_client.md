
# Configure Teller Client

The command line client uses environment values to configure how it communicates with the Teller agent.

The easiest method is to build a file, like in the `config` directory, containing the export commands below and then run `source <filename>` to export those to the local environment, then run the command line tool in that environment.

## Peer Channels

Peer channels provide a method of posting messages back and forth between the Teller client and the Teller service. Peer channels by default have a read token and a write token. The read token allows reading of messages, marking them as read, and deleting them. The write token allows posting new messages to the channel. To post request messages to the Teller service's peer channel it needs that channel URL and the write token. To receive responses from the Teller service, the Teller client needs to provide a peer channel URL with the write token so that the Teller service can post a response message to it, then it needs that same peer channel URL with the read token so that it can read that response message.

## Configuration Values

### Auth Key

`export AUTH_KEY=<WIF KEY>`

Specify a WIF private key to use to authenticate. The public key for this private key must be in the Teller config as an authorized client.

### Teller Key

`export TELLER_KEY=<Compressed Public Key>`

This is the public key associated with the private key that teller is using that authenticates responses from the Teller agent.

### Teller Peer Channel

`export TELLER_PEER_CHANNEL=`

This is the peer channel through which Teller agent is sent requests.

The format is `https://<peer_channel_host_domain>/api/v1/channel/<channel_id>?token=<write_token>`.

### Response Peer Channel

`export RESPONSE_PEER_CHANNEL=`

This is the peer channel through which Teller sends responses back to this client.

The format is `https://<peer_channel_host_domain>/api/v1/channel/<channel_id>?token=<read_token>`.

### Response Read Token

`export RESPONSE_READ_TOKEN=`

This is the read token for response peer channel that is used to retrieve messages posted to the peer channel.
