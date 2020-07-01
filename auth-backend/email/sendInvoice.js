const sendMail = require("./sendMail");
const invoiceMailTemplate = require("./templates/invoiceEmail");

module.exports = function sendInvoice({ to, text, data, fileName }) {
  sendMail({
    to,
    from: process.env.SERVICE_EMAIL,
    subject: "Billeroo | Neue Rechnung verfügbar",
    text:
      "Hallo,\n\ndies ist eine automatisch generierte Email von https://billeroo.de. Bei Fragen wenden Sie sich bitte an service@billeroo.de.\n\n" +
      text,
    html: invoiceMailTemplate(text),
    attachments: [
      {
        filename: fileName + ".pdf",
        contentType: "application/pdf",
        encoding: "base64",
        path: data,
      },
    ],
  });
};
