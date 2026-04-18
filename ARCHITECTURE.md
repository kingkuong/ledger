# Architecture Overview: Financial Ledger System

## 1. Core Purpose
The system is a **distributed financial ledger** designed to maintain an immutable, audit-ready record of financial transactions. It ensures the integrity of financial data through **double-entry accounting** principles, where every transaction is balanced by equal debits and credits.

## 2. Main Components and Responsibilities

### Service Layer (`services/`)
* **`LedgerService`**: Orchestrates complex business logic, such as initiating transactions, ensuring double-entry constraints are met, updating account balances, and managing the transaction outbox. Handles `Transfer` (cross-account) and `CreateTransaction` (single-account) flows.

### Repository Layer (`repos/`)
* **`AccountPostgresRepo`**: Manages the lifecycle and state of `Account` entities in PostgreSQL. Includes `ReadForUpdate` (SELECT FOR UPDATE) for optimistic concurrency during balance updates.
* **`TransactionRepo`**: Handles persistent storage of `Transaction` ledger entries in DynamoDB.
* **`OutboxPostgresRepo`**: Manages the `outbox` table in PostgreSQL. Tracks `PENDING` transactions awaiting settlement, supports `Update` (status + retry tracking), and `FetchPending` for the outbox worker.
* **No base repository** — each repo is implemented independently.

### Process Layer (`processes/`)
* **`OutboxProcess`**: Background worker that polls the `outbox` table at intervals (configurable ticker), writes pending entries to DynamoDB, and updates outbox status to `SETTLED` on success.

### Model Layer (`models/`)
* **`Account`**: Account entity with balance, account number, optional nickname.
* **`Transaction`**: Immutable ledger entry for DynamoDB. Embedded in `OutboxTransaction`.
* **`OutboxTransaction`**: Wraps `Transaction` with status (`PENDING`/`SETTLED`/`DEAD`), retry count, and retry timestamp for the outbox pattern.
* **`TransactionSide`**: Enum — `DR` (debit), `CR` (credit).
* **`OutboxTransactionStatus`**: Enum — `PENDING`, `SETTLED`, `DEAD`.

## 3. Data Flow

### Transfer Flow (Cross-Account)
1. **Begin PG Transaction**: `pool.Begin()`
2. **Lock Accounts**: `ReadForUpdate` on both source and destination accounts (`SELECT ... FOR UPDATE`)
3. **Validate**: Check `From.Balance >= amount`
4. **Update Balances**: `UpdateBalance` on both accounts (decrement from, increment to)
5. **Create Outbox Entries**: Two `OutboxTransaction` records inserted into the same PG transaction:
   - CR entry (destination account, amount positive)
   - DR entry (source account, amount negative)
   - Both share cross-references via `CounterPartID`
6. **Commit**: If successful, account balances and outbox entries are persisted atomically.
7. **Outbox Processing (Async)**: `OutboxProcess` polls for `PENDING` outbox entries, writes each to DynamoDB, then marks them `SETTLED`.

### Single-Account Transaction Flow
1. **Begin PG Transaction**
2. **Lock Account**: `ReadForUpdate` on the account
3. **Validate Balance**: Check sufficient balance for debit transactions
4. **Update Balance**: Apply amount delta
5. **Commit**: Balance persisted in PG (Note: DynamoDB write happens outside the PG transaction — an area for future improvement)

## 4. Key Design Patterns
* **Double-Entry Accounting**: Every transfer creates two outbox entries with equal but opposite debits/credits.
* **Transactional Outbox Pattern**: Solves the dual-write problem between PostgreSQL and DynamoDB. Account state and outbox entries are written atomically in a single PG transaction.
* **Immutability**: DynamoDB ledger entries are append-only. Reversing entries correct errors rather than modifying existing records.
* **Advisory Locking**: `ReadForUpdate` (SELECT FOR UPDATE) prevents concurrent balance modification.
* **Idempotency**: Outbox entries carry UUIDs generated at the application level, enabling safe retries.

## 5. Data Models
* **`Account`**: Stored in PostgreSQL. Has `id`, `account_number`, optional `nickname`, `balance` (bigint, smallest unit).
* **`Transaction`**: Stored in DynamoDB. Primary key components: `PK`, `SK`. Contains side, amount, description, counter-part ID, account ID.
* **`OutboxTransaction`**: Stored in PostgreSQL `outbox` table. Wraps `Transaction` with `status`, `retry_count`, `last_retry_at`. Tracks delivery state of each transaction to DynamoDB.

## 6. Technology Stack
* **Programming Language**: Go (Golang)
* **Relational Database**: PostgreSQL 18 (Accounts + Outbox)
* **NoSQL Database**: Amazon DynamoDB Local (Permanent Ledger)
* **Orchestration**: Docker & Docker Compose
* **Database Driver**: pgx v5

## 7. Current Limitations
* No HTTP/API layer (main.go uses manual seeding)
* Single-account `CreateTransaction` writes to DynamoDB outside the PG transaction — potential inconsistency if DynamoDB write fails
* No exponential backoff or DEAD marking on outbox retries
* No account type classification (Asset/Liability/Equity/etc.)
* No multi-currency support
* No reporting endpoints
* No RBAC or user authentication
