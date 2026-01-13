# Google Login Architecture

## Component Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         User's Browser                           │
│                    (http://localhost:3000)                       │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ Google OAuth
                         ▼
              ┌──────────────────────┐
              │   Google OAuth API   │
              │  (accounts.google)   │
              └──────────┬───────────┘
                         │ ID Token
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Next.js Frontend                               │
│                 (news-subscriber-frontend)                       │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Pages:                                                   │  │
│  │  • /login        - Google Sign-In + Invite Code          │  │
│  │  • /dashboard    - Main app (after auth)                 │  │
│  └──────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  API Routes (Proxy to gRPC):                             │  │
│  │  • /api/auth/google-login                                │  │
│  │  • /api/auth/complete-signup                             │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP (will be gRPC)
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Go gRPC Backend                               │
│                  (news-subscriber-core)                          │
│                     (localhost:8080)                             │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  gRPC Services:                                           │  │
│  │  • AuthService.GoogleLogin                                │  │
│  │  • AuthService.CompleteSignup                             │  │
│  └──────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Repositories:                                            │  │
│  │  • UserRepository                                         │  │
│  │  • InviteCodeRepository                                   │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
                ┌──────────────────────┐
                │   PostgreSQL DB      │
                │   (localhost:5432)   │
                │                      │
                │  Tables:             │
                │  • users             │
                │  • invite_codes      │
                └──────────────────────┘
```

## Authentication Flow

### Scenario A: Existing User (Direct Login)

```
1. User clicks "Sign in with Google"
   Browser → Google OAuth

2. Google verifies user, returns ID token
   Google → Browser

3. Browser sends ID token to frontend API
   Browser → Frontend (/api/auth/google-login)

4. Frontend proxies to backend via gRPC
   Frontend → Backend (AuthService.GoogleLogin)

5. Backend verifies token with Google
   Backend → Google (verify ID token)

6. Backend checks if user exists in DB
   Backend → PostgreSQL (SELECT ... WHERE google_id = ?)

7. User FOUND → Generate session token
   Backend → Frontend

8. Frontend stores token, redirects to dashboard
   Frontend → Dashboard
```

### Scenario B: New User (Requires Invite Code)

```
1-6. Same as Scenario A

7. User NOT found → return NEEDS_INVITE status
   Backend → Frontend

8. Frontend shows invite code input screen
   Frontend → User

9. User enters invite code
   User → Frontend

10. Frontend sends token + code to backend
    Frontend → Backend (AuthService.CompleteSignup)

11. Backend validates invite code
    Backend → PostgreSQL (SELECT ... WHERE code = ? AND used_by IS NULL)

12. If valid, create user account
    Backend → PostgreSQL (INSERT INTO users ...)

13. Mark invite code as used
    Backend → PostgreSQL (UPDATE invite_codes SET used_by = ?, used_at = NOW())

14. Generate session token
    Backend → Frontend

15. Frontend stores token, redirects to dashboard
    Frontend → Dashboard
```

## Database Schema

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    google_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Invite codes table
CREATE TABLE invite_codes (
    id UUID PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    used_by UUID REFERENCES users(id),  -- NULL = unused
    created_at TIMESTAMP DEFAULT NOW(),
    used_at TIMESTAMP                    -- NULL until used
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_invite_codes_code ON invite_codes(code);
```

## Proto Definitions

```protobuf
service AuthService {
  rpc GoogleLogin(GoogleLoginRequest) returns (GoogleLoginResponse);
  rpc CompleteSignup(CompleteSignupRequest) returns (CompleteSignupResponse);
}

enum LoginStatus {
  UNKNOWN = 0;
  LOGGED_IN = 1;        // Existing user
  NEEDS_INVITE = 2;     // New user
}

message GoogleLoginRequest {
  string id_token = 1;  // From Google OAuth
}

message GoogleLoginResponse {
  LoginStatus status = 1;
  string session_token = 2;  // Only if LOGGED_IN
  string user_id = 3;        // Only if LOGGED_IN
  string email = 4;          // Returned in both cases
}

message CompleteSignupRequest {
  string id_token = 1;     // From Google OAuth
  string invite_code = 2;  // User-entered
}

message CompleteSignupResponse {
  string session_token = 1;
  string user_id = 2;
}
```

## Security Model

### Token Verification
1. Frontend receives ID token from Google
2. Backend verifies with Google's public keys
3. Extracts user info (email, google_id, name)
4. Backend never trusts frontend-provided user data

### Session Management
1. Backend generates secure random token (32 bytes)
2. Frontend stores in localStorage
3. All subsequent requests include session token
4. Backend validates token before processing requests

### Invite Code Protection
- Database constraint: `used_by` can only be set once
- Atomic check-and-set in transaction
- Prevents race conditions on code usage

## State Transitions

```
┌─────────────┐
│   No User   │
│  (Landing)  │
└──────┬──────┘
       │
       │ Click "Sign in with Google"
       ▼
┌─────────────┐
│   Google    │
│   OAuth     │
└──────┬──────┘
       │
       │ Return ID Token
       ▼
┌─────────────────────────┐
│  Backend Check User     │
└──┬───────────────────┬──┘
   │                   │
   │ Exists            │ New User
   ▼                   ▼
┌─────────────┐   ┌──────────────────┐
│  Generate   │   │  Request Invite  │
│   Token     │   │      Code        │
└──────┬──────┘   └────────┬─────────┘
       │                   │
       │                   │ User Enters Code
       │                   ▼
       │          ┌──────────────────┐
       │          │  Validate Code   │
       │          │  Create Account  │
       │          └────────┬─────────┘
       │                   │
       │                   │ Success
       ▼                   ▼
    ┌──────────────────────────┐
    │   Dashboard (Logged In)  │
    └──────────────────────────┘
```

## Environment Configuration

### Backend (config.yaml)
```yaml
google:
  client_id: "YOUR_GOOGLE_CLIENT_ID"

database:
  host: "localhost"
  port: 5432
  user: "user"
  password: "password"
  dbname: "news_subscriber"
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID
GRPC_BACKEND_URL=http://localhost:8080
```

## Transition to Public Access

When ready to remove invite requirement, modify `GoogleLogin` service:

```go
// Current: Return NEEDS_INVITE for new users
if user == nil {
    return &pb.GoogleLoginResponse{
        Status: pb.LoginStatus_NEEDS_INVITE,
        Email:  payload.Email,
    }
}

// Future: Auto-create user
if user == nil {
    user, err = s.userRepo.Create(ctx, payload.Email, payload.GoogleID, payload.Name)
    // handle error
}
token := s.generateSessionToken()
return &pb.GoogleLoginResponse{
    Status:       pb.LoginStatus_LOGGED_IN,
    SessionToken: token,
    UserId:       user.ID,
    Email:        user.Email,
}
```

Simply remove the invite code check and auto-create users.
