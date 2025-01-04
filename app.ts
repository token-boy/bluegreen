const version = Deno.env.get("HOSTNAME");

const sockets: WebSocket[] = [];

Deno.addSignalListener("SIGTERM", () => {
  for (const socket of sockets) {
    socket.send("reconnect");
  }
});

Deno.serve(async (req) => {
  const url = new URL(req.url);

  if (req.headers.get("upgrade") == "websocket") {
    const { socket, response } = Deno.upgradeWebSocket(req);
    socket.addEventListener("open", () => {
      console.log("a client connected!");
    });
    socket.addEventListener("message", () => {
      socket.send(version as string);
    });
    sockets.push(socket);
    return response;
  }

  await new Promise((resolve) => setTimeout(resolve, 5000));

  console.log(url.searchParams.get("v"));

  return new Response(version, {
    status: 200,
  });
});
