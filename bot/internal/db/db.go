package db

import "gopkg.in/mgo.v2"

type DbStruct struct {
	DbSession *mgo.Session
	PicCollection *mgo.Collection
}

var DB DbStruct

func init() {

}
