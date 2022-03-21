// Copyright (c) 2022, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feedfinder

import (
	"fmt"
	"net/http"
	"time"
)

var defaultHTTPClient = &http.Client{
	Timeout: 2 * time.Minute,
}

func httpGet(rawURL string) (_ *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("HTTP response status %s", resp.Status)
	}

	return resp, nil
}
