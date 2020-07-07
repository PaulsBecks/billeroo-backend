require("dotenv").config();

const { MongoClient } = require("mongodb");

async function dataToAuthor() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const data = await client.db().collection("data").find({}).toArray();

  for (var d in data) {
    const userId = data[d].userId;
    const authors = data[d].authors;
    for (cI in authors) {
      const author = authors[cI];
      if (!author) {
        continue;
      }
      author["userId"] = userId;

      client.db().collection("authors").insertOne(author);
    }
  }
  client.close();
}

dataToAuthor();
