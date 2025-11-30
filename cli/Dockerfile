FROM node:18-alpine

WORKDIR /app

# Set npm registry and use host network
RUN npm config set registry https://registry.npmjs.org/

COPY package*.json ./
COPY tsconfig.json ./

# Install with network timeout
RUN npm install --network-timeout=600000

COPY src/ ./src/

RUN npm run build

CMD ["node", "dist/index.js"]