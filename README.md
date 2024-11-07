# hotels-data-merge
Simple Go application which cleans and merges hotel data from multiple sources.
## How to run app
### Prerequisites
- Docker
### Running the app
1. Build docker image
	`docker build -t hotel-data-merge .`
2. Run docker container
	`docker run -p 8080:8080 hotel-data-merge`
3. Application will be accessible on `http://localhost:8080`

## Exposed endpoints and filters
- `/hotels`
	- returns all hotels
- `/hotels?hotel_ids=iJhz,f8c9`
	- return hotels filtered by `hotel_ids` `iJhz` and `f8c9`. hotel ids are a list of comma separated ids 
- `/hotels?destination_ids=1122,5432`
	- return hotels filtered by `destination_ids` `1122` and `5432`. destination ids are a list of comma separated ids 
- if both `hotel_ids` and `destination_ids` are provided, `hotel_ids` will take precedence because the search is more specific

## Optimisations 
1. Caching of supplier endpoint responses using [gocache](https://github.com/eko/gocache).
2. Fetching of supplier hotel data parallelly using go routines

### Further optimisation considerations (not implemented)
1. Pagination can be implemented if the data size gets too big.
2. Depending on how frequently the supplier responses get updated, if it is a known frequency (eg. once every month), we can consider using a cronjob to pull and store the data in a database (eg. DynamoDB), and read from database instead.

## Testing pipeline
Tests are run on every PR create merging to `main`. Pipeline is executed using Github Actions. Example test pipeline [here](https://github.com/szeshen/hotels-data-merge/pull/2/checks). 

## Further Improvements
### Codebase
1. Adding config file 
	- Currently many static configurations such as cache expiry, and supplier hostname are stored as constants in the code. We can consider having a config file for it, particularly if there will be stg and prd environments with different configurations. 
	- This will also be better because any changes in config will only affect the config, and not code. 
2. Error responses and logging
	- Currently there is no proper format for returning errors implemented. Only success cases have been taking into consideration

### Data
1. Choosing of data
	- Description and hotel name were chosen by the supplier that provided the longest string for both fields. This is not the best way to select the description and hotel name.
	- We could look to implement a scoring system taking into account several factors such as sentiment and accuracy. It could be implemented using an external library or a separate service. 






