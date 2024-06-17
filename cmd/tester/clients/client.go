package clients

import (
	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/model"
)

type Id model.Id
type Status model.Id
type Record model.Record
type Time int

type Client interface {
	Delay(t Time) error
	Forward(t Time) error
	Put(records []Record) ([]Status, error)
	Get([]Id) ([]Status, error)
	Del([]Id) ([]Status, error)

	SoftTTL() Time
	HardTTL() Time
}

type client struct {
	softTtl Time
	hardTtl Time
	s       *service
	w       *workers
}

func NewClient(softTtl Time, hardTtl Time, s *service, w *workers) Client {
	return &client{
		softTtl: softTtl,
		hardTtl: hardTtl,
		s:       s,
		w:       w,
	}
}

func (c *client) Delay(t Time) error {
	return c.w.Delay(t)
}

func (c *client) Forward(t Time) error {
	return c.s.Forward(t)
}

func (c *client) Put(records []Record) ([]Status, error) {
	return c.w.Put(records)
}

func (c *client) Get(ids []Id) ([]Status, error) {
	return c.s.Get(ids)
}

func (c *client) Del(ids []Id) ([]Status, error) {
	return c.w.Del(ids)
}

func (c *client) SoftTTL() Time {
	return c.softTtl
}

func (c *client) HardTTL() Time {
	return c.hardTtl
}
