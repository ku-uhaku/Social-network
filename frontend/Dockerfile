FROM node:alpine

WORKDIR /app/frontend

COPY package.json package-lock.json ./

COPY . /app/frontend

RUN npm install
RUN npm run build

EXPOSE 8080

CMD ["npm", "start"]