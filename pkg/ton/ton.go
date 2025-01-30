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

package ton

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/dns"
	"strings"
	"time"
)

const ConfigUrl = "https://ton.org/global.config.json"

type Client struct {
	api ton.APIClientWrapped
	dns *dns.Client
}

func NewClient() (*Client, error) {
	client := liteclient.NewConnectionPool()
	err := client.AddConnectionsFromConfigUrl(context.Background(), ConfigUrl)
	if err != nil {
		panic(err)
	}

	apiClient := ton.NewAPIClient(client).WithRetry(3)

	root, err := dns.GetRootContractAddr(context.Background(), apiClient)
	if err != nil {
		return nil, err
	}

	dnsClient := dns.NewDNSClient(apiClient, root)

	return &Client{
		api: apiClient,
		dns: dnsClient,
	}, nil
}

func (c *Client) ResolveDomainAddress(domain string) (string, error) {
	if !strings.HasSuffix(domain, ".ton") {
		return "", fmt.Errorf("invalid domain \"%s\"", domain)
	}

	resolve, err := c.dns.Resolve(context.Background(), domain)
	if err != nil {
		return "", err
	}
	return resolve.GetNFTAddress().String(), nil
}

func (c *Client) LastFillUpTime(itemAddress *address.Address) (time.Time, error) {
	res, err := c.runMethod(itemAddress, "get_last_fill_up_time")
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(res.MustInt(0).Int64(), 0), err
}

func (c *Client) ItemOwner(itemAddress *address.Address) (*address.Address, error) {
	res, err := c.runMethod(itemAddress, "get_nft_data")
	if err != nil {
		return nil, err
	}

	return res.MustSlice(3).MustLoadAddr(), nil
}

func (c *Client) Api() ton.APIClientWrapped {
	return c.api
}

func (c *Client) runMethod(address *address.Address, method string, params ...any) (*ton.ExecutionResult, error) {
	block, err := c.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	return c.api.RunGetMethod(context.Background(), block, address, method, params...)
}
