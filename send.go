package teller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tokenized/channels"
	channelsExpandedTx "github.com/tokenized/channels/expanded_tx"
	envelopeV1 "github.com/tokenized/envelope/pkg/golang/envelope/v1"
	"github.com/tokenized/logger"
	"github.com/tokenized/pkg/bitcoin"
	"github.com/tokenized/pkg/expanded_tx"
	"github.com/tokenized/pkg/peer_channels"
	"github.com/tokenized/specification/dist/golang/print"
	"github.com/tokenized/threads"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrFailure        = errors.New("Failure")
	ErrWrongMessageID = errors.New("Wrong Message ID")
)

func ProcessRequest(ctx context.Context, requestID uuid.UUID, msg channels.Writer) error {
	peerChannelsFactory := peer_channels.NewFactory()

	key, err := bitcoin.KeyFromStr(os.Getenv("AUTH_KEY"))
	if err != nil {
		return errors.Wrap(err, "auth key")
	}

	tellerKey, err := bitcoin.PublicKeyFromStr(os.Getenv("TELLER_KEY"))
	if err != nil {
		return errors.Wrap(err, "teller key")
	}

	tellerPeerChannel, err := peer_channels.ParseChannel(os.Getenv("TELLER_PEER_CHANNEL"))
	if err != nil {
		return errors.Wrap(err, "teller peer channel")
	}

	tellerClient, err := peerChannelsFactory.NewClient(tellerPeerChannel.BaseURL)
	if err != nil {
		return errors.Wrap(err, "teller peer channel client")
	}

	responsePeerChannel, err := peer_channels.ParseChannel(os.Getenv("RESPONSE_PEER_CHANNEL"))
	if err != nil {
		return errors.Wrap(err, "teller peer channel")
	}

	responseReadToken, err := uuid.Parse(os.Getenv("RESPONSE_READ_TOKEN"))

	responseClient, err := peerChannelsFactory.NewClient(responsePeerChannel.BaseURL)
	if err != nil {
		return errors.Wrap(err, "response peer channel client")
	}

	replyTo := &channels.ReplyTo{
		PeerChannel: responsePeerChannel,
	}

	msgScript, err := wrapMessage(msg, key, channels.RandomHash(), replyTo, requestID)
	if err != nil {
		return errors.Wrap(err, "wrap message")
	}

	if err := tellerClient.WriteMessage(ctx, tellerPeerChannel.ChannelID, tellerPeerChannel.Token,
		peer_channels.ContentTypeBinary, bytes.NewReader(msgScript)); err != nil {
		return errors.Wrap(err, "post message")
	}

	js, _ := json.MarshalIndent(msg, "", "  ")
	fmt.Printf("Sent request : %s\n", js)

	var wait sync.WaitGroup

	incoming := make(chan peer_channels.Message, 10)
	listenThread, listenComplete := threads.NewInterruptableThreadComplete("Listen Peer Channel",
		func(ctx context.Context, interrupt <-chan interface{}) error {
			return responseClient.Listen(ctx, responseReadToken.String(), true, time.Second,
				incoming, interrupt)
		}, &wait)

	listenThread.Start(ctx)

	var handleErr error
	var timeout error
	done := false
	for !done {
		select {
		case msg := <-incoming:
			handleErr = handleMessage(ctx, tellerKey, requestID, msg)
			if handleErr != nil {
				if errors.Cause(handleErr) == ErrFailure {
					logger.Error(ctx, handleErr.Error())
					done = true
				}
			} else {
				logger.Info(ctx, "Handled message")
				done = true
			}

			if handleErr == nil || errors.Cause(handleErr) == ErrWrongMessageID ||
				errors.Cause(handleErr) == ErrFailure {
				if err := responseClient.MarkMessages(ctx, msg.ChannelID,
					responseReadToken.String(), msg.Sequence, true, true); err != nil {
					logger.Error(ctx, "Failed to mark message as read : %s", err)
				}
				handleErr = nil
			}

		case listenErr := <-listenComplete:
			logger.Error(ctx, "Listen completed : %s", listenErr)
			done = true

		case <-time.After(time.Second * 10):
			timeout = errors.New("Timed out")
			done = true
		}
	}

	listenThread.Stop(ctx)
	wait.Wait()

	combinedErr := threads.CombineErrors(listenThread.Error(), handleErr, timeout)
	if errors.Cause(combinedErr) != threads.Interrupted {
		return combinedErr
	}

	return nil
}

