# Security Configuration Guide

This document outlines how to properly configure your ChatApp for different environments while keeping sensitive information secure.

## 🔒 Security Best Practices

### 1. Configuration Files

**DO NOT commit sensitive configuration files to version control!**

- ✅ Commit: `config.example.yaml` (template with placeholders)
- ❌ Do not commit: `config.yaml` (contains real credentials)
- ❌ Do not commit: `.env` files with real values
- ✅ Commit: `.env.example` (template)

### 2. Environment-Specific Configurations

#### Development
```yaml
# config.development.yaml
database:
  host: "localhost"
  user: "dev_user"
  password: "dev_password"
  sslmode: "disable"

jwt:
  secret: "development-secret-not-for-production"

app:
  debug: true
  environment: "development"
```

#### Production
```yaml
# config.production.yaml (DO NOT COMMIT THIS FILE)
database:
  host: "your-production-db-host"
  user: "prod_user"
  password: "very-secure-production-password"
  sslmode: "require"

jwt:
  secret: "super-long-and-secure-production-jwt-secret-key-at-least-256-bits"

app:
  debug: false
  environment: "production"
```

### 3. Using Environment Variables

For production deployments, use environment variables instead of config files:

```bash
# Set in your deployment environment
export DATABASE_HOST=your-production-host
export DATABASE_USER=your-user
export DATABASE_PASSWORD=your-secure-password
export JWT_SECRET=your-very-long-and-secure-jwt-secret
export APP_ENVIRONMENT=production
export APP_DEBUG=false
```

### 4. JWT Secret Requirements

- **Development**: Minimum 32 characters
- **Production**: Minimum 64 characters
- Use cryptographically secure random strings
- Never reuse secrets across environments

Generate secure JWT secrets:
```bash
# Using openssl
openssl rand -base64 64

# Using Python
python -c "import secrets; print(secrets.token_urlsafe(64))"

# Using Node.js
node -e "console.log(require('crypto').randomBytes(64).toString('base64'))"
```

### 5. Database Security

#### SSL Configuration
- **Development**: `sslmode: "disable"` (local only)
- **Staging/Production**: `sslmode: "require"` or higher

#### Connection Pooling
```yaml
database:
  max_idle_conns: 10      # Adjust based on your needs
  max_open_conns: 100     # Don't exceed database limits
  conn_max_lifetime: 3600s # 1 hour
```

### 6. CORS Configuration

#### Development
```yaml
cors:
  allowed_origins:
    - "http://localhost:3000"
    - "http://localhost:8080"
```

#### Production
```yaml
cors:
  allowed_origins:
    - "https://yourdomain.com"
    - "https://www.yourdomain.com"
```

**Never use `"*"` in production!**

## 📋 Setup Checklist

### Initial Setup
1. [ ] Copy `config.example.yaml` to `config.yaml`
2. [ ] Update database credentials in `config.yaml`
3. [ ] Generate and set a secure JWT secret
4. [ ] Configure CORS for your domain
5. [ ] Add `config.yaml` to `.gitignore` (already done)

### Before Deploying to Production
1. [ ] Set `app.debug: false`
2. [ ] Set `app.environment: "production"`
3. [ ] Use environment variables for secrets
4. [ ] Enable database SSL
5. [ ] Configure specific CORS origins
6. [ ] Set up proper logging
7. [ ] Review all security settings

### Environment Variables Template
```bash
# Production environment variables
DATABASE_HOST=your-production-db-host
DATABASE_PORT=5432
DATABASE_USER=your-db-user
DATABASE_PASSWORD=your-secure-password
DATABASE_DBNAME=ChatApp
DATABASE_SSLMODE=require

JWT_SECRET=your-very-long-secure-jwt-secret-key
JWT_EXPIRE_HOURS=24

SERVER_PORT=:8080
APP_DEBUG=false
APP_ENVIRONMENT=production

CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

## 🚨 What NOT to Do

1. **Never commit passwords or secrets to version control**
2. **Never use default/weak JWT secrets in production**
3. **Never disable SSL in production**
4. **Never use debug mode in production**
5. **Never use wildcard CORS in production**
6. **Never hardcode credentials in source code**

## 🛡️ Security Monitoring

### Log Security Events
- Failed login attempts
- Invalid JWT tokens
- Database connection failures
- Unusual access patterns

### Regular Security Tasks
- Rotate JWT secrets periodically
- Update dependencies regularly
- Monitor database access logs
- Review CORS configuration

## 📞 Emergency Response

If credentials are accidentally committed:
1. **Immediately** rotate all affected secrets
2. Force push a commit that removes the sensitive data
3. Contact your hosting provider if database credentials were exposed
4. Review access logs for unauthorized usage
5. Consider the entire secret compromised and replace it

## 🔍 Verification

Test your security configuration:

```bash
# Verify environment variables are loaded
curl -H "Authorization: Bearer invalid-token" http://localhost:8080/api/profile

# Should return 401 Unauthorized

# Test CORS (replace with your domain)
curl -H "Origin: http://malicious-site.com" http://localhost:8080/api/chatrooms

# Should be blocked if CORS is properly configured
```

Remember: Security is not a feature, it's a requirement!