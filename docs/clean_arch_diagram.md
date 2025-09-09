```mermaid
flowchart LR
subgraph Domain (чистая доменная модель)
E[entity/*]
S[service/bandit]
R[(repository interfaces)]
end
subgraph UseCase (application)
U[usecase/banner.UseCase]
end
subgraph Infra (adapters)
DB[infra/adapters/repository/sql/pg]
MQ[infra/adapters/broker/kafka]
L[infra/logger]
end
subgraph InterfaceAdapters (delivery)
H[server/grpc/handler/*]
PB[server/grpc/pb]
SV[server/grpc/server]
end
M[cmd/.../main.go]

M --> SV
SV --> H
H --> U
U -->|depends on| R
U --> S
U --> MQ
DB -.implements .-> R
```