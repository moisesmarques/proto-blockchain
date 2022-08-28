### Getting Started

Ensure you have setup go in your local machine.
[https://go.dev/doc/install](https://go.dev/doc/install)

### Prerequisites

Clone the repo

```sh
git clone https://github.com/moisesmarques/judge-blockchain.git
```

### Building the binary

```sh
env GOOS=darwin GOARCH=amd64 go build -o ./build/
```

visit below link for more GOOS and GOARCH values
https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63

### Run GRPC Server

```sh
go run main.go startAsgRPC --blockchainPort 3080 --grpcPort 9000
```

blockchainPort is the port where httpServer of blockchain is running (kind of central node for time being)

### Usage

```sh
go run main.go startAsHTTP --port 3080

open http client e.g postman
base url: http://localhost:3080/
```

## Troubleshooting

If you are running the command adobe and you find an error that is similar to this one:

```
../../../go/pkg/mod/golang.org/x/sys@v0.0.0-20190626221950-04f50cda93cb/unix/zsyscall_darwin_amd64.go:28:3: //go:linkname must refer to declared function or variable
```

You need to run this first: `go get -u golang.org/x/sys`

#### Functions

**1. List Wallets**

```sh
GET 'localhost:3080/wallets'
```

```
Address: 14NCm5HBb8bCMiihicfh72Yor24Nde9WXx  Balance: 0
Address: 1MAD2sfEnSKx2PQP4qPss1sPHgxGEi8guF  Balance: 0
```

**2. Create Wallet**

```sh
POST 'localhost:3080/wallets'
```

```
New address is: 1KWJMpzQomVxFKqnmqpLycd8jsFqzyTZ4

```

**3. Create Chain**

Initialize a blockchain database if it does not exist

```sh
POST 'localhost:3080/chain' --form wallet="14NCm5HBb8bCMiihicfh72Yor24Nde9WXx"'
```

```
********************Printing chain data****************************************************************
***************************************************
Previous hash:
hash: 0006eb67d5ab51cbe53085b67e9127517cfc14311772e611c8c640405863f54d
Pow: true
 --- Transaction e0dffd2f014a678ea9b142ecb6ea728821a47cba1cf3049c216f3544395fb57e:
     Input 0:
       TXID:
       Out:       -1
       Signature:
       PubKey:    4669727374205472616e73616374696f6e2066726f6d2047656e65736973
     Output 0:
       Value:  100
       Script: 24eaede0e690b6956ce676f7a3e198430047fa15
     ValidationResult:
       Value:  {   [] {   } {0 0 {} [] 0}}
Height: 0
Nonce: 3161
Timestamp: 1655894144

```

**4. Print Chain**

Display chain blocks

```sh
GET 'localhost:3080/chain'
```

```
Previous hash: 000db3dc10d3849c48cd1318133b8d6df0940d37eb6292716674e1dcc4a8c29b
hash: 00024739d39c63ac973c107a56242134c8c5d3f4ff3206fba79bb0316ded195a
Pow: true
 --- Transaction 3bc503339d284d070217980bb115681b3b1695c8b282eb905f23d80850cf7052:
     Input 0:
       TXID:     50366ff145e6a5ebc25e599cc52396579ab0f8be50ef213e349ebb18f4d6fd25
       Out:       1
       Signature: 62948af249d337838a3d49d024fe68ae7cffe40ce35da800bd09e2237482700d490453d7e8b01d8d0e3dc2a46c0fb0df173bacc619d6eb3f540e4d1e21bbad18
       PubKey:    bdae9766607dcc6abaf2dfec286b7c5e1a67d95ca2396ead2f2b8eca652bbd1f18e66eb30cf9738cda2cf292d2e143931cd73dd3e05c011e18693bdf4a623da3
     Output 0:
       Value:  1
       Script: dd204daafe2ce3cb6a5b1cc4a1487da73404f898
     Output 1:
       Value:  96
       Script: 24eaede0e690b6956ce676f7a3e198430047fa15
     ValidationResult:
       Value:  {   [] {   } {0 0 {} [] 0}}
Height: 2
Nonce: 885
Timestamp: 1655894752
***********************************************************************************************
Previous hash: 0006eb67d5ab51cbe53085b67e9127517cfc14311772e611c8c640405863f54d
hash: 000db3dc10d3849c48cd1318133b8d6df0940d37eb6292716674e1dcc4a8c29b
Pow: true
 --- Transaction 50366ff145e6a5ebc25e599cc52396579ab0f8be50ef213e349ebb18f4d6fd25:
     Input 0:
       TXID:     e0dffd2f014a678ea9b142ecb6ea728821a47cba1cf3049c216f3544395fb57e
       Out:       0
       Signature: 82794366d3719f7e7b83433747b7a95f193823f1455e87f18504f1ce5c6792dd9970615f38d1e8bb860d5d17b8f6c33cf8b2bc12c009b63ddf3768fffec98f97
       PubKey:    bdae9766607dcc6abaf2dfec286b7c5e1a67d95ca2396ead2f2b8eca652bbd1f18e66eb30cf9738cda2cf292d2e143931cd73dd3e05c011e18693bdf4a623da3
     Output 0:
       Value:  3
       Script: dd204daafe2ce3cb6a5b1cc4a1487da73404f898
     Output 1:
       Value:  97
       Script: 24eaede0e690b6956ce676f7a3e198430047fa15
     ValidationResult:
       Value:  {   [] {   } {0 0 {} [] 0}}
Height: 1
Nonce: 6809
Timestamp: 1655894415
*********************************************************************************************
Previous hash:
hash: 0006eb67d5ab51cbe53085b67e9127517cfc14311772e611c8c640405863f54d
Pow: true
 --- Transaction e0dffd2f014a678ea9b142ecb6ea728821a47cba1cf3049c216f3544395fb57e:
     Input 0:
       TXID:
       Out:       -1
       Signature:
       PubKey:    4669727374205472616e73616374696f6e2066726f6d2047656e65736973
     Output 0:
       Value:  100
       Script: 24eaede0e690b6956ce676f7a3e198430047fa15
     ValidationResult:
       Value:  {   [] {   } {0 0 {} [] 0}}
Height: 0
Nonce: 3161
Timestamp: 1655894144
*********************************************************************************************

```

**5. Send Tokens**

```sh
POST 'localhost:3080/sendTokens' \
--form 'from="14NCm5HBb8bCMiihicfh72Yor24Nde9WXx"' \
--form 'to="1MAD2sfEnSKx2PQP4qPss1sPHgxGEi8guF"' \
--form 'amount="3"'

```

```

Tokens send successfully
```

**6. Wallet Balances**

```sh
GET 'localhost:3080/wallets'

```

```
Address: 1KWJMpzQomVxFKqnmqpLycd8jsFqzyTZ4  Balance: 0
Address: 14NCm5HBb8bCMiihicfh72Yor24Nde9WXx  Balance: 96
Address: 1MAD2sfEnSKx2PQP4qPss1sPHgxGEi8guF  Balance: 4
```

**7. Add Peers**

```sh
POST 'localhost:3080/peers' \
--form 'host="localhost"' \
--form 'port="3082"'
```

```
Peer added
```

**7. List Peers**

```sh
GET 'localhost:3080/peers'
```

```
 localhost:3080
 localhost:3082
 localhost:3081
All addresses printed
```

**7. Save ValidationResults**

```sh
POST 'localhost:3080/newActionResult' \
--header 'Content-Type: application/json' \
--data-raw '{
    "ValidatorID": "12345",
    "FAK":"anotherone",
    "Results": [
        {
            "ValidatorID": "A",
            "Result": true
        },
         {
            "ValidatorID": "B",
            "Result": false
        },
         {
            "ValidatorID": "C",
            "Result": true
        },
         {
            "ValidatorID": "18181881",
            "Result": true
        }
    ],
    "Action": {
        "ActionID": "18181881",
        "ActionType": "saveFile",
        "ActionData": "resources",
        "ResourceID": "122222"
    }
}'
```

```
Previous hash: 00024739d39c63ac973c107a56242134c8c5d3f4ff3206fba79bb0316ded195a
hash: 00082da3361afe9126c8cc370a2fa3d69da92b369d18cb990c8bc35a72d42d07
Pow: true
Previous hash: 000db3dc10d3849c48cd1318133b8d6df0940d37eb6292716674e1dcc4a8c29b
hash: 00024739d39c63ac973c107a56242134c8c5d3f4ff3206fba79bb0316ded195a
Pow: true
Previous hash: 0006eb67d5ab51cbe53085b67e9127517cfc14311772e611c8c640405863f54d
hash: 000db3dc10d3849c48cd1318133b8d6df0940d37eb6292716674e1dcc4a8c29b
Pow: true
Previous hash:
hash: 0006eb67d5ab51cbe53085b67e9127517cfc14311772e611c8c640405863f54d
Pow: true

```

**8. Running Multiple Nodes**
i. use cli commands on terminal to start example nodes as TCP server e.g

```
go run main.go startAsTCP --port 3081
go run main.go startAsTCP --port 3082
```

Start one node as an HTTP node to allow inputting of sample data into the network via http

```
go run main.go startAsHTTP --port 3080
```

Minimum is 1 tcp node and one http node. Then you can add other tcp nodes. Sync will require an already existing tcp node running.
