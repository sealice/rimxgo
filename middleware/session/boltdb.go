package session

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/boltdb"
)

func useBoltdb(sess *sessions.Sessions) {
	db, _ := boltdb.New("./sessions/store.db", 0666)

	iris.RegisterOnInterrupt(func() {
		db.Close()
	})

	sess.UseDatabase(db)
}
