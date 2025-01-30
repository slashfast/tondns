/*
 * Copyright (c) 2025 Vladislav Trofimenko <com@slashfast.dev>
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the
 * "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish,
 * distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to
 * the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 * LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
 * OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package gems

import (
	"context"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
)

const Endpoint = "https://api.getgems.io/graphql"
const firstWindowSize = 100

type Client struct {
	client graphql.Client
	config *Config
}

func NewClient(config *Config) *Client {
	httpClient := &http.Client{}

	if config.Proxy.String() != "" {
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(config.Proxy.URL),
		}
	}

	client := graphql.NewClient(Endpoint, httpClient)

	return &Client{
		client: client,
		config: config,
	}
}

func (c *Client) ItemsByOwner(address string) ([]ItemsNftItemsByOwnerNftItemConnectionItemsNftItem, error) {
	var items []ItemsNftItemsByOwnerNftItemConnectionItemsNftItem
	cursor := ""

	for {
		res, err := Items(context.Background(), c.client, address, firstWindowSize, cursor)
		if err != nil {
			return nil, err
		}
		items = append(items, res.NftItemsByOwner.Items...)

		cursor = res.NftItemsByOwner.Cursor
		if cursor == "" {
			break
		}

		time.Sleep(c.config.QueryDelay)
	}

	return items, nil
}

func (c *Client) ItemHistory(address string) ([]HistoryHistoryNftItemNftItemHistoryConnectionItemsNftItemHistory, error) {
	var items []HistoryHistoryNftItemNftItemHistoryConnectionItemsNftItemHistory
	cursor := ""

	for {
		res, err := History(context.Background(), c.client, address, firstWindowSize, cursor)
		if err != nil {
			return nil, err
		}
		items = append(items, res.HistoryNftItem.Items...)

		cursor = res.HistoryNftItem.Cursor
		if cursor == "" {
			break
		}

		time.Sleep(c.config.QueryDelay)
	}

	return items, nil
}
