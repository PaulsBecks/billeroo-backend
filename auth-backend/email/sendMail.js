"use strict";
const nodemailer = require("nodemailer");

module.exports = async function sendMail(emailData) {
  let transporter = nodemailer.createTransport({
    host: process.env.EMAIL_PROVIDER,
    port: 587,
    secure: false, // true for 465, false for other ports
    auth: {
      user: process.env.EMAIL_USER,
      pass: process.env.EMAIL_PASSWORD,
    },
    tls: {
      ciphers: "SSLv3",
    },
    requireTLS: true,
  });

  let info = await transporter.sendMail(emailData);
  console.log(info);
};
