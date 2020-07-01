FROM node:14.4.0-alpine3.10

WORKDIR /usr/billeroo-backend
COPY ./package.json ./
RUN npm install
COPY ./ ./

EXPOSE 8000

CMD ["node", "index.js"]