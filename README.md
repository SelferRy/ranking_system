# Ranking system service
The ranking system implements banner rotation microservice.
This service is designed to select the most effective (clickable) banners in
conditions of changing user preferences and a set of banners.
It consists of an API and a database that stores information about banners.
The service provides a gRPC API.

ranking_system has several slots and banners.
A slot is a specific API that a user can interact with.
Each slot can have any number of banners.
Each banner can be in different slots.
Customers are divided into socio-demographic groups. Banners are displayed according to their preferences.

The microservice sends click and impression events to a queue (Kafka) for further processing in analytics systems.

## Quick start
You can clone repo to your own machine:
```bash
git clone https://github.com/SelferRy/ranking_system.git 
```
And then, inside root directory of the project use the commands:
```bash
make run
make build
make test
```
`make run` - this command start to download all dependent containers, start them and run. 
After that it will start server. Then you can test in in another terminal. See commands below.

Select banner:
```bash
grpcurl -plaintext -import-path api/proto -proto api/proto/banner_rotator.proto -d '{"slot_id": 1, "group_id": 1}' localhost:5080 ranking_system.BannerRotatorService/SelectBanner
```

Click banner:
```bash
grpcurl -plaintext -import-path api/proto -proto api/proto/banner_interaction.proto -d '{"banner_id": 1, "slot_id": 1, "group_id": 1}' localhost:5080 ranking_system.BannerInteractionService/ClickBanner
```

You can create your banner:
```bash
grpcurl -plaintext \
  -import-path api/proto \
  -proto banner_management.proto \
  -d '{"banner": {"id": 3, "description": "Test banner"}, "slot_id": 1}' \
  localhost:5080 ranking_system.BannerManagementService/AddBannerToSlot
 ```
And remove it after that:
```bash
grpcurl -plaintext \
  -import-path api/proto \
  -proto banner_management.proto \
  -d '{"banner": {"id": 3, "description": "Test banner"}, "slot_id": 1}' \
  localhost:5080 ranking_system.BannerManagementService/RemoveBannerFromSlot
```
In current version of the project this is all API endpoints and no reflection grpc function added.

For finish your work with the project, you have to terminate server (Ctrl+C, for example) 
and then inactivate all dependencies:
```bash
make down 
```