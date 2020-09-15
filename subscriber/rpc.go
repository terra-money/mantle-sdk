package subscriber

import (
	"encoding/json"
	"errors"
	"log"

	websocket "github.com/gorilla/websocket"
	types "github.com/terra-project/mantle/types"
)

type (
	RPCSubscription struct {
		ws          *websocket.Conn
		id          int
		initialized bool
	}
	Request struct {
		JSONRPC string                    `json:"jsonrpc"`
		Method  string                    `json:"method"`
		ID      int                       `json:"id"`
		Params  SubscriptionJsonRpcParams `json:"params"`
	}
	SubscriptionJsonRpcParams struct {
		Query string `json:"query"`
	}
)

func NewRpcSubscription(endpoint string) *RPCSubscription {
	log.Print("Opening websocket...")
	ws, _, err := websocket.DefaultDialer.Dial(endpoint, nil)

	// panic as this should not fail
	if err != nil {
		panic(err)
	}

	return &RPCSubscription{
		ws:          ws,
		id:          0,
		initialized: false,
	}
}

func (c *RPCSubscription) Close() error {
	return c.ws.Close()
}

// Subscribe starts listening to tendermint RPC. It **must** be run as goroutine.
func (c *RPCSubscription) Subscribe() chan types.Block {
	var request = &Request{
		JSONRPC: "2.0",
		Method:  "subscribe",
		ID:      c.id,
		Params: SubscriptionJsonRpcParams{
			Query: "tm.event = 'NewBlock'",
		},
	}

	log.Print("Subscribing to tendermint rpc...")

	// should not fail here
	if err := c.ws.WriteJSON(request); err != nil {
		panic(err)
	}

	// handle initial message
	// by setting c.initialized to true, we prevent message mishandling
	if c.handleInitialHandhake() != nil {
		c.initialized = true
	}

	log.Print("Subscription and the first handshake done. Receiving blocks...")

	// make channel for receiving block events
	channel := make(chan types.Block)

	// run event receiver
	go c.receiveBlockEvents(channel)

	return channel
}

// tendermint rpc sends the "subscription ok" for the intiail response
// filter that out by only sending through channel when there is
// "data" field present
func (c *RPCSubscription) handleInitialHandhake() error {
	_, message, err := c.ws.ReadMessage()

	if err != nil {
		return err
	}

	if _, err := sanitizeMessage(message); err != nil {
		return err
	}

	return nil
}

// TODO: handle errors here
func (c *RPCSubscription) receiveBlockEvents(onBlock chan types.Block) {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			panic(err)
		}

		sanitized, err := sanitizeMessage(message)
		if err != nil {
			panic(err)
		}

		// unmarshaling to BlockEvent failed, handle me
		blockEvent := types.Block{}
		if err := json.Unmarshal([]byte(sanitized), &blockEvent); err != nil {
			panic(err)
		}

		// send!
		onBlock <- blockEvent
	}
}

func sanitizeMessage(message []byte) ([]byte, error) {
	// tendermint rpc sends the "subscription ok" for the intiail response
	// filter that out by only sending through channel when there is
	// "data" field present

	var response map[string]*json.RawMessage
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, err
	}

	var result map[string]*json.RawMessage
	if err := json.Unmarshal(*response["result"], &result); err != nil {
		return nil, err
	}

	// return data part of result if it exists
	if _, ok := result["data"]; ok {
		return *result["data"], nil
	}

	// if code ever reaches here it's an unknown case
	// return error
	return nil, errors.New("Unknown RPC response received")
}

// TODO: get specific block
