require("dotenv").config();

const { MongoClient } = require("mongodb");

async function dataToCustomer() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const data = await client.db().collection("data").find({}).toArray();

  for (var d in data) {
    const userId = data[d].userId;
    const customers = data[d].customers;
    for (cI in customers) {
      const customer = customers[cI];
      if (!customer) {
        continue;
      }
      customer["userId"] = userId;

      client.db().collection("customers").insertOne(customer);
    }
  }
  client.close();
}

dataToCustomer();
