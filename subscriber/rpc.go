package subscriber

import (
	"encoding/json"
	"github.com/terra-project/mantle-sdk/utils"
	"log"

	websocket "github.com/gorilla/websocket"
	types "github.com/terra-project/mantle-sdk/types"
)

type (
	RPCSubscription struct {
		ws          *websocket.Conn
		id          int
		initialized bool
		onWsError   OnWsError
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
	OnWsError func(error)
)

func NewRpcSubscription(endpoint string, onWsError func(err error)) *RPCSubscription {
	log.Print("Opening websocket...")
	ws, _, err := websocket.DefaultDialer.Dial(endpoint, nil)

	// panic as this should not fail
	if err != nil {
		panic(err)
	}

	ws.SetCloseHandler()

	return &RPCSubscription{
		ws:          ws,
		onWsError:   onWsError,
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
	_, _, err := c.ws.ReadMessage()

	if err != nil {
		return err
	}

	return nil
}

// TODO: handle errors here
func (c *RPCSubscription) receiveBlockEvents(onBlock chan types.Block) {
	for {
		_, message, err := c.ws.ReadMessage()

		if err != nil {
			if c.onWsError != nil {
				c.onWsError(err)
			} else {
				panic(err)
			}
		}

		data := new(struct {
			Result struct {
				Data struct {
					Value struct {
						Block json.RawMessage
					} `json:"value"`
				} `json:"data"`
			} `json:"result"`
		})

		if unmarshalErr := json.Unmarshal(message, data); unmarshalErr != nil {
			panic(unmarshalErr)
		}

		block := utils.ConvertBlockHeaderToTMHeader(data.Result.Data.Value.Block)

		// send!
		onBlock <- block
	}
}

// TODO: get specific block
