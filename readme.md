# GO-TAMBOON ไปทำบุญ

## Project Structure
```
    ├──cmd/
    |  ├──make_permit_runner/
    |  |  └──main.go
    └──config/
    └──internal/
    |  └──data/
    |  └──di/
    |  └──model/
    |  └──runner/
    |  └──service/
    |  |  └──donation/
    |  |  └──payment/
    |  └──validator/
    └──pkg/
    |  └──cipher/
    |  └──omise/
    └──testutils/
```
### cmd
A main directory for the application for entrypoint. I've created the sub directory inside and named it `make_permit_runner`
### config
A directory for configuration for the application
### internal
This directory is for all the packages that must not be shared into public or for other go packages or applications. There's five sub diractory inside. `data, di, model, runner` and `service`
1. `data` for storing encrypted donation data file
2. `di`  for codegen dependency injection of this project
3. `model` for domain model of this project
4. `runner` for application runner
5. `service` for all of service classes in this project
6. `validator` for validate struct inside this project.
### pkg
A directory for collecting the packages that can be shared into public or other go packages such as `cipher, omise`.
### testutils
A directory for collecting all of test util functions.

---

## Run application locally
1. Sign up and register Omise account on official website to get public key and private key  [Opn](https://sso-idp.omise.co/realms/engagement/protocol/openid-connect/auth?client_id=dashboard&redirect_uri=https%3A%2F%2Fdashboard.omise.co%2Fv2&state=60db43d3-75d3-4180-88c7-289130bf101a&response_mode=fragment&response_type=code&scope=openid&nonce=9d4f30e1-5b1c-4ba5-aab5-a350d3c3b1b1&code_challenge=Bv1ASBnaaADBRtps0WQVETUdD_dPjdcAozT6DYgYBnk&code_challenge_method=S256)
2. Set up environment varriable in `.env` file in root project (replace the public key and private key to `APP_OMISE_PUBLIC_KEY` and  `APP_OMISE_PRIVATE_KEY`)
```
APP_DONATION_FILE_ADDR=internal/data/fng.1000.csv.rot128
APP_MAX_CONCURRENT=8
APP_OMISE_PUBLIC_KEY=pkey_test_no1t4tnemucod0e51mo
APP_OMISE_PRIVATE_KEY=skey_test_no1t4tnemucod0e51mo
```
3. Run command `make run`. The script will be generating all of codegen and start the application.

## Unit test
```
make test
```
