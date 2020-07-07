require("dotenv").config();

const { MongoClient } = require("mongodb");

async function dataToArticle() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const data = await client.db().collection("data").find({}).toArray();

  for (var d in data) {
    const userId = data[d].userId;
    const company = data[d].company;

    company["userId"] = userId;

    client.db().collection("companies").insertOne(company);
  }
  client.close();
}

dataToArticle();
