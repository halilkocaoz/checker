# upsmo-checker

It runs in different locations on the Azure VMs and checks uptime monitors by region. Also, it takes actions according to the HTTP responses, such as sending messages to the Azure Service Bus.

## Prerequisites to run

* [Golang](https://golang.org/dl/)

## Installation & Run

* `git clone https://github.com/halilkocaoz/upsmo-checker.git`
* `export UPSMO_REGION="DE_Frankfurt"`  
Look [here](https://github.com/halilkocaoz/upsmo-server/blob/main/UpsMo.Common/Monitor/MonitorRegion.cs) for more
* `export AZURE_POSTGRES_CONNSTR="host=name.postgres.database.azure.com port=5432 dbname=- user=- password=-"`
* `export SERVICE_BUS_NAMESPACE="namespace"`
* `export SERVICE_BUS_SHARED_ACCESS_KEY_VALUE="key"`
* `cd upsmo-checker && go mod download && go run .`

## Related repositories

* [upsmo-server](https://github.com/halilkocaoz/upsmo-server)
* [upsmo-inserter](https://github.com/halilkocaoz/upsmo-inserter)
* [upsmo-notifier](https://github.com/halilkocaoz/upsmo-notifier)

## Screenshot

![cli](./assets/checker.png)
