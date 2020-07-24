package subscriber

import (
	"github.com/terra-project/mantle/types"
	"net/http"
)

const CloseRequest = 1
const CloseResponse = 2

type HTTPSubscription struct {
	getter HTTPSubscriptionGetter
	control chan int
}

type HTTPSubscriptionGetter func(height uint64) string

func NewHTTPSubscription(getter HTTPSubscriptionGetter) Subscriber {
	return &HTTPSubscription{
		getter: getter,
		control: make(chan int),
	}
}

func (c *HTTPSubscription) Subscribe() chan types.Block {
	ch := make(chan types.Block)

	go func() {
		for {
			select {
			case control := <-c.control:
				if control == CloseRequest {
					c.control <- CloseResponse
					return
				}
			default:
				// make request
				resp, err := http.Get(c.getter())
			}
		}
	}()

	return ch
}

func (c *HTTPSubscription) Close() error {
	c.control <- CloseRequest
	select {
		case control := <-c.control:
			if control == CloseResponse {
				return nil
			}
	}
	return nil
}

