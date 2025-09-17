# Configuration Guide

ChatApp now uses a YAML-based configuration system powered by Viper. This provides flexible configuration management with support for environment variables, default values, and multiple configuration sources.

## Configuration File

The main configuration file is `config.yaml` in the root directory. Here's the complete structure:

```yaml
server:
  port: ":8080"                    # Server port
  host: "localhost"                # Server host
  read_timeout: 30s               # HTTP read timeout
  write_timeout: 30s              # HTTP write timeout
  
database:
  host: "your-db-host"
  port: 5432
  user: "your-username"
  password: "your-password"
  dbname: "ChatApp"
  sslmode: "require"              # SSL mode: disable, require, verify-ca, verify-full
  timezone: "UTC"                 # Database timezone
  max_idle_conns: 10              # Maximum idle connections
  max_open_conns: 100             # Maximum open connections
  conn_max_lifetime: 3600s        # Connection maximum lifetime

jwt:
  secret: "your-super-secret-jwt-key-change-this-in-production"
  expire_hours: 24                # Token expiration in hours
  issuer: "chatapp"               # JWT issuer

websocket:
  read_buffer_size: 1024          # WebSocket read buffer size
  write_buffer_size: 1024         # WebSocket write buffer size
  read_deadline: 60s              # Read deadline
  write_deadline: 10s             # Write deadline
  ping_period: 54s                # Ping period

cors:
  allowed_origins:                # CORS allowed origins
    - "*"
  allowed_methods:                # CORS allowed methods
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:                # CORS allowed headers
    - "Origin"
    - "Content-Type"
    - "Authorization"

logging:
  level: "info"                   # Log level: debug, info, warn, error
  format: "json"                  # Log format: json, text
  output: "stdout"                # Log output: stdout, stderr, file

app:
  name: "ChatApp"                 # Application name
  version: "1.0.0"                # Application version
  debug: true                     # Debug mode
  environment: "development"      # Environment: development, staging, production
```

## Environment Variables

You can override any configuration value using environment variables. The variable names are uppercase and use underscores. Examples:

```bash
# Database configuration
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=myuser
export DATABASE_PASSWORD=mypassword
export DATABASE_DBNAME=chatapp

# JWT configuration
export JWT_SECRET=my-super-secret-key
export JWT_EXPIRE_HOURS=48

# Server configuration
export SERVER_PORT=:3000
export APP_DEBUG=false
export APP_ENVIRONMENT=production
```

## Configuration Loading

The configuration is loaded in the following order (later sources override earlier ones):

1. Default values (hardcoded in the application)
2. Configuration file (`config.yaml`)
3. Environment variables

## Usage in Code

```go
// Load configuration
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal("Failed to load configuration:", err)
}

// Access configuration values
fmt.Println("Server port:", cfg.Server.Port)
fmt.Println("Database host:", cfg.Database.Host)
fmt.Println("JWT secret:", cfg.JWT.Secret)

// Get database DSN
dsn := cfg.GetDatabaseDSN()
```

## Configuration for Different Environments

### Development
```yaml
app:
  debug: true
  environment: "development"
logging:
  level: "debug"
database:
  sslmode: "disable"
```

### Production
```yaml
app:
  debug: false
  environment: "production"
logging:
  level: "info"
  format: "json"
database:
  sslmode: "require"
jwt:
  secret: "production-secret-key-very-long-and-secure"
```

## Security Considerations

1. **JWT Secret**: Always use a strong, unique secret in production
2. **Database Password**: Use environment variables for sensitive data
3. **SSL Mode**: Use "require" or higher in production
4. **Debug Mode**: Disable debug mode in production
5. **CORS Origins**: Restrict allowed origins in production

## Configuration Validation

The application will validate configuration on startup and fail fast if required values are missing or invalid.

## Docker Configuration

When running in Docker, you can mount a configuration file or use environment variables:

```dockerfile
# Mount config file
COPY config.yaml /app/config.yaml

# Or use environment variables
ENV DATABASE_HOST=db
ENV DATABASE_USER=postgres
ENV JWT_SECRET=my-production-secret
```

## Troubleshooting

1. **Configuration not found**: Ensure `config.yaml` is in the application directory
2. **Invalid YAML**: Check YAML syntax with a validator
3. **Environment variables not working**: Ensure variable names match the expected format
4. **Database connection issues**: Verify database configuration and network connectivity