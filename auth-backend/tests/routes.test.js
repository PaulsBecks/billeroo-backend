const request = require("supertest");
const app = require("..");

describe("Index route", () => {
  it("should return correct status", async () => {
    const res = await request(app).get("/");
    expect(res.statusCode).toEqual(200);
  });

  afterAll((done) => {
    app.close();
  });
});
