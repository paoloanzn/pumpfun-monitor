#!/usr/bin/env node

const HOST = process.env.HOST ?? 'localhost'

const colors = Object.freeze({
    0: "\x1b[31m",
    1: "\x1b[32m",
    2: "\x1b[33m",
    3: "\x1b[34m",
    4: "\x1b[35m",
    5: "\x1b[36m",
    RESET: "\x1b[0m"
})

function log(...args) {
    if (process.env.LOG == 1) {
        console.log(...args)
    }
}

async function getUUIDs() {
    const response = await fetch(`http://${HOST}:8080/mint-consumers`);

    const data = await response.json();

    if (!data) {
        process.exit(1);
    }

    const uuids = data?.uuids

    return uuids
}

async function main() {
    const uuids = await getUUIDs();
    if (!uuids) {
        process.exit(1);
    }


    const wsConns = [];
    Array.from(uuids).forEach(uuid => {
        const clients = process.env.CLIENTS ?? 1
        for (let i = 0; i < clients; i++){
            const ws = new WebSocket(`ws://${HOST}:8080/${uuid}`);
            wsConns.push(ws);
        }
    });

    wsConns.forEach((conn, index) => {
        conn.onmessage = function (e) { log(colors[index]+e.data+colors.RESET) };
    })

    process.on('SIGINT', function() {
        console.log("\nGracefully closing all WebSocket connections...");
        wsConns.forEach((conn, _idx) => {
            conn.close();
        });
        console.log("All connections closed. Exiting.");
        process.exit(0);
    });
    
    console.log(`Connected to ${uuids.length} WebSocket endpoints. Press Ctrl+C to exit.`);
    console.log(`Total client connections: ${wsConns.length}.`)

    setInterval(() => {}, 5 * 60 * 1000);
}

main()