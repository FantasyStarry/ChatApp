# WebSocket Authentication Update

This document describes the updated WebSocket authentication mechanism for the ChatApp.

## Overview

The WebSocket authentication has been changed from Bearer token authentication during connection establishment to message-based authentication after connection.

## Key Changes

### 1. Connection Establishment
- **Before**: Required `Authorization: Bearer <token>` header during WebSocket handshake
- **After**: No authentication required during initial connection

### 2. Authentication Flow
- **Before**: User authenticated before WebSocket upgrade
- **After**: User sends authentication message after connection established

### 3. Authentication Message Format
```json
{
  "type": "auth",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "chatroomId": 1
}
```

### 4. Authentication Response
On successful authentication:
```json
{
  "type": "auth_success",
  "content": "Authentication successful",
  "timestamp": "2023-12-18T10:30:00Z"
}
```

On failed authentication: Connection is immediately closed.

## Implementation Details

### WebSocket Route
- **URL**: `ws://localhost:8080/api/ws/{chatroom_id}`
- **Authentication**: None required for initial connection
- **Middleware**: Removed authentication middleware

### Client State Management
- Clients start in unauthenticated state
- `isAuthenticated` flag tracks authentication status
- Only authenticated clients can send/receive regular messages

### Message Processing
1. **Auth messages**: Processed immediately, not broadcasted or saved
2. **Regular messages**: Only processed if client is authenticated
3. **Unauthenticated attempts**: Connection closed immediately

### Security Features
- Token validation uses existing JWT validation logic
- Failed authentication results in immediate disconnection
- Auth messages are not stored in database
- Auth messages are not broadcasted to other clients

## Usage Example

### 1. Connect to WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws/1');
```

### 2. Send Authentication Message
```javascript
const authMessage = {
    type: 'auth',
    token: 'your-jwt-token-here',
    chatroomId: 1
};
ws.send(JSON.stringify(authMessage));
```

### 3. Wait for Authentication Response
```javascript
ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    if (data.type === 'auth_success') {
        // Now authenticated, can send regular messages
        console.log('Authentication successful');
    }
};
```

### 4. Send Regular Messages
```javascript
const message = {
    type: 'message',
    content: 'Hello, world!'
};
ws.send(JSON.stringify(message));
```

## Testing

Use the provided `websocket_test.html` file to test the new authentication flow:

1. Open `examples/websocket_test.html` in a browser
2. Start the ChatApp server
3. Get a JWT token by logging in via REST API
4. Connect to WebSocket (no auth required)
5. Send authentication message with the JWT token
6. Send regular messages after successful authentication

## Migration Notes

- **Frontend clients** need to be updated to send auth messages after connection
- **Backend routes** no longer require authentication middleware for WebSocket endpoints
- **Existing JWT tokens** continue to work with the new authentication flow
- **API endpoints** remain unchanged, only WebSocket authentication is affected

## Benefits

1. **Flexibility**: Authentication can be re-attempted without reconnecting
2. **Error Handling**: Better control over authentication failures
3. **Protocol Consistency**: Authentication handled within WebSocket protocol
4. **Debugging**: Easier to debug authentication issues

## Security Considerations

- Connections start unauthenticated but cannot perform actions
- Invalid authentication attempts result in immediate disconnection
- JWT validation logic remains unchanged
- Auth messages are processed securely and not persisted