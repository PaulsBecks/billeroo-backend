require("dotenv").config();

const { MongoClient } = require("mongodb");

async function dataToArticle() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const data = await client.db().collection("articles").find({}).toArray();

  for (var d in data) {
    const userId = data[d].userId;
    const authors = data[d].authors;
    for (cI in authors) {
      let author = authors[cI];
      if (!author) {
        continue;
      }

      author = await client
        .db()
        .collection("authors")
        .findOne({ name: author.name });

      if (!author) {
        console.log(data);
        continue;
      }

      const articleAuthor = {
        userId,
        authorId: author._id,
        articleId: data[d]._id,
      };

      client.db().collection("articleAuthors").insertOne(articleAuthor);
    }
  }
  client.close();
}

dataToArticle();
