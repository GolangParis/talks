package main

import "log"

// Set logger config
func Logger(logger *log.Logger, mask LogMask) func(*Client) error {
	return func(c *Client) error {
		c.Config.logger = logger
		c.Config.loggerMask = mask
		return nil
	}
}

// NewClient takes variadic options to setup it's environment and returns a *Client instance
func NewClient(options ...func(*Client) error) (*Client, error) {
	// set the default config
	c := &Client{}
	dc := defaultClientConfig()
	c.Config = dc

	// override the config with the functional options
	for _, option := range options {
		err := option(c)
		if err != nil {
			panic(err)
		}
	}
	return c, nil
}
