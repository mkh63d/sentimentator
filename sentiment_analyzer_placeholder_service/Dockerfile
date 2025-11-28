FROM node:24-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY tsconfig.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY src ./src

# Build TypeScript
RUN npm run build

# Expose port for health checks (optional)
EXPOSE 3000

# Run the compiled service
CMD ["npm", "start"]