func Listen(ctx context.Context, requestID uuid.UUID) error {
	peerChannelsFactory := peer_channels.NewFactory()

	tellerKey, err := bitcoin.PublicKeyFromStr(os.Getenv("TELLER_KEY"))
	if err != nil {
		return errors.Wrap(err, "teller key")
	}

	responsePeerChannel, err := peer_channels.ParseChannel(os.Getenv("RESPONSE_PEER_CHANNEL"))
	if err != nil {
		return errors.Wrap(err, "teller peer channel")
	}

	responseReadToken, err := uuid.Parse(os.Getenv("RESPONSE_READ_TOKEN"))

	responseClient, err := peerChannelsFactory.NewClient(responsePeerChannel.BaseURL)
	if err != nil {
		return errors.Wrap(err, "response peer channel client")
	}

	var wait sync.WaitGroup

	incoming := make(chan peer_channels.Message, 10)
	listenThread, listenComplete := threads.NewInterruptableThreadComplete("Listen Peer Channel",
		func(ctx context.Context, interrupt <-chan interface{}) error {
			return responseClient.Listen(ctx, responseReadToken.String(), true, time.Second,
				incoming, interrupt)
		}, &wait)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	listenThread.Start(ctx)

	var handleErr error
	var timeout error
	done := false
	for !done {
		select {
		case msg := <-incoming:
			handleErr = handleMessage(ctx, tellerKey, requestID, msg)
			if handleErr != nil {
				if errors.Cause(handleErr) == ErrFailure {
					logger.Error(ctx, handleErr.Error())
					done = true
				}
			} else {
				logger.Info(ctx, "Handled message")
				done = true
			}

			if handleErr == nil || errors.Cause(handleErr) == ErrWrongMessageID ||
				errors.Cause(handleErr) == ErrFailure {
				if err := responseClient.MarkMessages(ctx, msg.ChannelID,
					responseReadToken.String(), msg.Sequence, true, true); err != nil {
					logger.Error(ctx, "Failed to mark message as read : %s", err)
				}
				handleErr = nil
			}

		case listenErr := <-listenComplete:
			logger.Error(ctx, "Listen completed : %s", listenErr)
			done = true

		case <-time.After(time.Minute * 10):
			timeout = errors.New("Timed out")
			done = true

		case <-osSignals:
			logger.Info(ctx, "Start shutdown...")
			done = true
		}
	}

	listenThread.Stop(ctx)
	wait.Wait()

	combinedErr := threads.CombineErrors(listenThread.Error(), handleErr, timeout)
	if errors.Cause(combinedErr) != threads.Interrupted {
		return combinedErr
	}

	return nil
}

