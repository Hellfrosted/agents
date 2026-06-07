const http = require("node:http");
const net = require("node:net");

function requestJson(url) {
  return new Promise((resolve, reject) => {
    const req = http.get(url, { timeout: 1000 }, (res) => {
      let body = "";
      res.setEncoding("utf8");
      res.on("data", (chunk) => {
        body += chunk;
      });
      res.on("end", () => {
        try {
          resolve({ ok: res.statusCode >= 200 && res.statusCode < 300, data: JSON.parse(body) });
        } catch (error) {
          reject(error);
        }
      });
    });
    req.on("timeout", () => {
      req.destroy(new Error("request timeout"));
    });
    req.on("error", reject);
  });
}

function requestOk(url) {
  return new Promise((resolve) => {
    const req = http.get(url, { timeout: 1000 }, (res) => {
      const ok = res.statusCode >= 200 && res.statusCode < 300;
      res.resume();
      res.on("end", () => resolve(ok));
    });
    req.on("timeout", () => {
      req.destroy(new Error("request timeout"));
    });
    req.on("error", () => resolve(false));
  });
}

async function isServerReady(origin) {
  try {
    const result = await requestJson(`${origin}/api/config`);
    if (!result.ok || !result.data) return { ready: false, iconReady: false, error: null };
    const iconReady = await requestOk(`${origin}/favicon.ico`);
    return { ready: true, iconReady, error: null };
  } catch (error) {
    return { ready: false, iconReady: false, error };
  }
}

function canBind(host, port) {
  return new Promise((resolve) => {
    const server = net.createServer();
    server.once("error", () => resolve(false));
    server.once("listening", () => {
      server.close(() => resolve(true));
    });
    server.listen(port, host);
  });
}

async function findAvailablePort(host, startPort) {
  for (let port = startPort; port <= 65535; port += 1) {
    if (await canBind(host, port)) return port;
  }
  throw new Error("No available local port found for Goal panel server.");
}

module.exports = {
  findAvailablePort,
  isServerReady,
};
