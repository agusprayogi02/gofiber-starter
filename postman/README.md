# Postman Collection Guide

## Import Collection

1. Open Postman
2. Click "Import" button
3. Select `Go-Fiber-Starter-API.postman_collection.json`
4. Import environment files:
   - `Development.postman_environment.json`
   - `Production.postman_environment.json`

## Setup Environment

1. Click the environment dropdown (top right)
2. Select "Development" or "Production"
3. The `base_url` is already configured
4. `access_token` and `refresh_token` will be auto-populated after login

## Using the Collection

### 1. Authentication Flow

**Step 1: Register a new user**
- Open "Authentication" > "Register"
- Click "Send"
- You should receive a success response

**Step 2: Login**
- Open "Authentication" > "Login"
- Click "Send"
- The response will automatically set `access_token` and `refresh_token` in your environment

**Step 3: Use Protected Endpoints**
- All requests with {{access_token}} will now work
- Token is automatically added to Authorization header

### 2. Refresh Token Flow

When your access token expires (after 1 hour):
1. Open "Authentication" > "Refresh Token"
2. Click "Send"
3. New tokens will be automatically saved

### 3. Password Management

**Forgot Password:**
1. Open "Authentication" > "Forgot Password"
2. Enter your email
3. Check console/logs for reset token (or email in production)

**Reset Password:**
1. Copy the reset token from step above
2. Open "Authentication" > "Reset Password"
3. Paste token and enter new password

**Change Password:**
1. Must be logged in (have access_token)
2. Open "Authentication" > "Change Password"
3. Enter old and new password

### 4. Session Management

**View All Active Sessions:**
- Open "Authentication" > "Get Active Sessions"
- Shows all devices/browsers where you're logged in

**Revoke a Specific Session:**
- Open "Authentication" > "Revoke Session"
- Replace `:sessionId` in URL with actual session ID
- That device will be logged out

**Logout All Devices:**
- Open "Authentication" > "Logout All Devices"
- All sessions will be revoked

### 5. Posts CRUD

**Get All Posts:**
- Open "Posts" > "Get All Posts"
- Supports pagination: `?page=1&per_page=10`

**Create Post:**
- Open "Posts" > "Create Post"
- Body is form-data (for file upload support)
- Fill in title, body
- Optionally add image file

**Update Post:**
- Open "Posts" > "Update Post"
- Replace `:id` in URL with actual post ID
- Modify title, body, or image

**Delete Post:**
- Open "Posts" > "Delete Post"
- Replace `:id` with post ID to delete

## Auto-Save Tokens

The collection includes test scripts that automatically:
1. Extract `token` and `refresh_token` from login/refresh responses
2. Save them to environment variables
3. Use them in subsequent requests

You don't need to manually copy tokens!

## Variables

Collection uses these environment variables:

| Variable | Description | Auto-populated |
|----------|-------------|----------------|
| `base_url` | API base URL | No (set manually) |
| `access_token` | JWT access token (1 hour) | Yes (from login) |
| `refresh_token` | Refresh token (30 days) | Yes (from login) |

## Testing Tips

### Sequential Testing
1. Register → Login → Use Protected Endpoints
2. Or: Login → Get Posts → Create Post → Update Post → Delete Post

### Multiple Users
- Duplicate the collection
- Rename to "User 2", etc.
- Each will have separate tokens

### Testing Expiration
1. Login and save tokens
2. Wait 1 hour for access token to expire
3. Try a protected endpoint (should fail)
4. Use Refresh Token endpoint
5. Try again (should work)

## Common Errors

### 401 Unauthorized
- Token expired → Use Refresh Token
- Not logged in → Login first
- Invalid token → Login again

### 404 Not Found
- Replace `:id` or `:sessionId` in URL with actual values
- Check if resource exists

### 400 Bad Request
- Check request body format
- Ensure all required fields are filled
- Validate field formats (email, password length, etc.)

## Production Use

1. Switch to "Production" environment
2. Update `base_url` to your production URL
3. Use real email for password reset/verification
4. Secure your tokens (don't share collections with saved tokens)

## Exporting/Sharing

When sharing this collection:
1. Clear sensitive data from environment
2. Export collection: Collection → ... → Export
3. Export environment: Environment → ... → Export
4. Share JSON files (tokens are empty by default)

## Runner (Automated Testing)

You can run all requests automatically:
1. Click "..." on collection
2. Select "Run collection"
3. Choose requests to run
4. Set delay between requests
5. Click "Run"

This is useful for regression testing!
