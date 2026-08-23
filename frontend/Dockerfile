FROM node:22-alpine

WORKDIR /app/frontend

COPY package.json package-lock.json ./

COPY . /app/frontend

# clean install
RUN npm ci 
RUN npm run build

EXPOSE 3000

CMD ["npm run dev"]