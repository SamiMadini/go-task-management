FROM node:18-alpine

WORKDIR /app

# Copy package.json and package-lock.json
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the application
COPY frontend/ ./

# Build the application
RUN npm run build

# Expose the port
EXPOSE 3010

# Start the application
CMD ["npm", "start"] 