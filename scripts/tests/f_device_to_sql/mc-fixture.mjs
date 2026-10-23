// Isolated MC 3E word-read device. Frames follow internal/protocol/mcprotocol.
import { createServer } from 'node:net';

export async function startMCFixture(port, initialWords) {
  const words = new Map(Object.entries(initialWords).map(([key, value]) => [Number(key), value]));
  const faults = new Set();
  const sockets = new Set();
  const requests = [];
  const server = createServer(socket => {
    sockets.add(socket);
    socket.on('close', () => sockets.delete(socket));
    socket.on('error', () => {});
    let pending = Buffer.alloc(0);
    socket.on('data', chunk => {
      pending = Buffer.concat([pending, chunk]);
      while (pending.length >= 9) {
        const length = pending.readUInt16LE(7);
        if (length > 4096) { socket.destroy(); return; }
        if (pending.length < 9 + length) return;
        const request = pending.subarray(0, 9 + length);
        pending = pending.subarray(9 + length);
        let code = 0, address = 0, count = 0;
        if (request.length !== 21 || request.readUInt16LE(11) !== 0x0401 || request.readUInt16LE(13) !== 0 || request[18] !== 0xa8) {
          code = 0xc059;
        } else {
          address = request.readUIntLE(15, 3);
          count = request.readUInt16LE(19);
          if (count < 1 || count > 125) code = 0xc061;
          if (faults.has(address)) code = 0xc051;
        }
        if (requests.length < 1024) requests.push({ address, count, code });
        const body = Buffer.alloc(2 + (code ? 0 : count * 2));
        body.writeUInt16LE(code);
        if (!code) for (let i = 0; i < count; i++) body.writeUInt16LE(words.get(address + i) ?? 0, 2 + i * 2);
        const header = Buffer.from(request.subarray(0, 9));
        header[0] = 0xd0; header[1] = 0;
        header.writeUInt16LE(body.length, 7);
        socket.write(Buffer.concat([header, body]));
      }
    });
  });
  await new Promise((resolve, reject) => server.once('error', reject).listen(port, '127.0.0.1', resolve));
  return {
    requests, words, faults,
    async close() { for (const socket of sockets) socket.destroy(); await new Promise(resolve => server.close(resolve)); },
  };
}
