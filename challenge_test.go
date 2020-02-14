/*******************************************************************************
** @Author:					Thomas Bouder <Tbouder>
** @Email:					Tbouder@protonmail.com
** @Date:					Friday 14 February 2020 - 12:44:56
** @Filename:				challenge.test.go
**
** @Last modified by:		Tbouder
** @Last modified time:		Friday 14 February 2020 - 13:44:26
*******************************************************************************/

package main

import (
	"reflect"
	"testing"
)

type	sTestIsMultiple struct {
	value		int
	multiple	int
	result		bool
}

type	sTestPerformChallenge struct {
	int1	int
	int2	int
	limit	int
	str1	string
	str2	string
	result	[]string
}

/******************************************************************************
**	More test to check if the IsMultiple functions works correctly.
******************************************************************************/
func	TestIsMultiple(t *testing.T) {
	var	toTest = []sTestIsMultiple{}

	/**************************************************************************
	**	A test with 0 as multiple is always positive
	**************************************************************************/
	toTest = append(toTest, []sTestIsMultiple{
		{0, 0, true},
		{1, 0, true},
		{42, 0, true},
	}...)

	/**************************************************************************
	**	A test with 1 as multiple is always positive
	**************************************************************************/
	toTest = append(toTest, []sTestIsMultiple{
		{1, 1, true},
		{1, 1, true},
		{42, 1, true},
	}...)

	/**************************************************************************
	**	Positive and negative values/multiple should give the same result
	**************************************************************************/
	toTest = append(toTest, []sTestIsMultiple{
		{2, 3, false},
		{6, 3, true},
		{15, 3, true},
		{17, 3, false},
		{2, -3, false},
		{6, -3, true},
		{15, -3, true},
		{17, -3, false},
		{-2, 3, false},
		{-6, 3, true},
		{-15, 3, true},
		{-17, 3, false},
		{-2, -3, false},
		{-6, -3, true},
		{-15, -3, true},
		{-17, -3, false},
	}...)

	/**************************************************************************
	**	Dummies values to test
	**************************************************************************/
	toTest = append(toTest, []sTestIsMultiple{
		{2, 2, true},
		{2, 3, false},
		{10, 3, false},
		{15, 3, true},
		{17, 3, false},
		{8, 3, false},
		{9, 3, true},
		{154, 3, false},
		{24, 3, true},
		{100, 3, false},
		{44, 3, false},
		{10, 2, true},
		{10, 5, true},
		{15, 2, false},
		{15, 5, true},
		{17, 5, false},
		{8, 5, false},
		{9, 5, false},
		{154, 2, true},
		{24, 2, true},
		{100, 5, true},
		{44, 5, false},
	}...)

	for _, each := range toTest {
		result := IsMultiple(each.value, each.multiple)
		if result != each.result {
			t.Errorf("FAIL : IsMultiple of (%d, %d), got: [%t] instead of [%t].", each.value, each.multiple, result, each.result)
		}
	}
}

/******************************************************************************
**	Helper for the TestPerformChallenge testing function to actually test
**	the PerformChallenge function and check if the result is correct
******************************************************************************/
func	testChallenge(t *testing.T, each sTestPerformChallenge) {
	result := PerformChallenge(each.int1, each.int2, each.limit, each.str1, each.str2)
	if (!reflect.DeepEqual(result, each.result)) {
		t.Errorf("FAIL : PerformChallenge with (%d, %d, %d, %s, %s), got: [%s] instead of [%s].", each.int1, each.int2, each.limit, each.str1, each.str2, result, each.result)
	}
}