func handleMessage(ctx context.Context, tellerKey bitcoin.PublicKey, requestID uuid.UUID,
	msg peer_channels.Message) error {

	protocols := channels.NewProtocols(NewProtocol(), channelsExpandedTx.NewProtocol())

	if msg.ContentType != peer_channels.ContentTypeBinary {
		return fmt.Errorf("Unsupported content type : %s", msg.ContentType)
	}

	wMsg, err := unwrapMessage(protocols, msg.Payload)
	if err != nil {
		return errors.Wrap(err, "unwrap message")
	}

	if wMsg.Signature == nil {
		return errors.New("Missing Signature")
	}

	if wMsg.ID == nil {
		return errors.New("Missing ID")
	}

	if !bytes.Equal(wMsg.ID[:], requestID[:]) {
		return errors.Wrapf(ErrWrongMessageID, "got %s, want %s", *wMsg.ID, requestID)
	}

	if wMsg.Signature.PublicKey == nil {
		wMsg.Signature.SetPublicKey(&tellerKey)
	} else {
		if !wMsg.Signature.PublicKey.Equal(tellerKey) {
			return fmt.Errorf("Wrong Signature Key : teller_key %s, sig_key %s", tellerKey,
				wMsg.Signature.PublicKey)
		}
	}

	if wMsg.Response != nil {
		js, _ := json.MarshalIndent(wMsg.Response, "", "  ")
		fmt.Printf("Response : %s\n", js)
	}

	if err := wMsg.Signature.Verify(); err != nil {
		return err
	}

	if wMsg.Message != nil {
		switch m := wMsg.Message.(type) {
		case *channelsExpandedTx.ExpandedTxMessage:
			etxd := expanded_tx.ExpandedTx(*m)
			etx := &etxd
			fmt.Printf("Expanded Transaction : %s", etx)

			if etx.Tx != nil {
				print.PrintActions(etx.Tx)
			}

		case *TokensToReceive:
			js, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("TokensToReceive : %s\n", js)

			if m.PaymentRequest != nil {
				print.PrintActions(m.PaymentRequest)
			}

		case *TokensReceived:
			js, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("TokensReceived : %s\n", js)

			if m.Tx.Tx != nil {
				fmt.Printf("TxID : %s\n", m.Tx.Tx.TxHash())
				print.PrintActions(m.Tx.Tx)
			}

		case *TokensSent:
			js, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("TokensSent : %s\n", js)

			if m.Tx.Tx != nil {
				fmt.Printf("TxID : %s\n", m.Tx.Tx.TxHash())
				print.PrintActions(m.Tx.Tx)
			}

		case *InstrumentCreated:
			js, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("InstrumentCreated : %s\n", js)

			if m.Tx.Tx != nil {
				fmt.Printf("TxID : %s\n", m.Tx.Tx.TxHash())
				print.PrintActions(m.Tx.Tx)
			}

		default:
			js, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("Message : %s\n", js)
		}
	}

	if wMsg.Message == nil {
		if wMsg.Response == nil {
			return errors.New("No Payload and No Response")
		}
	}

	return nil
}

type WrappedMessage struct {
	Signature *channels.Signature
	Response  *channels.Response
	ID        *uuid.UUID
	Message   channels.Message
}

func wrapMessage(msg channels.Writer, key bitcoin.Key, hash bitcoin.Hash32,
	replyTo *channels.ReplyTo, id uuid.UUID) (bitcoin.Script, error) {

	payload, err := msg.Write()
	if err != nil {
		return nil, errors.Wrap(err, "write")
	}

	payload, err = channels.WrapUUID(payload, id)
	if err != nil {
		return nil, errors.Wrap(err, "uuid")
	}

	if replyTo != nil {
		payload, err = replyTo.Wrap(payload)
		if err != nil {
			return nil, errors.Wrap(err, "reply to")
		}
	}

	payload, err = channels.WrapSignature(payload, key, &hash, true)
	if err != nil {
		return nil, errors.Wrap(err, "sign")
	}

	return envelopeV1.Wrap(payload).Script()
}

func unwrapMessage(protocols *channels.Protocols, script []byte) (*WrappedMessage, error) {
	payload, err := envelopeV1.Parse(bytes.NewReader(script))
	if err != nil {
		return nil, errors.Wrap(err, "envelope")
	}

	result := &WrappedMessage{}
	result.Signature, payload, err = channels.ParseSigned(payload)
	if err != nil {
		return nil, errors.Wrap(err, "sign")
	}

	var channelsUUID *channels.UUID
	channelsUUID, payload, err = channels.ParseUUID(payload)
	if err != nil {
		return nil, errors.Wrap(err, "uuid")
	}
	uuidValue := uuid.UUID(*channelsUUID)
	result.ID = &uuidValue

	result.Response, payload, err = channels.ParseResponse(payload)
	if err != nil {
		return nil, errors.Wrap(err, "response")
	}

	if len(payload.ProtocolIDs) == 0 {
		return result, nil
	}

	if len(payload.ProtocolIDs) > 1 {
		return nil, errors.Wrap(channels.ErrNotSupported, "more than one data protocol")
	}

	protocol := protocols.GetProtocol(payload.ProtocolIDs[0])
	if protocol == nil {
		return nil, errors.Wrap(channels.ErrUnsupportedProtocol, payload.ProtocolIDs[0].String())
	}

	msg, _, err := protocol.Parse(payload)
	if err != nil {
		return nil, errors.Wrap(err, "parse")
	}
	result.Message = msg

	return result, nil
}
