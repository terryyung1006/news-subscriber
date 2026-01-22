import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import path from 'path';

// Load proto file
// Load proto file
// In monorepo, protos are at root/protos/proto
// process.cwd() is root/frontend
const PROTO_PATH = path.join(process.cwd(), '..', 'protos', 'proto', 'auth.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
const authService = protoDescriptor.news_subscriber.v1;

// Create client
const GRPC_SERVER = process.env.GRPC_BACKEND_URL || 'localhost:8080';

let authClient: any = null;

function getAuthClient() {
  if (!authClient) {
    authClient = new authService.AuthService(
      GRPC_SERVER,
      grpc.credentials.createInsecure()
    );
  }
  return authClient;
}

// Helper to promisify gRPC calls
function callAsync(client: any, method: string, request: any): Promise<any> {
  return new Promise((resolve, reject) => {
    // Check if method exists
    if (typeof client[method] !== 'function') {
      // Try camelCase version
      const camelMethod = method.charAt(0).toLowerCase() + method.slice(1);
      if (typeof client[camelMethod] === 'function') {
        method = camelMethod;
      } else {
        // Debug: log available methods
        const availableMethods = Object.getOwnPropertyNames(Object.getPrototypeOf(client))
          .filter(name => typeof client[name] === 'function' && name !== 'constructor');
        reject(new Error(
          `Method '${method}' not found on client. Available methods: ${availableMethods.join(', ')}`
        ));
        return;
      }
    }
    
    client[method](request, (error: any, response: any) => {
      if (error) {
        reject(error);
      } else {
        resolve(response);
      }
    });
  });
}

// API functions
export async function googleLogin(idToken: string) {
  const client = getAuthClient();
  try {
    // Try both PascalCase and camelCase method names
    const response = await callAsync(client, 'googleLogin', {
      id_token: idToken,
    });
    return response;
  } catch (error: any) {
    throw new Error(error.details || error.message || 'Google login failed');
  }
}

export async function completeSignup(idToken: string, inviteCode: string) {
  const client = getAuthClient();
  try {
    // Try both PascalCase and camelCase method names
    const response = await callAsync(client, 'completeSignup', {
      id_token: idToken,
      invite_code: inviteCode,
    });
    return response;
  } catch (error: any) {
    throw new Error(error.details || error.message || 'Signup failed');
  }
}