/******************************************************************************
**	TestPerformChallenge defines a list a parameters to test if the fuzz-bizz
**	challenge is correct
******************************************************************************/
func	TestPerformChallenge(t *testing.T) {
	/**************************************************************************
	**	General test from 1 to 100 with
	**	mult of 3 = fizz
	**	mult of 5 = buzz
	**	mult of both = fizzbuzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{3, 5, 100, `fizz`, `buzz`, []string{`1`, `2`, `fizz`, `4`, `buzz`, `fizz`, `7`, `8`, `fizz`, `buzz`, `11`, `fizz`, `13`, `14`, `fizzbuzz`, `16`, `17`, `fizz`, `19`, `buzz`, `fizz`, `22`, `23`, `fizz`, `buzz`, `26`, `fizz`, `28`, `29`, `fizzbuzz`, `31`, `32`, `fizz`, `34`, `buzz`, `fizz`, `37`, `38`, `fizz`, `buzz`, `41`, `fizz`, `43`, `44`, `fizzbuzz`, `46`, `47`, `fizz`, `49`, `buzz`, `fizz`, `52`, `53`, `fizz`, `buzz`, `56`, `fizz`, `58`, `59`, `fizzbuzz`, `61`, `62`, `fizz`, `64`, `buzz`, `fizz`, `67`, `68`, `fizz`, `buzz`, `71`, `fizz`, `73`, `74`, `fizzbuzz`, `76`, `77`, `fizz`, `79`, `buzz`, `fizz`, `82`, `83`, `fizz`, `buzz`, `86`, `fizz`, `88`, `89`, `fizzbuzz`, `91`, `92`, `fizz`, `94`, `buzz`, `fizz`, `97`, `98`, `fizz`, `buzz`}})

	/**************************************************************************
	**	String test from 1 to 100 with
	**	mult of 3 = fAzz
	**	mult of 5 = bAzz
	**	mult of both = fAzzbAzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{3, 5, 100, `fAzz`, `bAzz`, []string{`1`, `2`, `fAzz`, `4`, `bAzz`, `fAzz`, `7`, `8`, `fAzz`, `bAzz`, `11`, `fAzz`, `13`, `14`, `fAzzbAzz`, `16`, `17`, `fAzz`, `19`, `bAzz`, `fAzz`, `22`, `23`, `fAzz`, `bAzz`, `26`, `fAzz`, `28`, `29`, `fAzzbAzz`, `31`, `32`, `fAzz`, `34`, `bAzz`, `fAzz`, `37`, `38`, `fAzz`, `bAzz`, `41`, `fAzz`, `43`, `44`, `fAzzbAzz`, `46`, `47`, `fAzz`, `49`, `bAzz`, `fAzz`, `52`, `53`, `fAzz`, `bAzz`, `56`, `fAzz`, `58`, `59`, `fAzzbAzz`, `61`, `62`, `fAzz`, `64`, `bAzz`, `fAzz`, `67`, `68`, `fAzz`, `bAzz`, `71`, `fAzz`, `73`, `74`, `fAzzbAzz`, `76`, `77`, `fAzz`, `79`, `bAzz`, `fAzz`, `82`, `83`, `fAzz`, `bAzz`, `86`, `fAzz`, `88`, `89`, `fAzzbAzz`, `91`, `92`, `fAzz`, `94`, `bAzz`, `fAzz`, `97`, `98`, `fAzz`, `bAzz`}})

	/**************************************************************************
	**	0 test from 1 to 20 with
	**	mult of 0 = fizz
	**	mult of 0 = buzz
	**	mult of both = fizzbuzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{0, 0, 20, `fizz`, `buzz`, []string{`fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`}})

	/**************************************************************************
	**	1 test from 1 to 20 with
	**	mult of 1 = fizz
	**	mult of 1 = buzz
	**	mult of both = fizzbuzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{1, 1, 20, `fizz`, `buzz`, []string{`fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`, `fizzbuzz`}})

	/**************************************************************************
	**	Positive values test from 1 to 20 with
	**	mult of 2 = fizz
	**	mult of 4 = buzz
	**	mult of both = fizzbuzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{2, 4, 20, `fizz`, `buzz`, []string{`1`, `fizz`, `3`, `fizzbuzz`, `5`, `fizz`, `7`, `fizzbuzz`, `9`, `fizz`, `11`, `fizzbuzz`, `13`, `fizz`, `15`, `fizzbuzz`, `17`, `fizz`, `19`, `fizzbuzz`}})

	/**************************************************************************
	**	Positive/Negative values test from 1 to 20 with
	**	mult of -2 = fizz
	**	mult of 4 = buzz
	**	mult of both = fizzbuzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{-2, 4, 20, `fizz`, `buzz`, []string{`1`, `fizz`, `3`, `fizzbuzz`, `5`, `fizz`, `7`, `fizzbuzz`, `9`, `fizz`, `11`, `fizzbuzz`, `13`, `fizz`, `15`, `fizzbuzz`, `17`, `fizz`, `19`, `fizzbuzz`}})

	/**************************************************************************
	**	Negative values test from 1 to 20 with
	**	mult of -2 = fizz
	**	mult of -4 = buzz
	**	mult of both = fizzbuzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{-2, -4, 20, `fizz`, `buzz`, []string{`1`, `fizz`, `3`, `fizzbuzz`, `5`, `fizz`, `7`, `fizzbuzz`, `9`, `fizz`, `11`, `fizzbuzz`, `13`, `fizz`, `15`, `fizzbuzz`, `17`, `fizz`, `19`, `fizzbuzz`}})

	/**************************************************************************
	**	Negative limit test from 1 to 20 with
	**	mult of -2 = fizz
	**	mult of -4 = buzz
	**	mult of both = fizzbuzz
	**************************************************************************/
	testChallenge(t, sTestPerformChallenge{-2, -4, -20, `fizz`, `buzz`, []string{}})
}