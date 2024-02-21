## Study Case
Create Money Transfer APIs using
Golang programming consist of 3 endpoints, namely:
1. Account Validation URL endpoint that used to validate account number and
   name of the bank account owner. To validate the account, please create a
   mock endpoint using https://mockapi.io/ or other preferred service that will act
   as a bank.
2. Transfer/Disbursement URL endpoint that used to transfer money to the
   destination account. To transfer the money, you can create a mock endpoint
   similar to point 1 that acts as a bank.
3. Transfer/Disbursement Callback URL endpoint that used to receive callback
   status of transfer from the bank.

### Functional Requirements
#### Disbursement Endpoint
- The system must be able to disburse money to the bank account
- The system must be able to handle duplicate reference
- The system must be able to track transaction status
- The system must be able to track bank transaction reference
- The system must be able to verify bank account number
- The system must be able to store information about bank account used in the transaction

#### Disbursement Callback Endpoint
- The system must be able to receive callback
- The system must be able to handle relative big traffic
- The system must provide a way for reconciliation process
- The system must be able to handle failure properly and restore user balance
- The system must be able to handle duplicated event / notification

#### Bank Validation Endpoint
- The system must be able to validate bank account name by calling external API

### Non Functional Requirements
- The system can be scaled up when the number of transactions is significantly increased

### Key Assumptions
- To create disbursement we call external service (BANK or 3rd Party Payment Service Provider)
- The service run in the internal network which can't be accessed publicly. And the authentication handled by API Gateway and Authentication service

### Out Of Scope
- DLQ Handling and Retry Mechanism
- Transaction Reconciliation
- Unit Test and Integration Test (time constrain)
- Customer notification
- Authentication/Authorization

## Implementation Detail
### Sequence Diagram
#### Disbursements
![img.png](img.png)
#### Bank Account Validation
![img_2.png](img_2.png)
### Database Schema
![img_1.png](img_1.png)

### Mock Service
To mock external service (in this case bank service) I use mockapi.io with the following details 
#### Account Validation
```json 
 {
   "url": "https://65d38c2f522627d501091ff2.mockapi.io/api/bank-accounts/:account_number",
   "method": "GET"
 }
```
example request:
```curl
   curl --location 'https://65d38c2f522627d501091ff2.mockapi.io/api/bank-accounts/511110'
```
example response:
```json
   {"account_number":"511110","account_name":"Ben Ward"}
```
#### Transfer
```json
   {
   "url": "https://65d38c2f522627d501091ff2.mockapi.io/api/transfers",
   "method": "POST"
   }
```
example request:
```curl
   curl --location --request POST 'https://65d38c2f522627d501091ff2.mockapi.io/api/transfers' \
   --data ''
```
example response:
```json
  {
   "transaction_id": "1"
  }
```
## Run Application in Local
### Prerequisite
1. go version ^1.21.6
2. docker-compose https://docs.docker.com/compose/install/
2. golang migrate cli https://pkg.go.dev/github.com/golang-migrate/migrate/cli#readme-installation

### Build The Project
1. run dependency (postgres, sqs) 
```$ docker-compose up -d```
2. run database migration for the first time
```$ migrate -path migrations -database  postgres://brick:brick@localhost:5432/brick?sslmode=disable up```
3. build project ```go build -o main ./app``` (If you are using windows, put .exe extension in the output binary name)
4. run web server ```./main start-rest-server```
5. run worker ```./main start-worker```

### Documentations
[Read Documentation](./docs)