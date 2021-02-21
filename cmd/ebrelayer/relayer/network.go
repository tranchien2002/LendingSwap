package relayer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/trinhtan/horizon-hackathon/hmyclient"
)

// IsWebsocketURL returns true if the given URL is a websocket URL
func IsWebsocketURL(rawurl string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}
	return u.Scheme == "ws" || u.Scheme == "wss"
}

// EthSetupWebsocketClient returns boolean indicating if a URL is valid websocket ethclient
func EthSetupWebsocketClient(ethURL string) (*ethclient.Client, error) {
	if strings.TrimSpace(ethURL) == "" {
		return nil, nil
	}

	if !IsWebsocketURL(ethURL) {
		return nil, fmt.Errorf("invalid websocket eth client URL: %s", ethURL)
	}

	client, err := ethclient.Dial(ethURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// HmySetupWebsocketClient returns boolean indicating if a URL is valid websocket hmyclient
func HmySetupWebsocketClient(hmyURL string) (*hmyclient.Client, error) {
	if strings.TrimSpace(hmyURL) == "" {
		return nil, nil
	}

	if !IsWebsocketURL(hmyURL) {
		return nil, fmt.Errorf("invalid websocket eth client URL: %s", hmyURL)
	}

	client, err := hmyclient.Dial(hmyURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}
