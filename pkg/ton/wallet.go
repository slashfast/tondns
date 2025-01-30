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
	"errors"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

type WalletVersion string

const (
	V1R1               WalletVersion = "V1R1"
	V1R2               WalletVersion = "V1R2"
	V1R3               WalletVersion = "V1R3"
	V2R1               WalletVersion = "V2R1"
	V2R2               WalletVersion = "V2R2"
	V3R1               WalletVersion = "V3R1"
	V3R2               WalletVersion = "V3R2"
	V4R1               WalletVersion = "V4R1"
	V4R2               WalletVersion = "V4R2"
	HighloadV2R2       WalletVersion = "HighloadV2R2"
	HighloadV2Verified WalletVersion = "HighloadV2Verified"
	HighloadV3         WalletVersion = "HighloadV3"
	V5R1Beta           WalletVersion = "V5R1Beta"
	V5R1Final          WalletVersion = "V5R1Final"
	V5R1               WalletVersion = "V5R1"
)

func (e *WalletVersion) String() string {
	return string(*e)
}

func (e *WalletVersion) Set(v string) error {
	switch v {
	case "V1R1", "V1R2", "V1R3", "V2R1",
		"V2R2", "V3R1", "V3R2", "V4R1",
		"V4R2", "HighloadV2R2", "HighloadV2Verified",
		"HighloadV3", "V5R1Beta", "V5R1Final", "V5R1":
		*e = WalletVersion(v)
		return nil
	default:
		return errors.New(
			`must be one of 
			"V1R1", "V1R2", "V1R3", "V2R1", "V2R2", 
			"V3R1", "V3R2", "V4R1", "V4R2", "HighloadV2R2", 
			"HighloadV2Verified", "HighloadV3", "V5R1Beta", 
			"V5R1Final", "V5R1"`,
		)
	}
}

func (e *WalletVersion) Type() string {
	return "string"
}

func (e *WalletVersion) WalletType() wallet.VersionConfig {
	switch *e {
	case V1R1:
		return wallet.V1R1
	case V1R2:
		return wallet.V1R2
	case V1R3:
		return wallet.V1R3
	case V2R1:
		return wallet.V2R1
	case V2R2:
		return wallet.V2R2
	case V3R1:
		return wallet.V3R1
	case V3R2:
		return wallet.V3R2
	case V4R1:
		return wallet.V4R1
	case V4R2:
		return wallet.V4R2
	case HighloadV2R2:
		return wallet.HighloadV2R2
	case HighloadV2Verified:
		return wallet.HighloadV2Verified
	case HighloadV3:
		return wallet.HighloadV3
	case V5R1Beta:
		return wallet.ConfigV5R1Beta{
			NetworkGlobalID: wallet.MainnetGlobalID,
		}
	case V5R1Final, V5R1:
		return wallet.ConfigV5R1Final{
			NetworkGlobalID: wallet.MainnetGlobalID,
		}
	default:
		return wallet.Unknown
	}
}
