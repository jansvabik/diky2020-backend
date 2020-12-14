# Díky 2020 (backend)
Backend for the Díky 2020 project by F. Soudek. Includes the one and only microservice (thanks API), build scripts and so on. Below is a guide to run this service on your server. It's really simple.

## Build and run in production mode
You can run the microservice in production mode by running this command:

    make run

## Run in development mode
To run the service (e.g. for testing purposes), just run this command:

    make up

There is still no difference between dev and prod mode, excluding the fact that in prod mode a binary file is built and ran.
