# Webmonitor

This is a website monitoring service that I wrote to monitor various
websites.

I recently rewritten it completely as an excuse to learn about some web
technologies.

## Backend
The backend uses a Monitor object that runs at predefined intervals, fetches the checks for that interval and runs them.

The checks are stored in a PostgreSQL database and the queries are done the the standard library with the [pq](https://pkg.go.dev/github.com/lib/pq@v1.9.0) driver, no ORM whatsoever.

If a difference is detected, the user is alerted with an email using [Sendgrid](https://sendgrid.com/) and saves the body of the web page.

The interaction with the frontend is done via a simple CRUD API using [Gorilla Mux](https://github.com/gorilla/mux).

There is no authorization or authentication at the moment, as this is something that is thought as selfhosted at home, but I might add it later on.


## Frontend
The frontend is a Typescript [React](https://reactjs.org/) App using [Chakra](https://chakra-ui.com/) for the user interface.

The interaction with the backend is done via [Axios](https://github.com/axios/axios) which uses interceptors to deserialize things like dates.

The state and cache handling is done via [React query](https://react-query.tanstack.com/).

## TODO
* [ ] Implement multiple notification services
* [ ] Make notification service per check
* [ ] Add a token to delete a check
* [x] Add persistent storage
* [x] Make API handler use storage instead of monitor
* [x] Add validation to API
* [ ] Use sensible defaults instead of erroring out
