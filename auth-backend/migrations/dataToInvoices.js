require("dotenv").config();

const { MongoClient } = require("mongodb");

async function dataToInvoice() {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const data = await client.db().collection("data").find({}).toArray();

  for (var d in data) {
    const userId = data[d].userId;
    const invoices = data[d].invoices;
    for (cI in invoices) {
      const invoice = invoices[cI];
      if (!invoice) {
        continue;
      }
      invoice["userId"] = userId;

      let customer = await client
        .db()
        .collection("customers")
        .findOne({ userId, name: invoice.customer.name });

      if (!customer) {
        customer = invoice.customer;
      }

      for (var i in invoice.articles) {
        const article = await client
          .db()
          .collection("articles")
          .findOne({ userId, name: invoice.articles[i].name });

        invoice.articles[i]["_id"] = article._id;
      }

      client
        .db()
        .collection("invoices")
        .insertOne({ ...invoice, customer });
    }
  }
  client.close();
}

dataToInvoice();
