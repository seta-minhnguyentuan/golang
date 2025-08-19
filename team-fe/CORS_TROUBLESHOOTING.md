# CORS Error Troubleshooting Guide

## Current Issue: CORS Error on Folders Endpoint

### 🔍 **Root Cause Analysis**

The CORS error occurs when your frontend (running on `http://localhost:5173`) tries to access your asset service (running on `http://localhost:7070`), but the server doesn't allow this cross-origin request.

### ✅ **What I've Fixed**

1. **Updated CORS configuration** in both services to include more common development ports
2. **Added additional headers** that might be required
3. **Added CORS debugging** to the frontend API interceptor

### 🔧 **Service Ports Configuration**

- **User Service**: `http://localhost:8080` (GraphQL + Teams REST API)
- **Asset Service**: `http://localhost:7070` (Folders/Notes REST API)  
- **Frontend**: `http://localhost:5173` (Vite default port)

### 🚀 **Steps to Fix CORS Error**

#### 1. **Restart Your Backend Services** (Important!)

The CORS configuration changes require restarting the services:

```bash
# Stop any running services first
sudo fuser -k 8080/tcp  # Kill user service
sudo fuser -k 7070/tcp  # Kill asset service

# Start User Service
cd user-service
go run cmd/api/main.go
# Should show: Server starting on :8080

# In another terminal, start Asset Service  
cd asset-service
go run cmd/api/main.go
# Should show: Server starting on :7070
```

#### 2. **Verify CORS Headers**

Open browser developer tools and check the network tab:
- Look for **OPTIONS** preflight requests
- Check **Access-Control-Allow-Origin** headers in responses
- Ensure **Authorization** header is in **Access-Control-Allow-Headers**

#### 3. **Test the Services Individually**

```bash
# Test Asset Service directly (should work without CORS)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:7070/api/v1/folders

# Test User Service
curl http://localhost:8080/health
```

#### 4. **Check Frontend Port**

Run your frontend and verify the port:
```bash
cd team-fe
npm run dev
# Note the port it starts on (usually 5173)
```

### 🔍 **Debugging Steps**

#### Check Browser Console

The updated API interceptor will now log:
```
Asset API Request: GET /folders with token: eyJhbGciOiJIUzI1NiIs...
Full request URL: http://localhost:7070/api/v1/folders
Current origin: http://localhost:5173
```

#### Common CORS Error Messages

1. **"Access to fetch at 'http://localhost:7070/...' from origin 'http://localhost:5173' has been blocked by CORS policy"**
   - ✅ Fixed: Updated allowed origins in backend

2. **"Request header field authorization is not allowed by Access-Control-Allow-Headers"**
   - ✅ Fixed: Added Authorization to allowed headers

3. **"Method DELETE is not allowed by Access-Control-Allow-Methods"**
   - ✅ Fixed: Added all HTTP methods

### 🛠 **Alternative Solution: Development Proxy**

If CORS issues persist, you can configure Vite to proxy API requests:

**Update `vite.config.ts`:**
```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api/v1': {
        target: 'http://localhost:7070',
        changeOrigin: true,
      },
      '/user': {
        target: 'http://localhost:8080', 
        changeOrigin: true,
      },
      '/teams': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
})
```

**Then update `src/services/api.ts`:**
```typescript
// Use relative URLs when using proxy
export const USER_SERVICE_URL = '';  // Proxied through Vite
export const ASSET_SERVICE_URL = '/api/v1';  // Proxied through Vite
```

### 🔍 **Current CORS Configuration**

**Both services now allow:**
- **Origins**: `http://localhost:3000`, `http://localhost:5173`, `http://localhost:5174`, `http://localhost:4173`
- **Methods**: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`
- **Headers**: `Origin`, `Content-Type`, `Authorization`, `Accept`, `X-Requested-With`
- **Credentials**: `true` (for JWT tokens)

### 📝 **Quick Test Script**

You can test the API integration using the browser console:

```javascript
// In browser console (after logging in)
import { testAssetAPI } from './src/services/testApi.ts';
await testAssetAPI();
```

### 🆘 **If CORS Still Doesn't Work**

1. **Check browser network tab** for exact error messages
2. **Verify services are running** on correct ports
3. **Try the Vite proxy solution** above
4. **Temporarily disable browser security** (development only):
   ```bash
   # Chrome with disabled security (development only!)
   google-chrome --disable-web-security --user-data-dir="/tmp/chrome_dev"
   ```

### 📞 **Need More Help?**

Share the exact error message from the browser console and I can provide more specific guidance!
