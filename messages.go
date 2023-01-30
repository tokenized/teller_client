package teller

import (
	"bytes"
	"fmt"

	"github.com/tokenized/channels"
	envelope "github.com/tokenized/envelope/pkg/golang/envelope/base"
	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/bsor"
	"github.com/tokenized/pkg/expanded_tx"
	"github.com/tokenized/pkg/wire"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Peer Channel flows
//
// Receive tokens from user
// 1. Receive `ReceiveTokens` from user.
// 2. Respond with `TokensToReceive` containing transfer destination.
// 3. Receive completed and signed `ExpandedTx` from user containing transfer action.
// 4. Broadcast transfer tx.
// 5. Respond with `ExpandedTx` containing settlement tx after it is received.
//
// Send tokens to user
// 1. Receive `SendTokens` from user containing payment destination.
// 2. Complete, sign, and broadcast transfer tx.
// 3. Respond with `ExpandedTx` containing transfer.
// 4. Respond with `ExpandedTx` containing settlement tx after it is received.

const (
	Version = uint8(0)

	MessageTypeInvalid           = MessageType(0)
	MessageTypeReceiveTokens     = MessageType(1)
	MessageTypeTokensToReceive   = MessageType(2)
	MessageTypeTokensReceived    = MessageType(3)
	MessageTypeSendTokens        = MessageType(4)
	MessageTypeTokensSent        = MessageType(5)
	MessageTypeCreateInstrument  = MessageType(6)
	MessageTypeInstrumentCreated = MessageType(7)
	MessageTypeReclaimBitcoin    = MessageType(8)
	MessageTypeBitcoinReclaimed  = MessageType(9)

	StatusInvalid               = uint32(1)
	StatusPaymentRequestInvalid = uint32(2)
	StatusInsufficientTxFunding = uint32(3)

	MessageGroupID = "teller"

	SubTypeReceiveTokens   = "receive-tokens"    // banking to teller
	SubTypeTokensToReceive = "tokens-to-receive" // teller to banking
	SubTypeTokensReceived  = "tokens-received"   // teller to banking

	SubTypeSendTokens = "send-tokens" // banking to teller
	SubTypeTokensSent = "tokens-sent" // teller to banking

	SubTypeCreateInstrument  = "create-instrument"  // to teller
	SubTypeInstrumentCreated = "instrument-created" // from teller

	SubTypeReclaimBitcoin   = "reclaim-bitcoin"   // to teller
	SubTypeBitcoinReclaimed = "bitcoin-reclaimed" // from teller
)

var (
	ProtocolID = envelope.ProtocolID("TELLER") // Protocol ID for teller
)

type MessageType uint8

type Protocol struct{}

func NewProtocol() *Protocol {
	return &Protocol{}
}

func (*Protocol) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (*Protocol) Parse(payload envelope.Data) (channels.Message, envelope.Data, error) {
	return Parse(payload)
}

func (*Protocol) ResponseCodeToString(code uint32) string {
	return ResponseCodeToString(code)
}

// ReceiveTokens represents a request for teller to receive tokens from a user.
// InstrumentID empty or protocol.BSVInstrumentID means send bitcoin and quantity is satoshis.
type ReceiveTokens struct {
	ID           uuid.UUID `bsor:"1" json:"id"`
	InstrumentID string    `bsor:"2" json:"instrument_id"`
	Quantity     uint64    `bsor:"3" json:"quantity"`
	Paymail      string    `bsor:"4" json:"paymail"`
}

func (*ReceiveTokens) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *ReceiveTokens) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeReceiveTokens)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

// TokensToReceive represents a response from teller with the payment request.
type TokensToReceive struct {
	ID             uuid.UUID   `bsor:"1" json:"id"` // matches ID from corresponding ReceiveTokens
	PaymentRequest *wire.MsgTx `bsor:"2" json:"payment_request"`
}

func (*TokensToReceive) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *TokensToReceive) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeTokensToReceive)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

// TokensReceived represents a response from teller that specifies the tokens for a request have
// been received.
type TokensReceived struct {
	ID   uuid.UUID      `bsor:"1" json:"id"` // matches ID from corresponding ReceiveTokens
	TxID bitcoin.Hash32 `bsor:"2" json:"txid"`

	// Tx is the settlement tx with the request tx and spent outputs included as ancestors.
	Tx *expanded_tx.ExpandedTx `bsor:"3" json:"tx"`
}

func (*TokensReceived) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *TokensReceived) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeTokensReceived)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

// SendTokens represents a request to teller to send tokens.
type SendTokens struct {
	ID             uuid.UUID   `bsor:"1" json:"id"`
	Paymail        string      `bsor:"2" json:"paymail"`
	PaymentRequest *wire.MsgTx `bsor:"3" json:"payment_request"`
	CurrencyCode   string      `bsor:"4" json:"currency_code"`
	Quantity       uint64      `bsor:"5" json:"quantity"`
}

