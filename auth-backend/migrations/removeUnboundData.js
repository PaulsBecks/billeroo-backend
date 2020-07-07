require("dotenv").config();

const { MongoClient } = require("mongodb");

async function removeUnboundData() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const data = await client.db().collection("data").find({}).toArray();

  for (var d in data) {
    const user = await client
      .db()
      .collection("users")
      .findOne({ _id: data[d].userId });

    if (user) {
    } else {
      dt = await client.db().collection("data").removeOne({ _id: data[d]._id });
      console.log(dt);
    }
  }
  client.close();
}

removeUnboundData();
