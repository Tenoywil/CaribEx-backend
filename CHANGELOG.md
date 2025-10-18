# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project scaffolding with Domain-Driven Design structure
- Domain models for User, Wallet, Product, Cart, and Order
- PostgreSQL database schema with migrations
- **Marketplace seed data migration with 40+ Caribbean-themed products**
  - 8 seller accounts with wallet addresses and balances
  - Products across all categories (Electronics, Fashion, Food & Beverages, etc.)
  - Realistic Jamaican pricing and descriptions
  - Ready-to-use test data for development
- Docker and Docker Compose configuration
- Configuration management with environment variables
- Structured logging with zerolog
- Comprehensive documentation:
  - API documentation (docs/API.md)
  - Architecture guide (docs/ARCHITECTURE.md)
  - Development guide (docs/DEVELOPMENT.md)
  - Contributing guidelines (CONTRIBUTING.md)
  - Migrations documentation (migrations/README.md)
- Makefile with common development tasks
- MIT License
- README with project overview and quick start guide

### Infrastructure
- Go module initialization (Go 1.23+)
- Directory structure following DDD principles:
  - `cmd/` - Application entry points
  - `internal/domain/` - Domain models and business logic
  - `internal/repository/` - Data access layer
  - `internal/usecase/` - Application business logic
  - `internal/controller/` - HTTP handlers
  - `internal/routes/` - Route definitions
  - `pkg/` - Reusable packages
  - `migrations/` - Database migrations
  - `tests/` - Test files
  - `docs/` - Documentation

### Database
- PostgreSQL schema with the following tables:
  - users (with role-based access)
  - wallets (multi-currency support)
  - transactions (immutable ledger)
  - products (marketplace listings)
  - categories (product categorization)
  - carts and cart_items (shopping cart)
  - orders and order_items (order management)
  - refresh_tokens (session management)
- Indexes for optimal query performance
- Foreign key constraints with proper cascading

## [0.1.0] - 2025-10-18

### Project Initialization
- Repository created
- Basic project structure established
- Initial documentation added
