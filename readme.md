# boatload
> oceonagraphic research data uploader

Tool for uploading oceanographic research data 
from research vessels to the `havvarsel-frost.met.no` API.

frontend is served at `/uploader`

## configuration

### deployment environment
Configure deployment enviroment through the environment variable `BOATLOAD_ENV`.
Acceptable values are `DEV` or `PROD`, and primarily define how to route outgoing traffic.

## deploy
This project runs on google app engine.

Deploy: `gcloud app deploy`