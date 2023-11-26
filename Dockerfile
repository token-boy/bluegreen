FROM node:16.14.0-alpine AS builder

WORKDIR /app

RUN npm config set registry https://registry.npm.taobao.org
RUN npm install -g pnpm@8.2.0
COPY package.json pnpm-lock.yaml ./
RUN pnpm i

COPY . .

EXPOSE 80

CMD ["npm", "start"]
