# Teller Client

Teller client provides functions to access and control a Teller agent.

## Requests

Requests are sent via peer channel to the Teller agent, then responses are sent back via a different peer channel. Message formats and serialization are defined in `messages.go`.

### Create Instrument

A create instrument request tells the Teller agent to create a new instrument to be distributed. The Teller agent then generates a new administration key and funds it with some bitcoin. It then requests a new hosted contract address. Then sends a contract offer. When it sees the corresponding contract formation on chain it sends an instrument definition. When it sees the corresponding instrument creation, it sends back a response containing the information about the instrument that was created.

### Send Tokens

A send token request tells the Teller agent to send tokens to a specified payment request. Teller responds by completing the payment request and transmitting the transaction. The Teller agent then waits for the settlement and sends back a response containing the txid of the settlement.

### Receive Tokens

A receive token request tells the Teller agent that it should receive tokens. Teller responds with a payment request to receive the tokens. The Teller agent then waits to see the payment request completed on chain, at which point it sends back a response containing the txid of the settlement.

### Reclaim Bitcoin

A reclaim bitcoin request tells the Teller agent to collect all of its UTXOs and spend them into a specified address. This can be used to close the Teller agent, or to pull bitcoin from unused instrument administration addresses. The Teller agent can continue to be used as long as the bitcoin funding address is re-funded.

## Configuration and Use

[Configure Client](./3_configure_client.md)

[Use Client](./4_use_client.md)
