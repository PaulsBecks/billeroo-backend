var LocalStrategy = require("passport-local").Strategy;
var bcrypt = require("bcryptjs");
const { MongoClient } = require("mongodb");

module.exports = new LocalStrategy(
  { usernameField: "email", passwordField: "password" },
  async (email, password, done) => {
    const client = new MongoClient(process.env.MONGO_URI);
    await client.connect();
    try {
      const user = await client.db().collection("users").findOne({ email });
      if (user && (await bcrypt.compare(password, user.password))) {
        client.close();
        return done(null, user);
      }
      client.close();
      return done(null, false);
    } catch (err) {
      client.close();
      return done(err);
    }
  }
);
