const sendMail = require("./sendMail");

module.exports = function sendRegistrationConfirmation(user) {
  const emailData = {
    to: user.email,
    from: process.env.SERVICE_EMAIL,
    bcc: process.env.SERVICE_EMAIL,
    subject: "Willkommen auf Billeroo",
    text:
      "Hallo, \n\n du hast dich gerade auf Billeroo angemeldet.\n Wir hoffe wir können dir die Buchhaltung für deinen Verlag erleichtern. \n \n https://billeroo.de \n\n Mit frendlichen Grüßen, \n Paul Beck",
    html: "",
  };

  sendMail(emailData);
};
