#!/usr/bin/env node
const repl = require('repl');
const qmux = require('qmux');
const qrpc = require('qrpc');
let cmdFinished = undefined;

async function readChannel(ch) {
    var linebuf = new Buffer.from([]);
    while (true) {
      var data = await ch.read(1);
      if (data === undefined) {
          console.log("|Server got EOF");
          break;
      }
      if (data.toString('ascii') === "\n") {
        console.log(linebuf.toString('ascii'))
        if (cmdFinished) {
            cmdFinished(null);
        }
        linebuf = new Buffer.from([]);
      } else {
        linebuf = Buffer.concat([linebuf, data]);
      }
    }
}

process.on('SIGINT', function() {
    console.log("Caught interrupt signal");
    process.exit();
});

(async function start() {
    try {
        var conn = await qmux.DialWebsocket("ws://localhost:4243");
    } catch {
        setTimeout(start, 500);
    }
    var session = new qmux.Session(conn);
    var client = new qrpc.Client(session);
    var resp = await client.call("repl");
    if (resp.hijacked === true) {
        readChannel(resp.channel);
        let r = repl.start({prompt: ">> ", eval: (cmd, context, filename, callback) => {
            cmdFinished = callback;
            resp.channel.write(Buffer.from(cmd));
        }});
        r.on('exit', () => {
            if (session) {
                session.close();
            }
            process.exit();
        });
    } else {
      console.log("repl stream not hijacked");
    }
})()
