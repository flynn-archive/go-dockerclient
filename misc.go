// Copyright 2014 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"bytes"
	"github.com/fsouza/go-dockerclient/engine"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Version returns version information about the docker server.
//
// See http://goo.gl/IqKNRE for more details.
func (c *Client) Version() (*engine.Env, error) {
	body, _, err := c.do("GET", "/version", nil)
	if err != nil {
		return nil, err
	}
	out := engine.NewOutput()
	remoteVersion, err := out.AddEnv()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(out, bytes.NewReader(body)); err != nil {
		return nil, err
	}
	return remoteVersion, nil
}

// Info returns system-wide information, like the number of running containers.
//
// See http://goo.gl/LOmySw for more details.
func (c *Client) Info() (*engine.Env, error) {
	body, _, err := c.do("GET", "/info", nil)
	if err != nil {
		return nil, err
	}
	var info engine.Env
	err = info.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return &info, nil
}


type EventErr struct {
        Code    int    `json:"code,omitempty"`
        Message string `json:"message,omitempty"`
}

type Event struct {
        Status   string    `json:"status,omitempty"`
        Progress string    `json:"progress,omitempty"`
        ID       string    `json:"id,omitempty"`
        From     string    `json:"from,omitempty"`
        Time     int64     `json:"time,omitempty"`
        Error    *EventErr `json:"errorDetail,omitempty"`
}

type EventStream struct {
	Events chan *Event
	// Error must only be read after Stream has closed
	Error error

	conn io.ReadCloser
}

func (s *EventStream) Close() error { return s.conn.Close() }

func (s *EventStream) stream() {
	decoder := json.NewDecoder(s.conn)
	for {
		event := &Event{}
		if err := decoder.Decode(event); err != nil {
			if err == io.EOF {
				err = nil
			}
			s.Error = err
			break
		}
		s.Events <- event
	}
	close(s.Events)
	s.conn.Close()
}

// Events returns a stream of container events.
func (c *Client) Events() (*EventStream, error) {
	req, err := http.NewRequest("GET", c.getURL("/events"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := c.client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			err = ErrConnectionRefused
		}
		return nil, err
	}
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return nil, newError(res.StatusCode, body)
	}

	stream := &EventStream{
		Events: make(chan *Event),
		conn:   res.Body,
	}
	go stream.stream()
	return stream, nil
}
