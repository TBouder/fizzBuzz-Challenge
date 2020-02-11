/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 11 February 2020 - 10:26:25
** @Filename:				main.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Tuesday 11 February 2020 - 12:44:52
*******************************************************************************/

package			main

import			"log"
import			"github.com/microgolang/logs"
import			"github.com/valyala/fasthttp"
import			"github.com/lab259/cors"
import			"github.com/buaazp/fasthttprouter"

func	initRouter() func(*fasthttp.RequestCtx) {
	router := fasthttprouter.New()
	router.POST("/performChallenge/", performChallengeHandler)
	router.GET("/getChallengeStatistiques/", getChallengeStatistiquesHandler)
	return router.Handler
}

func	expose() {
	defer DB.Close()
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods: []string{`GET`, `POST`, `DELETE`, `PUT`, `OPTIONS`, `OPTION`},
		AllowedHeaders:	[]string{
			`Access-Control-Allow-Origin`,
			`Access-Control-Allow-Credentials`,
			`Content-Type`,
			`Transfer-Encoding`,
		},
		AllowCredentials: true,
	})

	handler := c.Handler(initRouter())
	logs.Success(`Listening on :8000`)
	log.Fatal(fasthttp.ListenAndServe(`:8000`, handler))
}

func	main()	{
	expose()
}