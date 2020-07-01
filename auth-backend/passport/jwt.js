var JwtStrategy = require("passport-jwt").Strategy,
  ExtractJwt = require("passport-jwt").ExtractJwt;
var opts = {};
opts.jwtFromRequest = ExtractJwt.fromAuthHeaderAsBearerToken();
opts.secretOrKey = process.env.JWT_SECRET;

const { MongoClient } = require("mongodb");

module.exports = new JwtStrategy(opts, async (jwt_payload, done) => {
  const client = new MongoClient(process.env.MONGO_URI);
  await client.connect();
  try {
    const user = await client
      .db()
      .collection("users")
      .findOne({ email: jwt_payload.email });
    client.close();
    return done(null, user);
  } catch (err) {
    client.close();
    return done(err);
  }
});
