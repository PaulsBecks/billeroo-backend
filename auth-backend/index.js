require("dotenv").config();

const express = require("express");
const bcrypt = require("bcryptjs");
const bodyParser = require("body-parser");
const cors = require("cors");
const passport = require("passport");
const jwt = require("jsonwebtoken");
require("./passport")(passport);
const sendRegistrationConfirmation = require("./email/sendRegistrationConfirmation");
const sendInvoice = require("./email/sendInvoice");
const companySceleton = require("./sceletons/company");

const app = express();

function makeid(length) {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

// CORS
app.use(cors());
app.use(express.json({ limit: "50mb" }));

app.use(function (req, res, next) {
  res.header("Access-Control-Allow-Origin", "*");
  res.header(
    "Access-Control-Allow-Headers",
    "Origin, X-Requested-With, Content-Type, Accept"
  );
  next();
});

const { MongoClient } = require("mongodb");

app.use(bodyParser.json());

app.post(
  "/login",
  passport.authenticate("local", { session: false }),
  async (req, res) => {
    const { user } = req;
    const token = jwt.sign(req.user, process.env.JWT_SECRET);
    return res.json({ user: { email: user.email }, token });
  }
);

app.post(
  "/register",
  passport.authenticate("jwt", { session: false }),
  async (req, res) => {
    const { email, password } = req.body;
    const { user: _user } = req;

    if (!password || !email) {
      return res.status(400).json({});
    }

    const _password = await bcrypt.hash(password, 8);
    const client = new MongoClient(process.env.MONGO_URI);
    await client.connect();

    //check if user exists already
    const users = await client
      .db()
      .collection("users")
      .find({ email })
      .toArray();

    if (users.length >= 1) {
      return res.status(400).json({ message: "Email bereits registriert." });
    }

    // update dummy user to use correct email and password
    const { upsertedId } = await client
      .db()
      .collection("users")
      .updateOne(
        { email: _user.email },
        { $set: { email, password: _password, placeholder: false } }
      );

    const user = await client
      .db()
      .collection("users")
      .findOne({ email }, { email: 1, placeholder: true });

    if (!user) {
      console.log("Unable to retrieve user: ", user);
      return res.status(500);
    }

    const token = jwt.sign(user, process.env.JWT_SECRET);

    //send registration email
    sendRegistrationConfirmation(user);

    return res.json({ user, token });
  }
);

app.get(
  "/test/email",
  passport.authenticate("jwt", { session: false }),
  (req, res) => {
    sendRegistrationConfirmation({
      email: "test@billeroo.de",
    });
    res.end();
  }
);

app.post(
  "/email/invoice",
  passport.authenticate("jwt", { session: false }),
  (req, res) => {
    let { data, fileName, text = "", to } = req.body;
    console.log({ text });
    if (!data || !to) {
      res.status(400).end();
      return;
    }

    if (!fileName) {
      fileName = "billeroo_rechnung";
    }

    sendInvoice({ data, to, text, fileName });
  }
);

app.get("/", (req, res) => {
  return res.send("Hi there!");
});

app.get("/users/placeholder", async (req, res) => {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();

  let email;

  while (true) {
    email = makeid(20);
    const users = await client
      .db()
      .collection("users")
      .find({ email })
      .toArray();
    if (users.length < 1) {
      break;
    }
  }

  const password = makeid(20);

  const { insertedId } = await client
    .db()
    .collection("users")
    .insertOne({ email, password, placeholder: true });

  // add company to account
  client
    .db()
    .collection("companies")
    .insertOne({ ...companySceleton, userId: insertedId });

  const user = await client
    .db()
    .collection("users")
    .findOne({ _id: insertedId }, { email: 1, placeholder: 1 });
  const token = jwt.sign(user, process.env.JWT_SECRET);

  return res.json({ user, token });
});

function parsePrice(price) {
  return parseFloat((price + "").replace(",", "."));
}

app.get(
  "/stats",
  passport.authenticate("jwt", { session: false }),
  async function (req, res) {
    const client = new MongoClient(process.env.MONGO_URI);
    await client.connect();
    const { user } = req;
    const now = new Date();
    const { year = now.getFullYear() } = req.query
    const invoices = await client
      .db()
      .collection("invoices")
      .find({
        $or: [{ deleted: { $exists: false } }, { deleted: false }],
        userId: user._id,
        invoiceDate: { '$regex': year }
      })
      .toArray();

    const invoiceStats = [[], [], [], [], [], [], [], [], [], [], [], []];

    for (let i in invoices) {
      const invoice = invoices[i];
      const invoiceDate = new Date(invoice.invoiceDate);
      const { totalPrice } = invoice;
      const month = invoiceDate.getMonth();

      invoiceStats[parseInt(month)].push({
        totalPrice: parsePrice(totalPrice),
        totalPriceNet:
          parsePrice(totalPrice) / (1 + parsePrice(invoice.customer.ust) / 100),
      });
    }

    const articleStats = invoices.reduce((stats, i) => {
      for (const aId in i.articles) {
        const article = i.articles[aId];
        if (!stats[article._id]) {
          stats[article._id] = {
            totalSold: 0,
            invoices: [],
            totalSend: 0,
            totalTurnover: 0,
            totalTurnoverNet: 0,
            name: article.name,
          };
        }
        const toBePayed = parseInt(article.toBePayed + "");
        stats[article._id].totalSold += toBePayed;
        stats[article._id].totalSend += parseInt(article.toBeSend + "");
        stats[article._id].totalTurnover +=
          toBePayed * parsePrice(article.price);
        stats[article._id].totalTurnoverNet +=
          (toBePayed *
            parsePrice(article.price) *
            (100 + parsePrice(i.customer.ust))) /
          100;
        stats[article._id].invoices.push({
          _id: i._id,
          send: article.toBeSend,
          payed: article.toBePayed,
          invoiceNumber: i.invoiceNumber,
          customerName: i.customer.name,
          paymentDate: i.paymentDate,
        });
      }
      return { ...stats };
    }, {});

    return res.json({ body: { invoiceStats, articleStats } });
  }
);
// START SERVER
const port = 8000;
const server = app.listen(port, () =>
  console.log(`Example app listening at http://localhost:${port}`)
);

app["close"] = () => {
  server.close();
};

module.exports = app;
