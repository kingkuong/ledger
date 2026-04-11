# Distributed Systems Engineer Skills

## Core Expertise
- **Primary Language**: Go (Golang)
- **Distributed Systems**: Design and implementation of scalable, fault-tolerant systems
- **Database Systems**: PostgreSQL, DynamoDB, and hybrid database patterns

## Technical Skills

### Go Programming
- Modern Go (1.18+) with generics
- Concurrency patterns (goroutines, channels, sync)
- Error handling and resilient design
- Testing (unit, integration, table-driven tests)
- Code organisation and package design

### Distributed Systems Patterns
- Eventual consistency
- Two-phase commit and saga patterns
- Outbox pattern for data consistency
- Idempotency and retry mechanisms
- Circuit breakers and fault tolerance

### Database Integration
- PostgreSQL: relational data, ACID transactions
- DynamoDB: NoSQL, key-value/document storage
- Hybrid storage strategies
- Database migrations and schema management

## Project Context: Ledger System
This project implements a financial ledger with:
- Multi-side accounting (DEBIT/CREDIT)
- Cross-database consistency (PostgreSQL + DynamoDB)
- Eventual sync via outbox pattern
- Retry logic for fault tolerance
