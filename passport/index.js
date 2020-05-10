const local = require("./local");
const jwt = require("./jwt");

module.exports = (passport) => {
  passport.use("local", local);
  passport.use("jwt", jwt);
};