func (*SendTokens) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *SendTokens) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeSendTokens)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

// TokensSent represents a response from teller that requested tokens have been sent.
type TokensSent struct {
	ID   uuid.UUID      `bsor:"1" json:"id"` // matches ID from corresponding SendTokens
	TxID bitcoin.Hash32 `bsor:"2" json:"txid"`

	// Tx is the settlement tx with the request tx and spent outputs included as ancestors.
	Tx *expanded_tx.ExpandedTx `bsor:"3" json:"tx"`
}

func (*TokensSent) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *TokensSent) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeTokensSent)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

type CreateInstrument struct {
	ID                         uuid.UUID `bsor:"1" json:"id"`
	CurrencyCode               string    `bsor:"2" json:"currency_code"`
	UseIdentityOracle          bool      `bsor:"3" json:"use_identity_oracle"`
	EnforcementOrdersPermitted bool      `bsor:"4" json:"enforcement_orders_permitted"`

	InstrumentType string              `bsor:"5" json:"instrument_type"`
	Precision      uint                `bsor:"6" json:"precision"`
	EntityContract *bitcoin.RawAddress `bsor:"7" json:"entity_contract"`
}

func (*CreateInstrument) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *CreateInstrument) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeCreateInstrument)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

type InstrumentCreated struct {
	ID                          uuid.UUID      `bsor:"1" json:"id"` // matches ID from corresponding CreateInstrument
	ContractLockingScript       bitcoin.Script `bsor:"2" json:"contract_locking_script"`
	AdministrationLockingScript bitcoin.Script `bsor:"3" json:"administration_locking_script"`
	InstrumentID                string         `bsor:"4" json:"instrument_id"`

	// Tx is the instrument creation tx with the request tx and spent outputs included as ancestors.
	Tx *expanded_tx.ExpandedTx `bsor:"5" json:"tx"`
}

func (*InstrumentCreated) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *InstrumentCreated) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeInstrumentCreated)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

type ReclaimBitcoin struct {
	ID      uuid.UUID       `bsor:"1" json:"id"`
	Address bitcoin.Address `bsor:"2" json:"address"`
}

func (*ReclaimBitcoin) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *ReclaimBitcoin) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeReclaimBitcoin)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

type BitcoinReclaimed struct {
	ID       uuid.UUID      `bsor:"1" json:"id"`   // matches ID from corresponding ReclaimBitcoin
	TxID     bitcoin.Hash32 `bsor:"2" json:"txid"` // hex
	Quantity uint64         `bsor:"3" json:"quantity"`

	// Tx is the bitcoin tx with spent outputs included.
	Tx *expanded_tx.ExpandedTx `bsor:"4" json:"tx"`
}

func (*BitcoinReclaimed) ProtocolID() envelope.ProtocolID {
	return ProtocolID
}

func (m *BitcoinReclaimed) Write() (envelope.Data, error) {
	// Version
	payload := bitcoin.ScriptItems{bitcoin.PushNumberScriptItem(int64(Version))}

	// Message type
	payload = append(payload, bitcoin.PushNumberScriptItem(int64(MessageTypeBitcoinReclaimed)))

	// Message
	msgScriptItems, err := bsor.Marshal(m)
	if err != nil {
		return envelope.Data{}, errors.Wrap(err, "marshal")
	}
	payload = append(payload, msgScriptItems...)

	return envelope.Data{envelope.ProtocolIDs{ProtocolID}, payload}, nil
}

func Parse(payload envelope.Data) (channels.Message, envelope.Data, error) {
	if len(payload.ProtocolIDs) == 0 {
		return nil, payload, nil
	}

	if !bytes.Equal(payload.ProtocolIDs[0], ProtocolID) {
		return nil, payload, nil
	}

	if len(payload.ProtocolIDs) != 1 {
		return nil, payload, errors.Wrapf(channels.ErrInvalidMessage, "teller messages can't wrap")
	}

	if len(payload.Payload) == 0 {
		return nil, payload, errors.Wrapf(channels.ErrInvalidMessage, "payload empty")
	}

	version, err := bitcoin.ScriptNumberValue(payload.Payload[0])
	if err != nil {
		return nil, payload, errors.Wrap(err, "version")
	}
	if version != 0 {
		return nil, payload, errors.Wrap(channels.ErrUnsupportedVersion, fmt.Sprintf("teller %d", version))
	}

	messageType, err := bitcoin.ScriptNumberValue(payload.Payload[1])
	if err != nil {
		return nil, payload, errors.Wrap(err, "message type")
	}

	result := MessageForType(MessageType(messageType))
	if result == nil {
		return nil, payload, errors.Wrap(channels.ErrNotSupported,
			fmt.Sprintf("%d", MessageType(messageType)))
	}

	payloads, err := bsor.Unmarshal(payload.Payload[2:], result)
	if err != nil {
		return nil, payload, errors.Wrap(err, "unmarshal")
	}
	payload.Payload = payloads

	return result, payload, nil
}

