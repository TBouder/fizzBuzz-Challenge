/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 11 February 2020 - 12:26:25
** @Filename:				main.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Friday 14 February 2020 - 13:42:22
*******************************************************************************/

package			main

import			"fmt"
import			"context"
import			"testing"
import			"bytes"
import			"strings"
import			"net"
import			"net/http"
import			"io/ioutil"
import			"github.com/valyala/fasthttp"
import			"github.com/valyala/fasthttp/fasthttputil"

type	sTestFastHttpHandler struct {
	method	string
	uri		string
	body	[]byte
	result	string
}

/******************************************************************************
**	Serve a dummy server to test a specific route
******************************************************************************/
func	serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

/******************************************************************************
**	Call the Rest API server to test a specific route
******************************************************************************/
func	performRequest(t *testing.T, each sTestFastHttpHandler) {
	r, err := http.NewRequest(each.method, each.uri, bytes.NewReader(each.body))
	if err != nil {
		t.Error(err)
	}

	res, err := serve(initRouter(), r)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	result := strings.TrimSuffix(string(body), "\n")

	if (string(result) != each.result) {
		t.Errorf("FAIL : bad [%s] - [%s] route result, got: [%s] instead of [%s].", each.method, each.uri, result, each.result)
	}
}

/******************************************************************************
**	Test the FastHTTP server and the global program by performing
******************************************************************************/
func	TestFastHttpHandler(t *testing.T) {
	DB = initDB(`./badgerTest`)
	defer DB.Close()
	defer DB.DropAll()

	/**************************************************************************
	**	Check challengeStats when no Key
	**************************************************************************/
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques/`, nil, `"Key not found"`})
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques`, nil, `"Key not found"`})
	
	/**************************************************************************
	**	Check challengeStats with invalid request
	**************************************************************************/	
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/getChallengeStatistiques/`, nil, `Method Not Allowed`})
	performRequest(t, sTestFastHttpHandler{`DELETE`, `http://localhost:8000/getChallengeStatistiques/`, nil, `Method Not Allowed`})
	performRequest(t, sTestFastHttpHandler{`PUT`, `http://localhost:8000/getChallengeStatistiques/`, nil, `Method Not Allowed`})

	/**************************************************************************
	**	Perfom a successfull challenge request
	**************************************************************************/	
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fizz","str2":"buzz"}`), `["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16"]`})

	/**************************************************************************
	**	Perfom a bad challenge request 
	**************************************************************************/	
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":3,"limit":16,"str1":"fizz","str2":"buzz"}`), `"int1 is the same as int2 -- aborting"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"","str2":"buzz"}`), `"str1 is empty"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str2":"buzz"}`), `"str1 is empty"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fizz","str2":""}`), `"str2 is empty"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fizz"}`), `"str2 is empty"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":0,"str1":"fizz","str2":"buzz"}`), `"limit is not set"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"str1":"fizz","str2":"buzz"}`), `"limit is not set"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":-1,"str1":"fizz","str2":"buzz"}`), `"limit needs to be positive"`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"`), `"unexpected end of JSON input"`})

	/**************************************************************************
	**	Check the challenge stats
	**************************************************************************/	
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques/`, nil, `{"hits":1,"int1":3,"int2":5,"limit":16,"str1":"fizz","str2":"buzz"}`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"}`), `["1","2","fazz","4","buzz","fazz","7","8","fazz","buzz","11","fazz","13","14","fazzbuzz","16"]`})
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"}`), `["1","2","fazz","4","buzz","fazz","7","8","fazz","buzz","11","fazz","13","14","fazzbuzz","16"]`})
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques/`, nil, `{"hits":2,"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"}`})
	
	/**************************************************************************
	**	Check challengeStats when mostKey is an invalid JSON
	**************************************************************************/
	DB.DropAll()
	setCurrentKey([]byte(`mostHit`), []byte(`10`))
	setCurrentKey([]byte(`mostKey`), []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"`))
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques/`, nil, `"unexpected end of JSON input"`})

	/**************************************************************************
	**	Check challengeStats when mostHit is an invalid int
	**************************************************************************/
	DB.DropAll()
	setCurrentKey([]byte(`mostHit`), []byte(`notAnInt`))
	setCurrentKey([]byte(`mostKey`), []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"}`))
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques/`, nil, `"strconv.Atoi: parsing \"notAnInt\": invalid syntax"`})

	/**************************************************************************
	**	Check challengeStats when mostKey is not Set
	**************************************************************************/
	DB.DropAll()
	setCurrentKey([]byte(`mostHit`), []byte(`notAnInt`))
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques/`, nil, `"Key not found"`})

	/**************************************************************************
	**	Check challengeStats when mostHit is not Set
	**************************************************************************/
	DB.DropAll()
	setCurrentKey([]byte(`mostKey`), []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"}`))
	performRequest(t, sTestFastHttpHandler{`GET`, `http://localhost:8000/getChallengeStatistiques/`, nil, `"Key not found"`})

	/**************************************************************************
	**	Check the successfull process of the performChallenge even if
	**	the saveStats fails because of mostHit atoi failure
	**************************************************************************/
	DB.DropAll()
	setCurrentKey([]byte(`mostHit`), []byte(`notAnInt`))
	setCurrentKey([]byte(`mostKey`), []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"}`))
	performRequest(t, sTestFastHttpHandler{`POST`, `http://localhost:8000/performChallenge/`, []byte(`{"int1":3,"int2":5,"limit":16,"str1":"fazz","str2":"buzz"}`), `["1","2","fazz","4","buzz","fazz","7","8","fazz","buzz","11","fazz","13","14","fazzbuzz","16"]`})

}
