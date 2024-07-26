# Subad: Dynamic Subscription-based Access Control System

Subad is a powerful and flexible subscription-based access control system built with Go, templ, Tailwind CSS, and HTMX. It enables administrators to manage user access to various pages or content, while subscribers can request, trade, and earn access to specific resources using a unique status-token system.

## Core Features

- User authentication with role-based access control
- Dynamic page content management and ownership
- Subscription-based access requests and approvals
- Status-token based "bartering economy" for feature access
- User status system influencing page access and temporary ownership
- Interactive feature pages (e.g., roulette wheel) for earning status
- Comprehensive admin dashboard for user and content management
- Responsive and modern UI powered by Tailwind CSS
- Seamless, JavaScript-free interactions using HTMX

## Technical Stack

- Backend: Go
- Frontend: templ for HTML templating
- Styling: Tailwind CSS
- Interactivity: HTMX
- Database: SQLite (development), scalable to other RDBMS for production

## Getting Started

1. Clone the repository
2. Install Go (version 1.16 or later)
3. Install required dependencies:



go mod tidy

4. Set up your development environment (see Environment Setup section)
5. Run the application:



go run main.go


## Environment Setup

1. Create a `.env` file in the project root with the following content:



DB_PATH=./subad.db PORT=8080

2. Ensure you have SQLite installed on your system
3. For private repository access, configure Git to use SSH:



git config --global url."git@github.com:".insteadOf "https://github.com/"


## Module Management and Dependency Resolution

1. To add a new dependency:



go get github.com/example/package

2. To update dependencies:



go mod tidy

3. For private repositories, set the GOPRIVATE environment variable:



export GOPRIVATE=github.com/yourusername/privaterepo


## Code Structure and Organization

- `database/`: Database operations and connections
- `models/`: Data structures and business logic
- `handlers/`: HTTP request handlers
- `tokens/`: Token management system
- `sessions/`: User session management
- `templates/`: HTML templates using templ

## Development Workflow

1. Create a new feature branch
2. Implement changes in small, testable increments
3. Write and update tests for new functionality
4. Use the LLM agent for code review and suggestions
5. Update documentation, including godoc comments
6. Submit a pull request for review

## Testing

1. Write unit tests in `*_test.go` files alongside the code they're testing
2. Run tests:



go test ./...

3. For coverage report:



go test -cover ./...


## Development with LLM Agent Assistance

1. Provide the entire codebase context when asking for assistance
2. Ask the LLM to review the codebase before making suggestions
3. Instruct the LLM to maintain the OOP approach and modular structure
4. For big updates, ask the LLM to suggest creating new branches
5. Request the LLM to print complete updated files one at a time
6. Use this README as an instruction manual for the LLM

### Improved LLM Agent Guidelines

- Start each interaction by clearly stating the current task or problem
- Provide context about recent changes or decisions made in the project
- Ask the LLM to explain its reasoning for suggested changes or improvements
- Request step-by-step instructions for complex implementations
- Encourage the LLM to suggest alternative approaches when applicable
- Ask for code examples or pseudo-code to illustrate concepts
- Prompt the LLM to consider edge cases and potential issues in its suggestions
- Request that the LLM highlight any assumptions it's making about the project structure or requirements
- Periodically ask the LLM to summarize the current state of the project and suggest next steps
- Encourage the LLM to provide comments in the code for complex logic or design decisions

## Troubleshooting Guide

1. Compilation errors:
- Ensure all necessary packages are imported
- Check for typos in variable or function names
- Verify that all used packages are in the go.mod file
2. Runtime errors:
- Check database connection string
- Ensure all required environment variables are set
- Verify file permissions for database and log files
3. Dependency issues:
- Run `go mod tidy` to resolve and update dependencies
- For private repos, ensure SSH keys are correctly set up

## Deployment Strategy

1. Set up a CI/CD pipeline (e.g., GitHub Actions, GitLab CI)
2. Create Docker containers for consistent deployments
3. Use environment variables for configuration in different environments
4. Implement database migrations for schema changes
5. Set up monitoring and logging (e.g., Prometheus, ELK stack)

## Contributing

1. Fork the repository
2. Create a new branch for your feature
3. Commit changes with clear, descriptive messages
4. Push your branch and submit a pull request
5. Ensure your code adheres to the project's style guide and passes all tests

## License

[Specify the license under which the project is released.]