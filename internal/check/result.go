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

type Result struct {
	Address         string `json:"address,omitempty"`
	DomainName      string `json:"domainName,omitempty"`
	IsNeedToRenew   bool   `json:"isNeedToRenew,omitempty"`
	LastFillUp      string `json:"lastFillUp,omitempty"`
	ExpiringOn      string `json:"expiringOn,omitempty"`
	DaysToExpire    int    `json:"daysToExpire,omitempty"`
	MintDate        string `json:"mintDate,omitempty"`
	IsAssigned      bool   `json:"isAssigned,omitempty"`
	IsMine          bool   `json:"isMine,omitempty"`
	OwnerAddress    string `json:"ownerAddress,omitempty"`
	OwnerRawAddress string `json:"ownerRawAddress,omitempty"`

	HasRenewedNow bool   `json:"hasRenewedNow,omitempty"`
	RenewalAmount int64  `json:"renewalAmount,omitempty"`
	Error         string `json:"error,omitempty"`
}

func NewResultError(domain, error string) Result {
	return Result{
		DomainName: domain,
		Error:      error,
	}
}
