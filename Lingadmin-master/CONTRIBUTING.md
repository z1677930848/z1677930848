# Contributing to Lingadmin

Thank you for your interest in contributing to Lingadmin!

## Development Setup

### Prerequisites
- Go 1.22 or higher
- Node.js (for frontend development)
- Git

### Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd Lingadmin-master
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o ling-admin ./cmd/edge-admin
```

4. Run the application:
```bash
./ling-admin
```

## Project Structure

```
├── cmd/              # Application entry points
├── internal/         # Private application code
├── web/              # Frontend assets
│   ├── public/       # Static files
│   └── views/        # Templates
├── configs/          # Configuration files
├── doc/              # Documentation and diagrams
└── docs/             # Consolidated documentation
```

## Coding Standards

- Follow Go best practices and conventions
- Write clear, self-documenting code
- Add comments for complex logic
- Use meaningful variable and function names

## Submitting Changes

1. Create a new branch for your feature/fix
2. Make your changes
3. Test thoroughly
4. Submit a pull request with clear description

## Reporting Issues

Please report issues on the project issue tracker with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)

## Questions?

Join our community:
- QQ群: 659832182
- Telegram: https://t.me/+5kVCMGxQhZxiODY9
