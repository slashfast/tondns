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
	"strings"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

type WalletVersion string

const (
	V1R1               WalletVersion = "v1r1"
	V1R2               WalletVersion = "v1r2"
	V1R3               WalletVersion = "v1r3"
	V2R1               WalletVersion = "v2r1"
	V2R2               WalletVersion = "v2r2"
	V3R1               WalletVersion = "v3r1"
	V3R2               WalletVersion = "v3r2"
	V4R1               WalletVersion = "v4r1"
	V4R2               WalletVersion = "v4r2"
	HighloadV2R2       WalletVersion = "highloadv2r2"
	HighloadV2Verified WalletVersion = "highloadv2verified"
	HighloadV3         WalletVersion = "highloadv3"
	V5R1Beta           WalletVersion = "v5r1beta"
	V5R1Final          WalletVersion = "v5r1final"
	V5R1               WalletVersion = "v5r1"
)

func (e *WalletVersion) String() string {
	return string(*e)
}

func (e *WalletVersion) Set(v string) error {
	v = strings.ToLower(v)
	switch v {
	case "v1r1", "v1r2", "v1r3",
		"v2r1", "v2r2", "v3r1",
		"v3r2", "v4r1", "v4r2",
		"highloadv2r2", "highloadv2verified", "highloadv3",
		"v5r1beta", "v5r1final", "v5r1":
		*e = WalletVersion(v)
		return nil
	default:
		return errors.New("must be one of v1r1, v1r2, v1r3, v2r1, v2r2, v3r1, v3r2, v4r1, v4r2, highloadv2r2, highloadv2verified, highloadv3, v5r1beta, v5r1final, v5r1")
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
