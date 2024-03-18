# Digital Wallet

Digital Wallet is a simple web application that consists of two microservices,
**User** and **Transaction**. <br/>

The **User** microservice is responsible for creating users and retrieving user balance
from the Transaction microservice. <br/>

The **Transaction** microservice is responsible for executing transactions.
Type of transactions supported are: crediting money to a user account,
and transferring money between two users. <br/>

Both of the microservices use **PostgreSQL** as their database. <br/>

https://nats.io/ is used for synchronous communication between the microservices. <br/>
User service sends request for retrieving user balance to Transaction service,
and Transaction service sends response back to User service. <br/>

https://kafka.apache.org/ is used for event-driven communication between the microservices. <br/>
User service sends event to Transaction service when creating user, and Transaction service responds
to User service with Success or Failure message, depending on whether the user was created or not. <br/>

**Docker** is used for containerization of the PostgreSQL, NATS and Kafka. <br/>

Some other interesting technologies used in this project:
- Go modules for dependency management
- Gin Gonic for web framework
- GORM for ORM
- Wire for dependency injection
- Cobra for CLI

**Database initialization:**
File `./docker/postgres_init/create-databases.sh` is a script that initializes the databases for the microservices. <br/>

## How to start application?

1. Position yourself in the root directory of the project.
2. Execute `docker-compose` script.
```
cd docker
docker-compose up -d
```
4. **Wait for your containers to be started!**
5. Open 2 new terminal windows in the root directory of the project.
6. In the first terminal window, position yourself in the `user-service` directory.
```
cd user-service
go mod tidy
go run main.go reinit db
go run main.go serve
```
7. In the second terminal window, position yourself in the `transaction-service` directory.
```
cd transaction-service
go mod tidy
go run main.go reinit db
go run main.go serve
```

## API
### User Microservice
User Microservice runs on port **8080**.

#### Create User
```
[POST] 
localhost:8080/api/v1/users

[Request]
{
    "email": "user@email.com"
}

[Response]
{
    "data": {
        "user_id": "dc44c491-e50e-48b9-9213-06c7a84bd95e",
        "email": "user@email.com",
        "created_at": "2024-03-17T00:40:20.400046+01:00"
    }
}
```

#### Get User Balance
```
[GET] 
localhost:8080/api/v1/users/balance/:email

[Request]
localhost:8080/api/v1/users/balance/user@email.com

[Response]
{
    "data": {
        "email": "user@email.com",
        "balance": 1000
    }
}
```

### Transaction Microservice
Transaction Microservice runs on port **8000**.

#### Add Money
```
[POST] 
localhost:8000/api/v1/transactions/add-money

[Request]
{
    "user_id": "6d6e69bf-167b-41b8-8438-f7b839c760eb",
    "amount": 1000
}

[Response]
{
    "data": {
        "balance": 1600
    }
}
```

#### Transfer Money
```
[POST] 
localhost:8000/api/v1/transactions/transfer-money

[Request]
{
    "to_user_id": "3202585d-fa4a-46c3-9ad8-98b561da0795",
    "from_user_id": "32c1a111-6ee8-4c20-8dab-bee13746c56c",
    "amount": 200
}

[Response]
{} // if empty response is returned, it means that transfer was successful
```
