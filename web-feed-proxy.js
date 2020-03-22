const http = require("http");
const { URL, parse } = require("url");
const superagent = require("superagent");
const cheerio = require("cheerio");
const RSS = require("rss");

const server = new http.Server(async (req, res) => {
  try {
    const { url, is, ts, ls } = parse(req.url, true).query;

    await superagent.get(url).then(response => {
      const $ = cheerio.load(response.text);
      const feed = new RSS({ title: $("title").text() });

      $(is).map(function() {
        feed.item({
          title: $(this)
            .find(ts)
            .text(),
          url: new URL(
            $(this)
              .find(ls)
              .attr("href"),
            url
          ).toString()
        });
      });

      res.setHeader("Content-Type", "application/xml; charset=utf-8");
      res.end(feed.xml());
    });
  } catch (error) {
    res.writeHead(500);
    res.end(error.message);
  }
}).listen(parseInt(process.env["PORT"] || 8080), () => {
  console.log(`Server is listening on port ${server.address().port}`);
});
