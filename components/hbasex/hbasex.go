package hbasex

import (
	"context"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"io"
)

type ColumnFamily string
type Key string
type Value []byte
type KV map[string][]byte
type Data map[ColumnFamily]KV

type Client struct {
	Parent gohbase.Client
}

func (c *Client) Put(ctx context.Context, table, rowKey string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (*hrpc.Result, error) {
	putRequest, err := hrpc.NewPutStr(
		ctx, table, rowKey, values, options...,
	)

	if err != nil {
		return nil, err
	}

	return c.Parent.Put(putRequest)
}

func (c *Client) RemoveByRowKey(ctx context.Context, table, rowKey string, options ...func(hrpc.Call) error) (*hrpc.Result, error) {
	removeRequest, err := hrpc.NewDelStr(
		context.Background(),
		table, rowKey, make(map[string]map[string][]byte), options...,
	)

	if err != nil {
		return nil, err
	}

	return c.Parent.Delete(removeRequest)
}

func (c *Client) Scan(table string, options ...func(hrpc.Call) error) (map[string]map[string][]byte, error) {
	req, err := hrpc.NewScanStr(
		context.Background(), table, options...,
	)

	if err != nil {
		return nil, err
	}

	scanner := c.Parent.Scan(req)

	packed := make(map[string]map[string][]byte)

	for {
		res, err := scanner.Next()
		if err == io.EOF || res == nil {
			break
		}

		var rowKey string
		temp := make(map[string][]byte)

		for _, v := range res.Cells {
			rowKey = string(v.Row)
			temp[string(v.Qualifier)] = v.Value
		}

		packed[rowKey] = temp
	}

	return packed, nil
}
