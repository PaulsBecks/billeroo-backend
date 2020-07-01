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
const app = express();

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

app.get("/users", async (req, res) => {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const users = await client.db().collection("users").find().toArray();
  client.close();
  res.json(users);
});

app.post(
  "/login",
  passport.authenticate("local", { session: false }),
  async (req, res) => {
    const { user } = req;
    const token = jwt.sign(req.user, process.env.JWT_SECRET);
    return res.json({ user: { email: user.email }, token });
  }
);

app.post("/register", async (req, res) => {
  const { email, password } = req.body;
  if (!password || !email) {
    return res.status(400).json({});
  }

  const _password = await bcrypt.hash(password, 8);
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  const users = await client.db().collection("users").find({ email }).toArray();

  if (users.length >= 1) {
    return res.status(400).json({ message: "Email bereits registriert." });
  }

  const { insertedId } = await client
    .db()
    .collection("users")
    .insertOne({ email, password: _password });
  client.db().collection("data").insertOne({
    userId: insertedId,
    invoices: [],
    articles: [],
    customers: [],
    authors: [],
  });
  const user = await client
    .db()
    .collection("users")
    .findOne({ _id: insertedId }, { email: 1 });
  const token = jwt.sign(user, process.env.JWT_SECRET);

  //send registration email
  sendRegistrationConfirmation(user);

  return res.json({ user, token });
});

app.post(
  "/data",
  passport.authenticate("jwt", { session: false }),
  async (req, res) => {
    const { user, body } = req;
    const client = new MongoClient(process.env.MONGO_URI);

    await client.connect();
    try {
      const { value } = await client
        .db()
        .collection("data")
        .findOneAndUpdate({ userId: user._id }, { $set: body });

      const data = await client
        .db()
        .collection("data")
        .findOne({ _id: value._id });
      client.close();
      return res.json({ ...data });
    } catch (err) {
      console.log(err);
      client.close();
      return res.status(500).end();
    }
  }
);

app.get(
  "/data",
  passport.authenticate("jwt", { session: false }),
  async (req, res) => {
    const { user } = req;
    const client = new MongoClient(process.env.MONGO_URI);

    await client.connect();
    try {
      const data = await client
        .db()
        .collection("data")
        .findOne({ userId: user._id });

      client.close();
      return res.json({ ...data });
    } catch (err) {
      console.log(err);
      return res.status(500).end();
    }
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

app.post("/email/invoice", (req, res) => {
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
});

app.post("/webhook/:hookId", (req, res) => {
  console.log(Object.keys(req.body));
  const { billing, shipping, line_items, date_paid } = req.body;
  return res.send("Hello world");
});

app.get("/webhook/:hookId", (req, res) => {
  console.log(req.body);
  return res.send("Hello world");
});

app.get("/", (req, res) => {
  return res.send("Hi there!");
});

// START SERVER
const port = 8000;
const server = app.listen(port, () =>
  console.log(`Example app listening at http://localhost:${port}`)
);

app["close"] = () => {
  server.close();
};

module.exports = app;
