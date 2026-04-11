# Architecture Overview: Financial Ledger System

## 1. Core Purpose
The system is a **distributed financial ledger** designed to maintain an immutable, audit-ready record of financial transactions. It ensures the integrity of financial data through **double-entry accounting** principles, where every transaction is balanced by equal debits and credits.

## 2. Main Components and Responsibilities

### Service Layer (`services/`)
* **`LedgerService`**: Orchestrates complex business logic, such as initiating transactions, ensuring double-entry constraints are met, and updating both account balances and the transaction outbox.

### Repository Layer (`repos/`)
* **`AccountRepository`**: Manages the lifecycle and state of accounts (e.g., `Account` model) within PostgreSQL. Handles balance updates and account creation.
* **`TransactionRepository`**: Manages the persistent storage of ledger entries (e.g., `Transaction` model) in DynamoDB.
* **`OutboxRepository`**: Manages the `OutboxTransaction` queue in PostgreSQL. It tracks pending transactions that need to be "settled" (moved from the outbox to the permanent ledger).
* **`Repository` (Base)**: Provides common database interaction patterns.

### Model Layer (`models/`)
* Defines the core domain entities: `Account`, `Transaction`, and `OutboxTransaction`. It also defines domain-specific types like `TransactionSide` (DR/CR) and `OutboxTransactionStatus`.

## 3. Data Flow
The system implements a reliable transaction flow to ensure atomicity across different databases:

1.  **Initiation**: A request is received to perform a transaction.
2.  **Atomic Local Transaction (PostgreSQL)**:
    *   The system updates the `Account` balances (e.g., deducting from one account, adding to another) within a single PostgreSQL transaction.
    *   Simultaneously, it inserts a new `OutboxTransaction` into the `outbox` table within the **same database transaction**. This ensures that the account update and the intent to record a transaction are atomic.
3.  **Outbox Processing (Asynchronous)**:
    *   A background process (Worker) polls the `outbox` table for `PENDING` entries.
    *   For each pending entry, it attempts to write the transaction details to **DynamoDB** (the permanent Ledger).
4.  **Settlement**:
    *   Upon successful write to DynamoDB, the `OutboxTransaction` status is updated to `SETTLED` in PostgreSQL.
    *   If the write fails, the system uses a retry mechanism (tracked by `RetryCount`) or marks the entry as `DEAD` after exhaustion.

## 4. Key Design Patterns
* **Double-Entry Accounting**: The fundamental domain pattern where every entry requires a corresponding debit and credit to keep the books balanced.
* **Transactional Outbox Pattern**: Used to solve the "dual-write" problem. It ensures that the PostgreSQL state (accounts) and the DynamoDB state (ledger entries) stay synchronized without requiring distributed transactions (2PC).
* **Immutability**: The ledger entries in DynamoDB are designed to be append-only. Changes to financial state are handled via reversing entries rather than modifying existing records.
* **Repository Pattern**: Decouples the business logic (Service layer) from the underlying data access logic (PostgreSQL/DynamoDB).

## 5. Data Models
* **`Account`**: Holds balance, metadata, and hierarchy information in PostgreSQL.
* **`Transaction`**: Represents an immutable ledger entry in DynamoDB.
* **`OutboxTransaction`**: Represents a pending transaction in the PostgreSQL outbox, acting as a buffer for the asynchronous sync to DynamoDB.

## 6. Technology Stack
* **Programming Language**: Go (Golang)
* **Relational Database**: PostgreSQL (for Accounts and Outbox)
* **NoSQL Database**: Amazon DynamoDB (for Transaction/Ledger history)
* **Orchestration**: Docker & Docker Compose
