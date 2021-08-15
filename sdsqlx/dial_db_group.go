package sdsqlx

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdrand"
	"github.com/jmoiron/sqlx"
)

type GroupAddress struct {
	Driver     string   `json:"driver" toml:"driver"`
	MasterAddr string   `json:"master" toml:"master"`
	SlaveAddrs []string `json:"slaves" toml:"slaves"`
}

type DBGroup struct {
	Master *sqlx.DB
	Slaves []*sqlx.DB
}

func DialGroup(addr GroupAddress) (*DBGroup, error) {
	var errs []error
	var slaves []*sqlx.DB

	master, err := Dial(Address{
		Driver: addr.Driver,
		Addr:   addr.MasterAddr,
	})
	if err != nil {
		errs = append(errs, sderr.WithStack(err))
	}
	if len(addr.SlaveAddrs) > 0 {
		for _, slaveAddr := range addr.SlaveAddrs {
			slave, err := Dial(Address{
				Driver: addr.Driver,
				Addr:   slaveAddr,
			})
			if err != nil {
				errs = append(errs, sderr.WithStack(err))
			} else {
				slaves = append(slaves, slave)
			}
		}
	}
	dbg := DBGroup{
		Master: master,
		Slaves: slaves,
	}
	if err1 := sderr.Multi(errs); err1 != nil {
		_ = dbg.Close()
		return nil, err1
	}
	return &dbg, nil
}

func (dbg *DBGroup) Close() error {
	var errs []error
	if dbg.Master != nil {
		if err := dbg.Master.Close(); err != nil {
			errs = append(errs, sderr.WithStack(err))
		}
	}
	for _, slave := range dbg.Slaves {
		if slave != nil {
			if err := slave.Close(); err != nil {
				errs = append(errs, sderr.WithStack(err))
			}
		}
	}
	dbg.Master, dbg.Slaves = nil, nil
	return sderr.Multi(errs)
}

func (dbg *DBGroup) HasSlave() bool {
	return len(dbg.Slaves) > 0
}

func (dbg *DBGroup) W() *sqlx.DB {
	return dbg.Master
}

func (dbg *DBGroup) R() *sqlx.DB {
	if dbg.HasSlave() {
		return sdrand.ChoiceOne(dbg.Slaves).(*sqlx.DB)
	} else {
		return dbg.Master
	}
}

func (dbg *DBGroup) S(forWrite bool) *sqlx.DB {
	if forWrite {
		return dbg.W()
	} else {
		return dbg.R()
	}
}

func (dbg *DBGroup) RList() []*sqlx.DB {
	if dbg.HasSlave() {
		return dbg.Slaves
	} else {
		return []*sqlx.DB{dbg.Master}
	}
}

func (dbg *DBGroup) REach(f func(*sqlx.DB)) {
	for _, db := range dbg.RList() {
		f(db)
	}
}
