FROM node:18-alpine

WORKDIR /app

COPY ./package.json ./
COPY ./pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm install -g pnpm && pnpm install

COPY . .

