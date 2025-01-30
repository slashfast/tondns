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

package check

import (
	"context"
	"fmt"
	"time"
	"github.com/slashfast/tondns/internal/config"
	"github.com/slashfast/tondns/pkg/gems"
	"github.com/slashfast/tondns/pkg/ton"

	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/dns"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const (
	DomainCollection = "EQC3dNlesgVD8YbAazcauIrXBPfiVhMMr5YYk2in0Mtsz0Bz"
	RenewPrice       = "0.005"
)

var renewAmount = tlb.MustFromTON(RenewPrice)

type Checker struct {
	client   *ton.Client
	wallet   *wallet.Wallet
	config   config.Config
	addr     string
	gems     *gems.Client
	itemsIdx map[string]gems.ItemsNftItemsByOwnerNftItemConnectionItemsNftItem
}

func NewChecker(cfg config.Config) (*Checker, error) {
	client, err := ton.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create TON client: %w", err)
	}

	gemsClient := gems.NewClient(&gems.Config{Proxy: cfg.Proxy, QueryDelay: cfg.CheckDelay})

	checker := &Checker{
		client: client,
		gems:   gemsClient,
		config: cfg,
	}

	if err := checker.initWallet(); err != nil {
		return nil, err
	}

	if !cfg.Lite {
		if err := checker.loadItems(); err != nil {
			return nil, err
		}
	}

	return checker, nil
}

func (c *Checker) initWallet() error {
	if !c.config.Lite {
		w, err := wallet.FromSeed(c.client.Api(), c.config.Seed, c.config.WalletVersion.WalletType())
		if err != nil {
			return fmt.Errorf("failed to initialize wallet: %w", err)
		}
		c.wallet = w
		c.addr = w.WalletAddress().Bounce(true).String()
	}
	return nil
}

func (c *Checker) loadItems() error {
	items, err := c.gems.ItemsByOwner(c.addr)
	if err != nil {
		return fmt.Errorf("failed to load items: %w", err)
	}

	c.itemsIdx = make(map[string]gems.ItemsNftItemsByOwnerNftItemConnectionItemsNftItem)
	for _, item := range items {
		if item.Collection.Address == DomainCollection {
			c.itemsIdx[item.Name] = item
		}
	}
	return nil
}

func (c *Checker) SmartCheck(domain string) Result {
	if c.config.Lite {
		panic("unsupported in lite mode")
	}

	item, isDomainMine := c.itemsIdx[domain]
	if !isDomainMine {
		itemAddr, err := c.client.ResolveDomainAddress(domain)
		if err != nil {
			return NewResultError(domain, err.Error())
		}
		item = gems.ItemsNftItemsByOwnerNftItemConnectionItemsNftItem{
			Address:    itemAddr,
			Name:       domain,
			Collection: gems.ItemsNftItemsByOwnerNftItemConnectionItemsNftItemCollectionNftCollection{Address: DomainCollection},
		}
	}

	return c.checkItem(item, isDomainMine)
}

func (c *Checker) SmartCheckAll() []Result {
	if c.config.Lite {
		panic("unsupported in lite mode")
	}

	var results []Result
	for _, item := range c.itemsIdx {
		results = append(results, c.checkItem(item, true))
	}
	return results
}

func (c *Checker) checkItem(item gems.ItemsNftItemsByOwnerNftItemConnectionItemsNftItem, isDomainMine bool) Result {
	history, err := c.gems.ItemHistory(item.Address)
	if err != nil {
		return NewResultError(item.Name, err.Error())
	}

	itemAddr := address.MustParseAddr(item.Address)
	mintDate := time.Unix(int64(history[len(history)-1].CreatedAt), 0)
	lastFillUpDate, err := c.client.LastFillUpTime(itemAddr)
	if err != nil {
		return NewResultError(item.Name, err.Error())
	}
	expiresOnDate := lastFillUpDate.Add(365 * 24 * time.Hour)
	timeToExpires := time.Until(expiresOnDate)
	isNeedToRenew := timeToExpires < c.config.Threshold
	itemOwner, err := c.client.ItemOwner(itemAddr)
	if err != nil {
		return NewResultError(item.Name, err.Error())
	}
	isAssigned := c.addr == itemOwner.String()

	result := Result{
		Address:       itemAddr.String(),
		DomainName:    item.Name,
		IsNeedToRenew: isNeedToRenew,
		LastFillUp:    lastFillUpDate.Format("02.01.2006"),
		ExpiringOn:    expiresOnDate.Format("02.01.2006"),
		DaysToExpire:  int(timeToExpires.Hours() / 24),
		MintDate:      mintDate.Format("02.01.2006"),
		IsAssigned:    isAssigned,
		IsMine:        isDomainMine,
		OwnerAddress:  itemOwner.String(),
	}

	isRenewAllowed := !c.config.Lite
	isRenewAllowed = isRenewAllowed && isNeedToRenew
	isRenewAllowed = isRenewAllowed && c.config.Renew
	isRenewAllowed = isRenewAllowed || c.config.Force
	isRenewAllowed = isRenewAllowed && isDomainMine

	if isRenewAllowed {
		if isAssigned {
			err := c.Renew(itemAddr)
			if err != nil {
				return NewResultError(item.Name, err.Error())
			}
			result.HasRenewedNow = true
		} else {
			err := c.Assign(itemAddr)
			if err != nil {
				return NewResultError(item.Name, err.Error())
			}
			result.IsAssigned = true
			result.HasRenewedNow = true
		}
	}

	return result
}

func (c *Checker) Assign(address *address.Address) error {
	if c.config.Lite {
		panic("unsupported in lite mode")
	}

	payload := (&dns.Domain{}).BuildSetWalletRecordPayload(address)
	return c.wallet.Send(
		context.Background(),
		wallet.SimpleMessage(address, renewAmount, payload), c.config.WaitTx,
	)
}

func (c *Checker) Renew(address *address.Address) error {
	if c.config.Lite {
		panic("unsupported in lite mode")
	}

	return c.wallet.Transfer(
		context.Background(),
		address, renewAmount, "Renew", c.config.WaitTx,
	)
}

func (c *Checker) Check(domain string) Result {
	item := gems.ItemsNftItemsByOwnerNftItemConnectionItemsNftItem{
		Name:       domain,
		Collection: gems.ItemsNftItemsByOwnerNftItemConnectionItemsNftItemCollectionNftCollection{Address: DomainCollection},
	}

	itemAddr, err := c.client.ResolveDomainAddress(domain)
	if err != nil {
		return NewResultError(item.Name, err.Error())
	}
	item.Address = itemAddr

	return c.checkItem(item, false)
}
