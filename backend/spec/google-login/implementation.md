# Implementation Specification: Google Login

## Changes to Proto Definitions

**File**: `news-subscriber-grpc-proto/proto/auth.proto`

Replace existing auth methods with:

```protobuf
service AuthService {
  rpc GoogleLogin(GoogleLoginRequest) returns (GoogleLoginResponse);
  rpc CompleteSignup(CompleteSignupRequest) returns (CompleteSignupResponse);
}

enum LoginStatus {
  UNKNOWN = 0;
  LOGGED_IN = 1;
  NEEDS_INVITE = 2;
}

message GoogleLoginRequest {
  string id_token = 1;
}

message GoogleLoginResponse {
  LoginStatus status = 1;
  string session_token = 2;
  string user_id = 3;
  string email = 4;
}

message CompleteSignupRequest {
  string id_token = 1;
  string invite_code = 2;
}

message CompleteSignupResponse {
  string session_token = 1;
  string user_id = 2;
}
```

## Database Schema

**File**: `migrations/001_init.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    google_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS invite_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    used_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    used_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_invite_codes_code ON invite_codes(code);
```

## Backend Implementation

### 1. Data Models

**File**: `src/models/user.go`

```go
type User struct {
    ID        string    `db:"id"`
    Email     string    `db:"email"`
    GoogleID  string    `db:"google_id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type InviteCode struct {
    ID        string     `db:"id"`
    Code      string     `db:"code"`
    UsedBy    *string    `db:"used_by"`
    CreatedAt time.Time  `db:"created_at"`
    UsedAt    *time.Time `db:"used_at"`
}
```

### 2. Repositories

**File**: `src/repository/postgres/user_repository.go`

```go
type UserRepository struct {
    db *DB
}

func (r *UserRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error)
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error)
func (r *UserRepository) Create(ctx context.Context, email, googleID, name string) (*models.User, error)
```

**File**: `src/repository/postgres/invite_repository.go`

```go
type InviteCodeRepository struct {
    db *DB
}

func (r *InviteCodeRepository) IsValid(ctx context.Context, code string) (bool, error)
func (r *InviteCodeRepository) MarkAsUsed(ctx context.Context, code, userID string) error
```

### 3. Service Implementation

**File**: `src/service/auth.go`

```go
func (s *Service) GoogleLogin(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.GoogleLoginResponse, error) {
    // 1. Verify Google ID token
    payload, err := s.verifyGoogleToken(ctx, req.IdToken)
    
    // 2. Check if user exists
    user, err := s.userRepo.GetByGoogleID(ctx, payload.GoogleID)
    
    // 3a. Existing user - generate token and return
    if user != nil {
        token := s.generateSessionToken()
        return &pb.GoogleLoginResponse{
            Status:       pb.LoginStatus_LOGGED_IN,
            SessionToken: token,
            UserId:       user.ID,
            Email:        user.Email,
        }
    }
    
    // 3b. New user - request invite code
    return &pb.GoogleLoginResponse{
        Status: pb.LoginStatus_NEEDS_INVITE,
        Email:  payload.Email,
    }
}

func (s *Service) CompleteSignup(ctx context.Context, req *pb.CompleteSignupRequest) (*pb.CompleteSignupResponse, error) {
    // 1. Verify Google ID token
    payload, err := s.verifyGoogleToken(ctx, req.IdToken)
    
    // 2. Validate invite code
    valid, err := s.inviteRepo.IsValid(ctx, req.InviteCode)
    if !valid {
        return nil, status.Errorf(codes.InvalidArgument, "invalid invite code")
    }
    
    // 3. Create user
    user, err := s.userRepo.Create(ctx, payload.Email, payload.GoogleID, payload.Name)
    
    // 4. Mark invite code as used
    err = s.inviteRepo.MarkAsUsed(ctx, req.InviteCode, user.ID)
    
    // 5. Generate session token
    token := s.generateSessionToken()
    
    return &pb.CompleteSignupResponse{
        SessionToken: token,
        UserId:       user.ID,
    }
}
```

### 4. Configuration

**File**: `config.yaml`

```yaml
database:
  host: "localhost"
  port: 5432
  user: "user"
  password: "password"
  dbname: "news_subscriber"
  sslmode: "disable"

google:
  client_id: "YOUR_GOOGLE_CLIENT_ID"
```

## Frontend Implementation

### Login Page

**File**: `src/app/login/page.tsx`

```tsx
'use client';

export default function LoginPage() {
  const [showInviteInput, setShowInviteInput] = useState(false);
  const [pendingToken, setPendingToken] = useState('');
  
  const handleGoogleCallback = async (response) => {
    const res = await fetch('/api/auth/google-login', {
      method: 'POST',
      body: JSON.stringify({ idToken: response.credential }),
    });
    
    const data = await res.json();
    
    if (data.status === 'LOGGED_IN') {
      localStorage.setItem('session_token', data.sessionToken);
      router.push('/dashboard');
    } else if (data.status === 'NEEDS_INVITE') {
      setPendingToken(response.credential);
      setShowInviteInput(true);
    }
  };
  
  const handleCompleteSignup = async (inviteCode) => {
    const res = await fetch('/api/auth/complete-signup', {
      method: 'POST',
      body: JSON.stringify({ idToken: pendingToken, inviteCode }),
    });
    
    const data = await res.json();
    localStorage.setItem('session_token', data.sessionToken);
    router.push('/dashboard');
  };
  
  return (
    <div>
      {!showInviteInput ? (
        <GoogleSignInButton onCallback={handleGoogleCallback} />
      ) : (
        <InviteCodeInput onSubmit={handleCompleteSignup} />
      )}
    </div>
  );
}
```

### API Routes

**File**: `src/app/api/auth/google-login/route.ts`

```ts
export async function POST(request: NextRequest) {
  const { idToken } = await request.json();
  
  // Call gRPC backend
  const response = await grpcClient.AuthService.GoogleLogin({ idToken });
  
  return NextResponse.json({
    status: response.status === 1 ? 'LOGGED_IN' : 'NEEDS_INVITE',
    sessionToken: response.session_token,
    userId: response.user_id,
    email: response.email,
  });
}
```

**File**: `src/app/api/auth/complete-signup/route.ts`

```ts
export async function POST(request: NextRequest) {
  const { idToken, inviteCode } = await request.json();
  
  // Call gRPC backend
  const response = await grpcClient.AuthService.CompleteSignup({ 
    idToken, 
    inviteCode 
  });
  
  return NextResponse.json({
    sessionToken: response.session_token,
    userId: response.user_id,
  });
}
```

## Docker Setup

**File**: `docker-compose.yaml`

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: news_subscriber_postgres
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: news_subscriber
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Setup Steps

1. **Start Database**: `docker-compose up -d`
2. **Run Migrations**: `psql -U user -d news_subscriber -f migrations/001_init.sql`
3. **Configure Google OAuth**: Get Client ID from Google Cloud Console
4. **Update Configs**: Set `NEXT_PUBLIC_GOOGLE_CLIENT_ID` and `config.yaml`
5. **Start Services**: Run backend and frontend

## Testing

Test invite codes pre-loaded:
- `WELCOME2024`
- `BETA_USER_001`
- `EARLY_ACCESS`

## Future: Transition to Public

To remove invite requirement, modify `GoogleLogin` in `auth.go`:

```go
if user != nil {
    // existing user logic
}

// Instead of returning NEEDS_INVITE, create user directly:
user, err := s.userRepo.Create(ctx, payload.Email, payload.GoogleID, payload.Name)
// ... generate token and return
```
