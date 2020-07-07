require("dotenv").config();

const { MongoClient } = require("mongodb");

async function dataToArticle() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const data = await client.db().collection("data").find({}).toArray();

  for (var d in data) {
    const userId = data[d].userId;
    const articles = data[d].articles;
    for (cI in articles) {
      const article = articles[cI];
      if (!article) {
        continue;
      }
      article["userId"] = userId;

      client.db().collection("articles").insertOne(article);
    }
  }
  client.close();
}

dataToArticle();
