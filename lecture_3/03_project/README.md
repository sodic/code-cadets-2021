# Lecture 3 Project

We are building a simple betting platform. You can find further information in the lecture presentation and in this file.

## Project structure
- `betgenerator`: Utility script which generates and produces bets into `bets-received` queue.
- `calculator`: Service which settles bets and calculates payout value for winning bets.
- `controller`: Service which owns the database of bets and updates the DB based on incoming event updates.
- `db`: Utility directory where all databases should be placed.
- `dbinitializer`: Utility scripts which sets up databases. If a database already exists, it will be overridden.
- `eventsettler`: Utility script which generates and produces event updates into `event-updates` queue.

## RabbitMQ queues
- The following queues are used:
    - `bets-received`: a queue where new bets are incoming
    - `bets`: a queue where all bets are produced (on each creation/update)
    - `event-updates`: a queue where updates of all bets are incoming
    - `bets-calculated`: a queue where updates of all bets are incoming
- Example flow of messages for a bet:
    1. `bets-received`: `{"id":"e7adac7e-6a7a-4fa8-5687-8df976dc9a92","customerId":"f02c288b-0abe-48f9-4196-5516f4575e4c","selectionId":"67273877-29bd-47b7-5f7e-f7779682ad5a","selectionCoefficient":7.7,"payment":345.47}`
    1. `bets`: `{"id":"e7adac7e-6a7a-4fa8-5687-8df976dc9a92","customerId":"f02c288b-0abe-48f9-4196-5516f4575e4c","status":"active",selectionId":"67273877-29bd-47b7-5f7e-f7779682ad5a","selectionCoefficient":7.7,"payment":345.47}`
    1. `event-updates`: `{"id":"67273877-29bd-47b7-5f7e-f7779682ad5a","outcome":"won"}`
    1. `bets-calculated`: `{"id":"e7adac7e-6a7a-4fa8-5687-8df976dc9a92","status":"won","payout":2660.12}`
    1. `bets`: `{"id":"e7adac7e-6a7a-4fa8-5687-8df976dc9a92","customerId":"f02c288b-0abe-48f9-4196-5516f4575e4c","status":"won",selectionId":"67273877-29bd-47b7-5f7e-f7779682ad5a","selectionCoefficient":7.7,"payment":345.47,"payout":2660.12}`

## Databases:
- SQLite is used; databases are stored as files in `db` directory.
- `bets` db:
    - This is the main database of bets. It is a single source of truth for bets data.
    - table `bets`:
        - data model:
            - `id TEXT NOT NULL PRIMARY KEY` (this is UUID of the bet)
            - `customer_id TEXT NOT NULL` (this is UUID of the customer)
            - `status TEXT NOT NULL`
            - `selection_id TEXT NOT NULL` (this is UUID of the selection)
            - `selection_coefficient INTEGER NOT NULL` (e.g. coefficient 12.34 will be represented as 1234)
            - `payment INTEGER NOT NULL` (e.g. payment of 7.43 will be represented as 743)
            - `payout INTEGER` (same thing as selection_coefficient and payment)
        - does not have a secondary index
- `calc_bets` db:
    - This is a local database of the Calculator service. It contains a subset of bets data relevant for the Calculator.
    - table `bets`:
        - data model:
            - `id TEXT NOT NULL PRIMARY KEY` (this is UUID of the bet)
            - `selection_id TEXT NOT NULL` (this is UUID of the selection)
            - `selection_coefficient INTEGER NOT NULL` (e.g. coefficient 12.34 will be represented as 1234)
            - `payment INTEGER NOT NULL` (e.g. payment of 7.43 will be represented as 743)
        - secondary index:
            - `selection_idx`, indexing on `selection_id`; useful in order to fetch all bets with the desired selection

## Other information:
- Valid bet statuses are:
    - `active` - selection is not settled yet
    - `won` - selection was settled as won
    - `lost` - selection was settled as lost
- Floating-point values should be always rounded to 2 decimal places.

## Help:
- In order to run your application, configure "run configurations" in Goland; example for Controller service:
    - go to "Run" -> "Edit configurations..."
    - under "Go Build" click "+"
    - Name: "Run Controller"
    - Run kind: "Package"
    - Package path: "github.com/superbet-group/code-cadets-2021/lecture_3/03_project/controller/cmd"
    - Working directory (replace with your path!): "/Users/oreskd/Axilis/code-cadets/2021/code-cadets-2021/lecture_3/03_project/controller"
    - under "EnvFile" check "Enable EnvFile" and click "+" under the table (which currently contains only one entry)
    - choose `.env` file from `controller/` directory (note: the file may be hidden, you need to show hidden files in that case)
    