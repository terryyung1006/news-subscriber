// This file will handle the communication with the backend (State Machine)
// It currently mocks the responses but should be replaced with real API calls (gRPC-Web or HTTP proxy)

export interface Preference {
  id: string;
  content: string;
}

export interface Memory {
  id: string;
  title: string;
  content: string;
  category: string;
  created_at: string;
}

export interface OnboardingStatus {
  needs_onboarding: boolean;
  has_active_session: boolean;
  session_id?: string;
}

export interface OnboardingStartResponse {
  session_id: string;
  first_message: string;
}

export interface OnboardingMessageResponse {
  response: string;
  is_complete: boolean;
  memories?: Memory[];
}

export interface Report {
  id: string;
  content: string;
  createdAt: string;
}

export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
}

export const apiClient = {
  auth: {
    googleLogin: async (idToken: string) => {
      const res = await fetch('/api/auth/google-login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ idToken }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Login failed');
      return data;
    },
    completeSignup: async (idToken: string, inviteCode: string) => {
      const res = await fetch('/api/auth/complete-signup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ idToken, inviteCode }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Signup failed');
      return data;
    },
  },
  user: {
    getPreferences: async (userId: string): Promise<Preference[]> => {
      // Call UserService.GetPreferences
      return [
        { id: "1", content: "Tech news about AI and LLMs" },
        { id: "2", content: "Global economic trends affecting Asia" },
      ];
    },
    addPreference: async (userId: string, content: string) => {
      // Call UserService.AddPreference
      return { id: Date.now().toString(), content };
    },
  },
  report: {
    getLatest: async (userId: string): Promise<Report> => {
      // Call ReportService.GetLatestReport
      return {
        id: "r1",
        content: "# Daily Digest\n\n...",
        createdAt: new Date().toISOString(),
      };
    },
  },
  chat: {
    sendMessage: async (userId: string, message: string, contextReportId?: string) => {
      const res = await fetch('/api/chat/send', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          // In a real app, you'd attach the session token here
        },
        body: JSON.stringify({ userId, message, contextReportId }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Chat failed');
      return data;
    },
  },
  onboarding: {
    getStatus: async (userId: string): Promise<OnboardingStatus> => {
      const res = await fetch('/api/onboarding/status', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to get status');
      return data;
    },
    start: async (userId: string): Promise<OnboardingStartResponse> => {
      const res = await fetch('/api/onboarding/start', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to start onboarding');
      return data;
    },
    sendMessage: async (userId: string, sessionId: string, message: string): Promise<OnboardingMessageResponse> => {
      const res = await fetch('/api/onboarding/message', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, session_id: sessionId, message }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to send message');
      return data;
    },
  },
  memories: {
    getAll: async (userId: string): Promise<{ memories: Memory[] }> => {
      const res = await fetch(`/api/memories?userId=${userId}`);
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Failed to get memories');
      return data;
    },
  },
};
