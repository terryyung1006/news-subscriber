import { NextResponse } from 'next/server';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { userId, message, contextReportId } = body;

    // Call the gRPC backend
    // Note: In a real app, you would use the generated gRPC client methods
    // For now, we are assuming the gRPC client has a method exposed or we are proxying
    
    // Since we don't have the full gRPC-web setup in this file yet, 
    // we will mock the call to the backend via a direct fetch if we were running in the same network,
    // OR we should use the grpcClient we have.
    
    // However, looking at grpc-client.ts, it seems to be a placeholder.
    // Let's assume for this step that we are calling the Go backend via HTTP/REST if it supports it,
    // OR we are using a gRPC client.
    
    // Given the complexity of setting up gRPC-Web in one step, 
    // and that the user asked for "simple", I will assume the Go backend ALSO exposes HTTP 
    // or we are using a simple fetch to the Go backend's HTTP port if it has one.
    
    // WAIT: The Go backend is gRPC only (net.Listen tcp).
    // Next.js cannot call gRPC directly easily without a proxy (like Envoy) or gRPC-Web.
    
    // ALTERNATIVE: I will make the Go backend expose a simple HTTP JSON endpoint for this chat feature
    // to make it easier for the Next.js API route to call it.
    
    // BUT, for now, let's assume we are just passing the request through.
    // I will implement a temporary mock response here that simulates the backend delay
    // to test the loading state, while we fix the backend connection.
    
    // ACTUALLY, the user wants the full flow.
    // I will implement a simple HTTP client here to talk to a new HTTP endpoint I will add to the Go backend.
    
    const backendUrl = process.env.BACKEND_URL || 'http://localhost:50051'; // This is gRPC port
    // We need an HTTP port on the backend. Let's add one on 8080.
    
    const res = await fetch('http://localhost:8081/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, message }),
    });
    
    if (!res.ok) {
        const err = await res.text();
        return NextResponse.json({ error: err }, { status: res.status });
    }
    
    const data = await res.json();
    return NextResponse.json(data);

  } catch (error) {
    console.error('Chat error:', error);
    return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }
}
