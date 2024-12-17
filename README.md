# Auth app

## Discription
Тестовое задание: https://medods.notion.site/Test-task-BackDev-623508ed85474f48a721e43ab00e9916

### Has been used
1. Go
2. Postgres
2. Docker
1. GNU Make

### How to use
1. Clone the repository:
    ```bash
    git clone https://github.com/l1qwie/JWTAuth.git
2. Go to the major enter:
    ```bash
    cd Auth
3. Create a docker-network for an app and a database
    ```bash
    make net
4. Build an app image
    ```bash
    make build.app
5. Build a database image
    ```bash
    make build.db
6. Run the database container
    ```bash
    make run.db.innet
7. Run the app container
    ```bash
    make run.app

Ta-daa! The app is working!

#### A few recommendations
1. If you want to use this app to make sure everything works well, you're supposted to create a row about your device in the database. Ecpesially there columns: guid, email and ip. To do it, do this:
    ```bash
    make get.into
2. If you want to run tests (my own tests), you could do it, only if you run the database container not in the just created network. To do it, do this:
    ```bash
    make run.db.outnet
3. And the finals step is:
    ```bash
    cd tests
and choose whatever you want to see: e2e, unit or integrated. Then go to this folder and use command    
    ```bash
    go test
to be done all tests in the folder
