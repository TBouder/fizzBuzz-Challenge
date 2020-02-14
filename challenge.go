/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Tuesday 11 February 2020 - 11:30:49
** @Filename:				challenge.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Friday 14 February 2020 - 14:21:41
*******************************************************************************/

package			main

import			"errors"
import			"strconv"
import			"github.com/valyala/fasthttp"
import			"encoding/json"

/******************************************************************************
**	Function to check if a value is a multiple of another one.
**	0 is a multiple for every values.
******************************************************************************/
func	IsMultiple(value, multiple int) bool {
	if (multiple == 0) {
		return true
	}
	return value % multiple == 0
}

/******************************************************************************
**	Structure containing the arguments for the challenge.
**	Int1 is an integer and will be replaced by Str1
**	Int2 is an integer and will be replaced by Str2
**	Limit is a non-empty positive integer, otherwise an error is returned
**	Str1 is a non-empty string, otherwise an error is returned
**	Str2 is a non-empty string, otherwise an error is returned
**	All others arguments are ignored.
******************************************************************************/
type	sPerformChallenge struct {
	Int1	int		`json:"int1"`
	Int2	int		`json:"int2"`
	Limit	int		`json:"limit"`
	Str1	string	`json:"str1"`
	Str2	string	`json:"str2"`
}

/******************************************************************************
**	Argument checker for the performChallenge handler
**	Will check different case of error and returns it.
**	Returns nil if there is no error
******************************************************************************/
func	performChallengeCheckArguments(body *sPerformChallenge) error {
	if (body.Int1 == body.Int2) {
		return errors.New(`int1 is the same as int2 -- aborting`)
	} else if (body.Str1 == ``) {
		return errors.New(`str1 is empty`)
	} else if (body.Str2 == ``) {
		return errors.New(`str2 is empty`)
	} else if (body.Limit == 0) {
		return errors.New(`limit is not set`)
	} else if (body.Limit < 0) {
		return errors.New(`limit needs to be positive`)
	}
	return nil
}

/******************************************************************************
**	Get the string to append to the result for the given index and parameters
**	informations
******************************************************************************/
func	getValue(index, int1, int2 int, str1, str2 string) string {
	result := ``
	if (IsMultiple(index, int1)) {
		result += str1
	}
	if (IsMultiple(index, int2)) {
		result += str2
	}
	if (result == ``) {
		result = strconv.Itoa(index)
	}
	return result
}
/******************************************************************************
**	From 1 to limit :
**	- will replace i by str1 when i is a multiple of int1
**	- will replace i by str2 when i is a multiple of int2
**	- will replace i by str1str2 when i is a multiple of int1 and int2
**	- will not replace i if the above conditions do not match
******************************************************************************/
func	performChallenge(int1, int2, limit int, str1, str2 string) []string {
	results := []string{}
	for i := 1; i <= limit; i++ {
		results = append(results, getValue(i, int1, int2, str1, str2))
	}
	return results
}

/******************************************************************************
**	Router handler to perform the Fizz-Buzz challenge
******************************************************************************/
func	performChallengeHandler(ctx *fasthttp.RequestCtx) {
	body := &sPerformChallenge{}
	
	err := json.Unmarshal(ctx.PostBody(), &body)
	if (err != nil) {
		resolveError(ctx, err)
		return
	}

	err = performChallengeCheckArguments(body)
	if (err != nil) {
		resolveError(ctx, err)
		return
	}

	saveStats(ctx.PostBody())

	results := performChallenge(body.Int1, body.Int2, body.Limit, body.Str1, body.Str2)
	ctx.Response.Header.SetContentType(`application/json`)
	ctx.Response.SetStatusCode(200)
	json.NewEncoder(ctx).Encode(results)
}
