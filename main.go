/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 11 February 2020 - 10:26:25
** @Filename:				main.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Friday 14 February 2020 - 13:08:28
*******************************************************************************/

package			main

import			"log"
import			"encoding/json"
import			"github.com/microgolang/logs"
import			"github.com/valyala/fasthttp"
import			"github.com/buaazp/fasthttprouter"
import			"github.com/dgraph-io/badger"

var DB *badger.DB

/******************************************************************************
**	To get the statistiques and get the number of hits per parameters in
**	order to find the most used parameters, we need to store theses
**	informations somewhere.
**	It could be cache, with an array of struct, but the data will not be
**	persistent.
**	In order to get persistent data, we are using BadgerDB, which is easy to
**	setup and effective for this type of challenge.
**	=> https://github.com/dgraph-io/badger
**	
**	The DB instance is opened on init.
******************************************************************************/
func	initDB(path string) (*badger.DB) {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	db, err := badger.Open(options)
	if (err != nil) {
		log.Fatal(err)
	}
	return db
}

func	resolveError(ctx *fasthttp.RequestCtx, err error) {
	ctx.Response.Header.SetContentType(`application/json`)
	ctx.Response.SetStatusCode(400)
	json.NewEncoder(ctx).Encode(err.Error())
}

func	initRouter() func(*fasthttp.RequestCtx) {
	router := fasthttprouter.New()
	router.POST("/performChallenge/", performChallengeHandler)
	router.GET("/getChallengeStatistiques/", getChallengeStatistiquesHandler)
	return router.Handler
}


func	main() {
	DB = initDB(`./badger`)
	defer DB.Close()
	logs.Success(`Listening on :8000`)
	log.Fatal(fasthttp.ListenAndServe(`:8000`, initRouter()))
}
