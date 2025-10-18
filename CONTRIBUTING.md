# Contributing to CaribX Backend

Thank you for your interest in contributing to CaribX! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce**
- **Expected vs actual behavior**
- **Screenshots** (if applicable)
- **Environment details** (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the proposed enhancement
- **Explain why this enhancement would be useful**
- **List any alternatives** you've considered

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the coding standards** described below
3. **Write tests** for your changes
4. **Ensure all tests pass** (`make test`)
5. **Run linters** (`make lint`)
6. **Update documentation** if needed
7. **Write a clear commit message** following conventional commits

## Development Process

### Setting Up Your Development Environment

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/CaribEx-backend.git
cd CaribEx-backend

# Add upstream remote
git remote add upstream https://github.com/Tenoywil/CaribEx-backend.git

# Install dependencies
make deps

# Start services
make docker-up

# Run migrations
make migrate-up
```

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clean, readable code
   - Follow Go best practices
   - Add tests for new functionality
   - Update documentation

3. **Test your changes**
   ```bash
   make test
   make lint
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request**

## Coding Standards

### Go Style Guide

- Follow the [official Go style guide](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `user`, `wallet`)
- **Files**: lowercase with underscores (e.g., `user_repository.go`)
- **Interfaces**: noun or adjective (e.g., `Repository`, `Validator`)
- **Functions**: camelCase, descriptive (e.g., `CreateUser`, `validateInput`)
- **Constants**: CamelCase with prefix (e.g., `DefaultTimeout`)

### Code Organization

Follow the established project structure:

```
internal/
  domain/       # Business logic, no infrastructure dependencies
  repository/   # Data access implementations
  usecase/      # Application logic
  controller/   # HTTP handlers
  routes/       # Route definitions
pkg/            # Reusable packages
```

### Comments

- Public functions and types must have godoc comments
- Comments should explain **why**, not **what**
- Use complete sentences with proper punctuation

```go
// CreateUser creates a new user in the system and returns the created user.
// It validates the input, checks for duplicates, and persists to the database.
func CreateUser(input UserInput) (*User, error) {
    // Implementation
}
```

### Error Handling

- Always handle errors explicitly
- Don't panic in library code
- Use custom error types for domain errors
- Add context to errors

```go
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

### Testing

- Write table-driven tests where appropriate
- Use meaningful test names (e.g., `TestCreateUser_WithValidInput_ReturnsUser`)
- Mock external dependencies
- Aim for >80% code coverage

```go
func TestWallet_Debit(t *testing.T) {
    tests := []struct {
        name    string
        balance float64
        amount  float64
        wantErr bool
    }{
        {"sufficient balance", 100.0, 50.0, false},
        {"insufficient balance", 30.0, 50.0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            wallet := &Wallet{Balance: tt.balance}
            err := wallet.Debit(tt.amount)
            if (err != nil) != tt.wantErr {
                t.Errorf("Debit() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **refactor**: Code refactoring
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks
- **ci**: CI/CD changes

### Examples

```
feat(wallet): add transaction history endpoint

Implement GET /v1/wallet/transactions endpoint with pagination
support. Includes caching and proper error handling.

Closes #123
```

```
fix(auth): prevent session hijacking

Add IP binding to sessions and rotate session IDs after
authentication to prevent session fixation attacks.
```

## Review Process

1. **Automated checks** run on all PRs (tests, linting, build)
2. **At least one maintainer** must approve the PR
3. **All conversations must be resolved** before merging
4. **Squash and merge** is preferred for cleaner history

## Documentation

Update documentation when:

- Adding new features
- Changing APIs
- Modifying configuration
- Updating dependencies

Documentation locations:

- **API changes**: `docs/API.md`
- **Architecture changes**: `docs/ARCHITECTURE.md`
- **Setup changes**: `README.md`
- **Code comments**: Inline godoc comments

## Questions?

If you have questions:

- Check existing [issues](https://github.com/Tenoywil/CaribEx-backend/issues)
- Ask in [discussions](https://github.com/Tenoywil/CaribEx-backend/discussions)
- Contact maintainers

Thank you for contributing to CaribX! 🚀
