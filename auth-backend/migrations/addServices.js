require("dotenv").config();

const { MongoClient } = require("mongodb");

async function dataToArticle() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const invoices = await client
    .db()
    .collection("invoices")
    .updateMany({ services: { $exists: false } }, { $set: { services: [] } });
  client.close();
}

dataToArticle();
