# Google Login with Invitation Code

## Overview

This feature implements Google OAuth authentication with an invitation code system for controlled access during the beta phase.

## Goals

- **Secure Authentication**: Use Google OAuth for user verification
- **Controlled Access**: Only users with valid invitation codes can sign up
- **Easy Transition**: Simple to remove invite requirement when ready for public launch
- **Database-Backed**: User accounts and invite codes stored in PostgreSQL

## User Flow

1. User clicks "Sign in with Google" on login page
2. Google verifies user identity and returns ID token
3. Backend checks if user exists:
   - **Existing user**: Generate session token → redirect to dashboard
   - **New user**: Prompt for invitation code
4. User enters invitation code
5. Backend validates code (must be unused)
6. Create user account and mark code as used
7. Generate session token → redirect to dashboard

## Technical Approach

- **Frontend**: Google Sign-In JavaScript SDK
- **Backend**: gRPC service with Google token verification
- **Database**: PostgreSQL for users and invite codes
- **Proto**: Centralized definitions in `news-subscriber-grpc-proto`

## Success Criteria

- ✅ Users can sign in with Google accounts
- ✅ New users must provide valid invitation code
- ✅ Invitation codes are single-use only
- ✅ Easy to disable invite requirement later

## Documentation

- **Architecture** → See [`architecture.md`](./architecture.md) for detailed flow diagrams
- **Implementation** → See [`implementation.md`](./implementation.md) for code examples

