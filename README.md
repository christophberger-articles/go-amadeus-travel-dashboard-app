# Travel Dashboard Demo


## About this app

This app demonstrates the capabilities of the Amadeus services for building a travel dashboard for your customers.

## Prerequisites for running the app

[Sign up](https://developers.amadeus.com/register) for a free developer account and follow the steps in the [documentation](https://developers.amadeus.com/get-started/get-started-with-self-service-apis-335) to generate an API key and secret.

You need both in order to run the app (see below). The app needs the key and the secret to generate and refresh ephemeral access tokens. 

## How to run the app

You can run the demo app in a few easy steps.

1. Clone this repository.
2. `cd` into the repository.
3. Add the API key and secret as environment variables. For example, on Unix-like systems:
	```sh
	export AMADEUS_CLIENT_ID=...
	export AMADEUS_CLIENT_SECRET=...
	```

4. Execute `go run .`
5. Open the browser and navigate to http://localhost:8020.

## What can I do in this app?

This app displays travel information for a city of your choice.

1. Enter a city name and click the search button.
