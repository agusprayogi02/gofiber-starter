# API Authentication Documentation

## Endpoints Authentication & Authorization

### 1. Register
**POST** `/api/register`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "user"  // Optional: "admin" atau "user", default: "user"
}
```

**Response (201):**
```json
{
  "code": 201,
  "message": "User registered successfully",
  "timestamp": "2025-12-31T10:00:00Z"
}
```

---

### 2. Login
**POST** `/api/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Login Success",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2025-12-31T09:00:00Z",
      "updated_at": "2025-12-31T09:00:00Z"
    },
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": "2025-12-31T10:00:00Z"
}
```

**Note:**
- `token` expires in 1 hour
- `refresh_token` expires in 30 days

---

### 3. Refresh Token
**POST** `/api/refresh-token`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": "2025-12-31T11:00:00Z"
}
```

---

### 4. Logout
**POST** `/api/logout` 🔒

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Logout successful",
  "timestamp": "2025-12-31T12:00:00Z"
}
```

---

### 5. Logout All Devices
**POST** `/api/logout-all` 🔒

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Logged out from all devices",
  "timestamp": "2025-12-31T12:00:00Z"
}
```

---

### 6. Forgot Password
**POST** `/api/forgot-password`

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "If email exists, password reset link has been sent",
  "timestamp": "2025-12-31T12:00:00Z"
}
```

**Note:** Email berisi token reset password (implementasi email service diperlukan)

---

### 7. Reset Password
**POST** `/api/reset-password`

**Request Body:**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_password": "newpassword123"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Password reset successful",
  "timestamp": "2025-12-31T12:30:00Z"
}
```

**Note:** Semua sessions user akan di-revoke setelah reset password

---

### 8. Change Password
**POST** `/api/change-password` 🔒

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "old_password": "password123",
  "new_password": "newpassword123"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Password changed successfully",
  "timestamp": "2025-12-31T13:00:00Z"
}
```

**Note:** Semua sessions lain akan di-revoke setelah change password

---

### 9. Verify Email
**POST** `/api/verify-email`

**Request Body:**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Email verified successfully",
  "timestamp": "2025-12-31T13:30:00Z"
}
```

---

### 10. Resend Verification Email
**POST** `/api/resend-verification`

**Request Body:**
```json
{
  "email": "john@example.com"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Verification email sent",
  "timestamp": "2025-12-31T14:00:00Z"
}
```

---

### 11. Get Active Sessions
**GET** `/api/sessions` 🔒

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Active sessions retrieved",
  "data": [
    {
      "id": 1,
      "device_id": "",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2025-12-31T10:00:00Z",
      "expires_at": "2026-01-30T10:00:00Z"
    },
    {
      "id": 2,
      "device_id": "",
      "ip_address": "192.168.1.2",
      "user_agent": "PostmanRuntime/7.29.2",
      "created_at": "2025-12-31T11:00:00Z",
      "expires_at": "2026-01-30T11:00:00Z"
    }
  ],
  "timestamp": "2025-12-31T14:30:00Z"
}
```

---

### 12. Revoke Session
**DELETE** `/api/sessions/:sessionId` 🔒

**Headers:**
```
Authorization: Bearer <access_token>
```

**Path Parameters:**
- `sessionId` (integer): ID of the session to revoke

**Response (200):**
```json
{
  "code": 200,
  "message": "Session revoked successfully",
  "timestamp": "2025-12-31T15:00:00Z"
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "code": 400,
  "message": "Email already exists",
  "order": "S1",
  "timestamp": "2025-12-31T15:00:00Z"
}
```

### 401 Unauthorized
```json
{
  "code": 401,
  "message": "Invalid refresh token",
  "order": "S1",
  "timestamp": "2025-12-31T15:00:00Z"
}
```

### 404 Not Found
```json
{
  "code": 404,
  "message": "User not found",
  "order": "S1",
  "timestamp": "2025-12-31T15:00:00Z"
}
```

### 422 Unprocessable Entity
```json
{
  "code": 422,
  "message": "Validation error",
  "order": "H1",
  "timestamp": "2025-12-31T15:00:00Z"
}
```

### 500 Internal Server Error
```json
{
  "code": 500,
  "message": "Internal server error",
  "order": "S3",
  "timestamp": "2025-12-31T15:00:00Z"
}
```

---

## Notes

🔒 = Requires authentication (Bearer token in Authorization header)

### Token Lifetimes:
- **Access Token**: 1 hour
- **Refresh Token**: 30 days
- **Password Reset Token**: 1 hour
- **Email Verification Token**: 24 hours

### Security Features:
1. Refresh tokens are stored in database and can be revoked
2. Each login creates a new session (refresh token)
3. Sessions track IP address and User-Agent for security monitoring
4. Password reset and change password revoke all active sessions
5. Email verification status is tracked in user table
6. All tokens expire and can be invalidated

### TODO:
- [ ] Implement email sending service for:
  - Email verification
  - Password reset
  - Account notifications
- [ ] Add rate limiting per user
- [ ] Add CAPTCHA for sensitive operations
- [ ] Implement device fingerprinting for better session tracking
