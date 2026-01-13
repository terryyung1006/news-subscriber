import { NextRequest, NextResponse } from 'next/server';
import { completeSignup } from '@/lib/grpc-client';

export async function POST(request: NextRequest) {
  try {
    const { idToken, inviteCode } = await request.json();

    if (!idToken || !inviteCode) {
      return NextResponse.json(
        { error: 'ID token and invite code are required' },
        { status: 400 }
      );
    }

    // Call the gRPC backend
    const response = await completeSignup(idToken, inviteCode);

    return NextResponse.json({
      sessionToken: response.session_token,
      userId: response.user_id,
    });
  } catch (error: any) {
    console.error('Complete signup error:', error);
    
    // Check for specific error types
    if (error.message?.includes('invalid') && error.message?.includes('code')) {
      return NextResponse.json(
        { error: 'Invalid or already used invite code' },
        { status: 400 }
      );
    }
    
    if (error.message?.includes('already exists')) {
      return NextResponse.json(
        { error: 'User already exists' },
        { status: 409 }
      );
    }
    
    if (error.message?.includes('invalid') || error.message?.includes('token')) {
      return NextResponse.json(
        { error: 'Invalid Google token' },
        { status: 401 }
      );
    }
    
    return NextResponse.json(
      { error: error.message || 'Internal server error' },
      { status: 500 }
    );
  }
}
