FROM node:18-alpine AS deps

WORKDIR /deps

# Copy only package files needed for installing dependencies
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci

FROM node:18-alpine AS builder

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /deps/node_modules ./node_modules
COPY frontend/ ./

# Build the application
RUN npm run build

FROM node:18-alpine

WORKDIR /app

# Copy built files and dependencies needed for production
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json

EXPOSE 3010

CMD ["npm", "start"]
