package main

import "gopkg.in/mgo.v2"

type DbStruct struct {
	DbSession       *mgo.Session
	UsersCollection *mgo.Collection
	AdminCollection *mgo.Collection
}

type obj map[string]interface{}

func connectToDB() *mgo.Session {
	d, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return d
}

func connectToUsersCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("users")
}

func connectToAdminCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("admin")
}

func (d *DbStruct) getUserFromDbById(id int) (result User, err error) {
	err = DB.UsersCollection.Find(obj{"id": id}).One(&result)
	return
}

func (d *DbStruct) userExists(id int) (bool, error){
	var u User
	if err := d.UsersCollection.Find(obj{"telebot.id":id}).One(&u);err != nil{
		if err.Error() == "not found"{
			return false, nil
		}
		return false, err
	}
	return true, nil
}