func MessageForType(messageType MessageType) channels.Message {
	switch MessageType(messageType) {
	case MessageTypeReceiveTokens:
		return &ReceiveTokens{}
	case MessageTypeTokensToReceive:
		return &TokensToReceive{}
	case MessageTypeTokensReceived:
		return &TokensReceived{}
	case MessageTypeSendTokens:
		return &SendTokens{}
	case MessageTypeTokensSent:
		return &TokensSent{}
	case MessageTypeCreateInstrument:
		return &CreateInstrument{}
	case MessageTypeInstrumentCreated:
		return &InstrumentCreated{}
	case MessageTypeReclaimBitcoin:
		return &ReclaimBitcoin{}
	case MessageTypeBitcoinReclaimed:
		return &BitcoinReclaimed{}
	case MessageTypeInvalid:
		return nil
	default:
		return nil
	}
}

func MessageTypeFor(message channels.Message) MessageType {
	switch message.(type) {
	case *ReceiveTokens:
		return MessageTypeReceiveTokens
	case *TokensToReceive:
		return MessageTypeTokensToReceive
	case *TokensReceived:
		return MessageTypeTokensReceived
	case *SendTokens:
		return MessageTypeSendTokens
	case *TokensSent:
		return MessageTypeTokensSent
	case *CreateInstrument:
		return MessageTypeCreateInstrument
	case *InstrumentCreated:
		return MessageTypeInstrumentCreated
	case *ReclaimBitcoin:
		return MessageTypeReclaimBitcoin
	case *BitcoinReclaimed:
		return MessageTypeBitcoinReclaimed
	default:
		return MessageTypeInvalid
	}
}

func (v *MessageType) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("Too short for MessageType : %d", len(data))
	}

	return v.SetString(string(data[1 : len(data)-1]))
}

func (v MessageType) MarshalJSON() ([]byte, error) {
	s := v.String()
	if len(s) == 0 {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}

func (v MessageType) MarshalText() ([]byte, error) {
	s := v.String()
	if len(s) == 0 {
		return nil, fmt.Errorf("Unknown MessageType value \"%d\"", uint8(v))
	}

	return []byte(s), nil
}

func (v *MessageType) UnmarshalText(text []byte) error {
	return v.SetString(string(text))
}

func (v *MessageType) SetString(s string) error {
	switch s {
	case SubTypeReceiveTokens:
		*v = MessageTypeReceiveTokens
	case SubTypeTokensToReceive:
		*v = MessageTypeTokensToReceive
	case SubTypeTokensReceived:
		*v = MessageTypeTokensReceived
	case SubTypeSendTokens:
		*v = MessageTypeSendTokens
	case SubTypeTokensSent:
		*v = MessageTypeTokensSent
	case SubTypeCreateInstrument:
		*v = MessageTypeCreateInstrument
	case SubTypeInstrumentCreated:
		*v = MessageTypeInstrumentCreated
	case SubTypeReclaimBitcoin:
		*v = MessageTypeReclaimBitcoin
	case SubTypeBitcoinReclaimed:
		*v = MessageTypeBitcoinReclaimed
	default:
		*v = MessageTypeInvalid
		return fmt.Errorf("Unknown MessageType value \"%s\"", s)
	}

	return nil
}

func (v MessageType) String() string {
	switch v {
	case MessageTypeReceiveTokens:
		return SubTypeReceiveTokens
	case MessageTypeTokensToReceive:
		return SubTypeTokensToReceive
	case MessageTypeTokensReceived:
		return SubTypeTokensReceived
	case MessageTypeSendTokens:
		return SubTypeSendTokens
	case MessageTypeTokensSent:
		return SubTypeTokensSent
	case MessageTypeCreateInstrument:
		return SubTypeCreateInstrument
	case MessageTypeInstrumentCreated:
		return SubTypeInstrumentCreated
	case MessageTypeReclaimBitcoin:
		return SubTypeReclaimBitcoin
	case MessageTypeBitcoinReclaimed:
		return SubTypeBitcoinReclaimed
	default:
		return ""
	}
}

func ResponseCodeToString(code uint32) string {
	switch code {
	case StatusInvalid:
		return "invalid"
	case StatusPaymentRequestInvalid:
		return "payment_request_invalid"
	case StatusInsufficientTxFunding:
		return "insufficient_tx_funding"
	default:
		return "parse_error"
	}
}
