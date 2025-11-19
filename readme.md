# Loan Management System

## Architecture Overview

```mermaid
flowchart TD
    A[HTTP Client] --> B[Delivery Layer]
    B --> C[Use Case Layer]
    C --> D[Repository Layer]
    D --> E[Domain Model]
    style E fill:#f9f,stroke:#333

    subgraph Clean Architecture
        direction TB
        B{{Echo HTTP Router}}
        C{{Loan Use Cases}}
        D{{Loan Repository}}
        E[(Loan Model)]
    end
```

## Core Components

### Loan State Machine

```mermaid
flowchart TB
    Start([Start]) --> Proposed
    Proposed -->|Admin Approval| Approved
    Approved -->|Full Funding| Invested
    Invested -->|Payout| Disbursed
    Disbursed --> Finish([End])

    subgraph Proposed_State [Proposed]
        direction TB
        PEntry([Entry: Validate borrower info])
    end

    subgraph Approved_State [Approved]
        direction TB
        AEntry([Entry: Set approval metadata])
    end

    subgraph Invested_State [Invested]
        direction TB
        IEntry([Entry: Track investments])
    end
```

## Key Packages

| Package | Responsibility |
|---------|----------------|
| `internal/model` | Domain entities and business rules |
| `internal/usecase` | Business transaction orchestration |
| `internal/repository` | Data persistence (memory implementation) |
| `internal/delivery/http` | Echo web handlers and routes |

## Sequence Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant H as HTTP Handler
    participant U as Loan Use Case
    participant R as Repository
    
    C->>H: POST /loans
    H->>U: CreateLoan()
    U->>R: Save()
    R-->>U: Loan
    U-->>H: Loan
    H-->>C: 200 OK

    C->>H: PUT /loans/:id/approve
    H->>U: ApproveLoan()
    U->>R: Update()
    R-->>U: Loan
    U-->>H: Loan
    H-->>C: 200 OK

    C->>H: PUT /loans/:id/invest
    H->>U: InvestLoan()
    U->>R: Update()
    R-->>U: Loan
    U-->>H: Loan
    H-->>C: 200 OK

    C->>H: PUT /loans/:id/disburse
    H->>U: DisburseLoan()
    U->>R: Update()
    R-->>U: Loan
    U-->>H: Loan
    H-->>C: 200 OK
```

## Development Setup

```bash
# Install dependencies
go mod tidy

# Run service
go run main.go
```

## Testing

```bash
# Run all tests with coverage
go test -cover ./...
```
