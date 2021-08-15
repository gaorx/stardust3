package sdzk

import (
	"fmt"
	"time"

	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdhash"
	"github.com/go-zookeeper/zk"
)

type NodeContent struct {
	Filename string
	Data     []byte
}

func (nc *NodeContent) String() string {
	if nc.Filename == "" {
		return "<nil>"
	}
	dataHash := ""
	if nc.Data != nil {
		dataHash = sdhash.Md5(nc.Data).HexL()
	}
	return fmt.Sprintf("%s (%dByte, MD5:%s)", nc.Filename, len(nc.Data), dataHash)
}

func (nc *NodeContent) AsText(def string) string {
	if nc.Data == nil {
		return def
	}
	return string(nc.Data)
}

func (c *Client) Content(filename string) (*NodeContent, error) {
	data, _, err := c.Get(filename)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return &NodeContent{Filename: filename, Data: data}, nil
}

func (c *Client) Contents(filenames []string, ignoreErr bool) ([]*NodeContent, error) {
	if len(filenames) == 0 {
		return []*NodeContent{}, nil
	}
	ncl := make([]*NodeContent, 0, len(filenames))
	for _, filename := range filenames {
		nc, err := c.Content(filename)
		if err != nil {
			if ignoreErr {
				continue
			} else {
				return nil, sderr.WithStack(err)
			}
		}
		ncl = append(ncl, nc)
	}
	return ncl, nil
}

func (c *Client) DirContents(dirname string, ignoreErr bool) ([]*NodeContent, error) {
	children, _, err := c.Children(dirname)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	contents, err := c.Contents(children, ignoreErr)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return contents, nil
}

func (c *Client) ContentW(filename string) (*NodeContent, <-chan *NodeContent, error) {
	data, _, startC, err := c.GetW(filename)
	if err != nil {
		return nil, nil, sderr.WithStack(err)
	}
	wc := make(chan *NodeContent)
	go func(ec <-chan zk.Event) {
		defer close(wc)
	exit:
		for {
			e := <-ec
			switch e.Type {
			case zk.EventNodeDeleted:
				break exit
			case zk.EventNodeDataChanged:
				data, _, ec0, err := c.GetW(e.Path)
				ec = ec0
				if err == nil {
					wc <- &NodeContent{Filename: e.Path, Data: data}
				}
			}
			time.Sleep(time.Millisecond * 200) // 稍微休息一下，别Busy loop了
		}
	}(startC)
	return &NodeContent{Filename: filename, Data: data}, wc, nil
}

type DirContentWatchOptions struct {
	WatchChildDataChanged bool
	StopWhenNoChildren    bool
}

func (c *Client) DirContentW(dir string, opts DirContentWatchOptions) ([]*NodeContent, <-chan []*NodeContent, error) {
	children, _, startC, err := c.ChildrenW(dir)
	if err != nil {
		return nil, nil, sderr.WithStack(err)
	}
	if opts.WatchChildDataChanged {
		// TODO: ..
		return nil, nil, sderr.New("todo: not impl")
	} else {
		ncl, err := c.Contents(JoinPathsDir(dir, children...), true)
		if err != nil {
			return nil, nil, sderr.WithStack(err)
		}
		wc := make(chan []*NodeContent)
		go func(ec <-chan zk.Event) {
			defer close(wc)
		exit:
			for {
				e := <-ec
				switch e.Type {
				case zk.EventNodeDeleted:
					break exit
				case zk.EventNodeChildrenChanged:
					children, _, ec0, err := c.ChildrenW(dir)
					ec = ec0
					if err == nil {
						ncl, err := c.Contents(JoinPathsDir(dir, children...), true)
						if err == nil {
							empty := len(ncl) == 0
							wc <- ncl
							if empty && opts.StopWhenNoChildren {
								break exit
							}
						}
					}
				}
				time.Sleep(time.Millisecond * 300)
			}
		}(startC)
		return ncl, wc, nil
	}
}
