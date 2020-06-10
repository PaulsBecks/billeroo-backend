module.exports = function invoiceEmailTemplate(content) {
  return `<html> 
  <body>
  <p>Hallo,</p>
  <p>dies ist eine automatisch generierte Email von https://billeroo.de. Bei Fragen wenden Sie sich bitte an service@billeroo.de.</p>
  ${content}
  </body>
  </html>`;
};
