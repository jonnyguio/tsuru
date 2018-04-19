// Copyright 2018 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tsuru/tsuru/db"
	"github.com/tsuru/tsuru/types/app"
)

var _ app.AppQuotaStorage = &AppQuotaStorage{}

type AppQuotaStorage struct{}

type _app struct {
	appName string `bson:"_id"`
	quota   *app.AppQuota
}

func (s *AppQuotaStorage) IncInUse(quota *app.AppQuota, quantity int) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Apps().Update(
		bson.M{"name": quota.AppName, "quota.inuse": quota.InUse},
		bson.M{"$inc": bson.M{"quota.inuse": quantity}},
	)
	return err
}

func (s *AppQuotaStorage) SetLimit(appName string, limit int) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Apps().Update(
		bson.M{"name": appName},
		bson.M{"$set": bson.M{"quota.limit": limit}},
	)
	return err
}

func (s *AppQuotaStorage) SetInUse(appName string, inUse int) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Apps().Update(
		bson.M{"name": appName},
		bson.M{"$set": bson.M{"quota.inuse": inUse}},
	)
	return err
}

func (s *AppQuotaStorage) FindByAppName(appName string) (*app.AppQuota, error) {
	var a _app
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = conn.Apps().FindId(appName).One(&a)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, app.ErrAppNotFound
		}
		return nil, err
	}
	return a.quota, nil
}
