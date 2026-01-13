import { NextRequest, NextResponse } from 'next/server';
import { googleLogin } from '@/lib/grpc-client';

export async function POST(request: NextRequest) {
  try {
    const { idToken } = await request.json();

    if (!idToken) {
      return NextResponse.json(
        { error: 'ID token is required' },
        { status: 400 }
      );
    }

    // Call the gRPC backend
    const response = await googleLogin(idToken);

    // Map LoginStatus enum: 0=UNKNOWN, 1=LOGGED_IN, 2=NEEDS_INVITE
    // The gRPC client might return the enum as a string or a number depending on configuration
    let status = 'UNKNOWN';
    if (response.status === 1 || response.status === 'LOGGED_IN') {
      status = 'LOGGED_IN';
    } else if (response.status === 2 || response.status === 'NEEDS_INVITE') {
      status = 'NEEDS_INVITE';
    }

    return NextResponse.json({
      status,
      sessionToken: response.session_token || '',
      userId: response.user_id || '',
      email: response.email || '',
    });
  } catch (error: any) {
    console.error('Google login error:', error);
    
    // Check if it's an authentication error
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
