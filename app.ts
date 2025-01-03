import { connect } from "@db/redis";

const r = await connect({
  hostname: "217.69.5.134",
  port: "6379",
  password: "cxbhna89yr389bfwf74wtg69",
});

const version = Deno.env.get("HOSTNAME");
await r.set("version", version as string);

Deno.serve(async (req) => {
  const url = new URL(req.url);

  if (url.pathname === "/healthcheck") {
    const newVersion = await r.get("version");
    if (newVersion === version) {
      return new Response(null, {
        status: 200,
      });
    } else {
      return new Response(null, {
        status: 503,
      });
    }
  }

  if (req.headers.get("upgrade") == "websocket") {
    const { socket, response } = Deno.upgradeWebSocket(req);
    socket.addEventListener("open", () => {
      console.log("a client connected!");
    });
    socket.addEventListener("message", () => {
      socket.send(version as string);
    });
    return response;
  }

  await new Promise((resolve) => setTimeout(resolve, 5000));

  console.log(url.searchParams.get("v"));
  3;
  return new Response(version, {
    status: 200,
  });
});
