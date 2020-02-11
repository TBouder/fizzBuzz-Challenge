## FizzBuzz Challenge

> *The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz". The output would look like this:
> `1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...`.
>
> Your goal is to implement a web server that will expose a REST API endpoint
> ##### REST API :
> - Accepts five parameters : three integers int1, int2 and limit, and two strings str1 and str2.
> - Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.
>
> #### The server needs to be:
> - Ready for production
> - Easy to maintain by other developers
>
> #### Bonus question :
> - Add a statistics endpoint allowing users to know what the most frequent request has been. This endpoint should:
> - Accept no parameter
> - Return the parameters corresponding to the most used request, as well as the number of hits for this request

## Setup
The program will expose a Rest API which will listen to the port `:8000`.
**With Docker**
1. `cd /tmp && git clone https://github.com/TBouder/fizzBuzz-Challenge.git`
2. `cd fizzBuzz-Challenge`
3. `docker-compose up --build`

**Without Docker**
1. Make sure you have `go` 
2. `cd /tmp && git clone https://github.com/TBouder/fizzBuzz-Challenge.git`
3. `cd fizzBuzz-Challenge`
4. `go build -o exec`
5. `./exec`


## Usage
> ##### getChallengeStatistiques
> Get the most used request and the number of hits for this request. Request must be a GET request. 
`curl -X GET localhost:8000/getChallengeStatistiques/`
>
>Reponse example : `{"hits":1,"int1":3,"int2":5,"limit":10,"str1":"fizz","str2":"buzz"}`

> ##### performChallenge
> Perfom the challenge with the given parameters. Parameters should be in JSON and as a POST request
`curl -d '{"int1":3,"int2":5,"limit":16,"str1":"fizz","str2":"buzz"}' -H "Content-Type: application/json" -X POST localhost:8000/performChallenge/`
>
>Reponse example : `["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16"]`
>
>**Note** :
> * int1 and int2 must be different
> * if int1 is not set, int1 is set to 0
> * if int2 is not set, int2 is set to 0
> * str1 and str2 must not be empty
> * limit must not be 0
> * limit must not be negative
