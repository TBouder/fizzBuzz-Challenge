/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 11 February 2020 - 11:30:49
** @Filename:				challenge.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Tuesday 11 February 2020 - 14:46:38
*******************************************************************************/

package			main

import			"log"
import			"sync"
import			"strconv"
import			"encoding/json"
import			"github.com/valyala/fasthttp"
import			"github.com/dgraph-io/badger"

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
var DB *badger.DB
func	init() {
	db, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}

/******************************************************************************
**	In order to avoid data not being update with a lot of simultanous update,
**	we will need to lock and unlock the database write operation
******************************************************************************/
var	RWMut sync.RWMutex

/******************************************************************************
**	Structure containing the informations about the most used request
**	Request containing the request as []byte
**	Hits is an interger containing the number of hit for the Request
******************************************************************************/
type	sStatsChallengeDB struct {
	Request	[]byte
	Hits	int		`json:"hits"`
}
/******************************************************************************
**	Structure containing the response value for the getChallengeStatistiques
**	endpoint. It's the same as for performChallenge, but with the number of
**	hits added to it.
******************************************************************************/
type	sStatsChallenge struct {
	Hits	int		`json:"hits"`
	Int1	int		`json:"int1"`
	Int2	int		`json:"int2"`
	Limit	int		`json:"limit"`
	Str1	string	`json:"str1"`
	Str2	string	`json:"str2"`
}


/******************************************************************************
**	Helper to get a specific Key from the DB
******************************************************************************/
func	getValueFromKey(bodyAsKey []byte) ([]byte, error) {
	toReturn := []byte{}

	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(bodyAsKey)
		if (err != nil) {
			return err
		}

		var value []byte
		item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			toReturn = value
			return nil
		})
		return nil
	})
	return toReturn, err
}

/******************************************************************************
**	Helper to set a Key to a specific Value in the DB
******************************************************************************/
func	setCurrentKey(bodyAsKey []byte, currentValue []byte) error {
	err := DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(bodyAsKey, currentValue)
		return err
	})
	return err
}

/******************************************************************************
**	To get the most used parameters, we can access the mostHit and mostKey
**	special elements from the DB, containing this information, and returns it
**	as a sStatsChallengeDB struct
******************************************************************************/
func	checkMostUsed() (sStatsChallengeDB, error) {
	mostUsedRequest := sStatsChallengeDB{}

	hits, err := getValueFromKey([]byte(`mostHit`))
	if (err != nil) {
		return mostUsedRequest, err
	}

	request, err := getValueFromKey([]byte(`mostKey`))
	if (err != nil) {
		return mostUsedRequest, err
	}


	hitsInt, err := strconv.Atoi(string(hits))
	if (err != nil) {
		return mostUsedRequest, err
	}
	
	mostUsedRequest.Hits = hitsInt
	mostUsedRequest.Request = request

	return mostUsedRequest, nil
}

/******************************************************************************
**	When a request is made and the parameters are valid, we are storing them
**	in the badger database.
**	We first need to lock the mutex for write access, to avoid overlap
**	Then, we need to get the mostHit hit counts and the current parameters hit
**	count
**	Then, we check if the current parameters hit count +1 (this hit) is above
**	the previous mostHit hit count. In this case, we can update the `mostHit`
**	and `mostKey` values with the current parameters information and hit count
**	Then, we update in the DB the current parameters hit count 
******************************************************************************/
func	saveStats(body []byte) error {
	RWMut.Lock()
	defer RWMut.Unlock()

	currentMostUsedByte, errMostUsed := getValueFromKey([]byte(`mostHit`))
	if (errMostUsed != nil && errMostUsed.Error() != `Key not found`) {
		return errMostUsed
	} else if (errMostUsed != nil && errMostUsed.Error() == `Key not found`) {
		currentMostUsedByte = []byte(`0`)
	}

	currentUsedByte, errUsed := getValueFromKey(body)
	if (errUsed != nil && errUsed.Error() != `Key not found`) {
		return errUsed
	} else if (errUsed != nil && errUsed.Error() == `Key not found`) {
		currentUsedByte = []byte(`0`)
	}

	currentMostUsedInt, err := strconv.Atoi(string(currentMostUsedByte))
	if (err != nil) {
		return err
	}

	currentUsedInt, err := strconv.Atoi(string(currentUsedByte))
	if (err != nil) {
		return err
	}

	if (currentUsedInt + 1 > currentMostUsedInt) {
		setCurrentKey([]byte(`mostHit`), []byte(strconv.Itoa(currentUsedInt + 1)))
		setCurrentKey([]byte(`mostKey`), body)
	}
	setCurrentKey(body, []byte(strconv.Itoa(currentUsedInt + 1)))

	return nil
}

/******************************************************************************
**	Router handler to perform the Fizz-Buzz challenge statistiques
******************************************************************************/
func	getChallengeStatistiquesHandler(ctx *fasthttp.RequestCtx) {
	response := sStatsChallenge{}
	mostUsed, err := checkMostUsed()
	if (err != nil) {
		resolveError(ctx, err)
		return
	}

	err = json.Unmarshal(mostUsed.Request, &response)
	if (err != nil) {
		resolveError(ctx, err)
		return
	}

	response.Hits = mostUsed.Hits

	ctx.Response.Header.SetContentType(`application/json`)
	ctx.Response.SetStatusCode(200)
	json.NewEncoder(ctx).Encode(response)
}