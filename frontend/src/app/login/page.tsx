'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { apiClient } from '@/lib/api-client';

declare global {
  interface Window {
    google?: any;
  }
}

export default function LoginPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [showInviteInput, setShowInviteInput] = useState(false);
  const [inviteCode, setInviteCode] = useState('');
  const [pendingToken, setPendingToken] = useState('');
  const [userEmail, setUserEmail] = useState('');

  useEffect(() => {
    // Load Google Sign-In script
    const script = document.createElement('script');
    script.src = 'https://accounts.google.com/gsi/client';
    script.async = true;
    script.defer = true;
    document.body.appendChild(script);

    script.onload = () => {
      if (window.google) {
        window.google.accounts.id.initialize({
          client_id: process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || 'YOUR_CLIENT_ID',
          callback: handleGoogleCallback,
        });

        window.google.accounts.id.renderButton(
          document.getElementById('googleSignInButton'),
          {
            theme: 'outline',
            size: 'large',
            width: '100%',
            text: 'signin_with',
          }
        );
      }
    };

    return () => {
      const existingScript = document.querySelector('script[src="https://accounts.google.com/gsi/client"]');
      if (existingScript) {
        document.body.removeChild(existingScript);
      }
    };
  }, []);

  const handleGoogleCallback = async (response: any) => {
    console.log('Google callback received', response);
    setIsLoading(true);
    setError('');

    try {
      console.log('Calling googleLogin API...');
      const data = await apiClient.auth.googleLogin(response.credential);
      console.log('API response:', data);

      if (data.status === 'LOGGED_IN') {
        console.log('Status is LOGGED_IN, redirecting...');
        // Store session token
        localStorage.setItem('session_token', data.sessionToken);
        localStorage.setItem('user_id', data.userId);
        router.push('/dashboard');
      } else if (data.status === 'NEEDS_INVITE') {
        console.log('Status is NEEDS_INVITE, showing input...');
        // Show invite code input
        setPendingToken(response.credential);
        setUserEmail(data.email);
        setShowInviteInput(true);
      } else {
        console.warn('Unknown status:', data.status);
        setError('Unknown login status: ' + data.status);
      }
    } catch (err: any) {
      console.error('Login error:', err);
      setError(err.message || 'An error occurred during login');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCompleteSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const data = await apiClient.auth.completeSignup(pendingToken, inviteCode);

      // Store session token
      localStorage.setItem('session_token', data.sessionToken);
      localStorage.setItem('user_id', data.userId);
      router.push('/dashboard');
    } catch (err: any) {
      setError(err.message || 'An error occurred during signup');
    } finally {
      setIsLoading(false);
    }
  };

  if (showInviteInput) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-muted/50 px-4">
        <Card className="w-full max-w-sm">
          <CardHeader>
            <CardTitle className="text-2xl">Enter Invite Code</CardTitle>
            <CardDescription>
              Welcome, {userEmail}! Please enter your invitation code to complete signup.
            </CardDescription>
          </CardHeader>
          <form onSubmit={handleCompleteSignup}>
            <CardContent className="grid gap-4">
              {error && (
                <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
                  {error}
                </div>
              )}
              <div className="grid gap-2">
                <Label htmlFor="invite-code">Invite Code</Label>
                <Input
                  id="invite-code"
                  placeholder="WELCOME2024"
                  value={inviteCode}
                  onChange={(e) => setInviteCode(e.target.value)}
                  required
                />
              </div>
            </CardContent>
            <CardFooter className="flex flex-col gap-2">
              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? 'Processing...' : 'Complete Signup'}
              </Button>
              <Button
                type="button"
                variant="ghost"
                className="w-full"
                onClick={() => {
                  setShowInviteInput(false);
                  setPendingToken('');
                  setInviteCode('');
                }}
              >
                Back to Login
              </Button>
            </CardFooter>
          </form>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/50 px-4">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-2xl">Login</CardTitle>
          <CardDescription>
            Sign in with your Google account to access your personalized news feed.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          {error && (
            <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
              {error}
            </div>
          )}
          <div id="googleSignInButton" className="flex justify-center"></div>
          {isLoading && (
            <div className="text-center text-sm text-muted-foreground">
              Signing in...
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
