package sdzk

import (
	"strings"

	"github.com/gaorx/stardust3/sderr"
	"github.com/go-zookeeper/zk"
)

func (c *Client) MkdirP(dirname string, flags int32, acl []zk.ACL) error {
	paths, dir := strings.Split(dirname, "/"), ""
	for _, p := range paths {
		if p == "" {
			continue
		}
		dir = JoinPath(dir, p)
		if dir == "/" {
			continue
		}
		_, err := c.Create(dir, []byte("dir"), flags, acl)
		if err != nil && err != zk.ErrNodeExists {
			return sderr.WithStack(err)
		}
	}
	return nil
}

func (c *Client) SetOrCreate(filename string, data []byte, createFlags int32, createACL []zk.ACL) error {
	ok, _, err := c.Exists(filename)
	if err != nil {
		return sderr.WithStack(err)
	}
	if ok {
		_, err = c.Set(filename, data, -1)
		return sderr.WithStack(err)
	} else {
		_, err := c.Create(filename, data, createFlags, createACL)
		return sderr.WithStack(err)
	}
}

func JoinPath(paths ...string) string {
	paths1 := make([]string, 0)
	for _, p := range paths {
		if p != "" {
			p = strings.TrimSuffix(strings.TrimPrefix(p, "/"), "/")
			paths1 = append(paths1, p)
		}
	}
	return "/" + strings.Join(paths1, "/")
}

func JoinPathsDir(dir string, children ...string) []string {
	if len(children) == 0 {
		return nil
	}
	var r []string
	for _, child := range children {
		r = append(r, JoinPath(dir, child))
	}
	return r
}

func PermAnyone() []zk.ACL {
	return zk.WorldACL(zk.PermAll)
}
