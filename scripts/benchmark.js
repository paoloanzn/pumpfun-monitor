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

const LOG_LEVEL = process.env.LOG_LEVEL || 'info';

const stats = {
    total: 0,
    active: 0,
    messagesReceived: 0,
    errors: 0
};

let messageBuffer = [];
const BATCH_INTERVAL = 5000; // 5 seconds

function log(level, message, data = {}) {
    const levels = ['error', 'warn', 'info', 'debug'];
    if (levels.indexOf(level) <= levels.indexOf(LOG_LEVEL)) {
        console.log(JSON.stringify({
            timestamp: new Date().toISOString(),
            level,
            message,
            ...data
        }));
    }
}

function shouldSample() {
    return Math.random() < 0.01; // 1% sampling rate
}

async function getUUIDs() {
    const response = await fetch(`http://${HOST}:8080/mint-consumers`);
    const data = await response.json();
    if (!data) {
        process.exit(1);
    }
    return data?.uuids;
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
            stats.total++;
        }
    });

    wsConns.forEach((conn, index) => {
        conn.onopen = () => {
            stats.active++;
            log('debug', `Connection ${index} opened`);
        };

        conn.onclose = () => {
            stats.active--;
            log('debug', `Connection ${index} closed`);
        };

        conn.onerror = (error) => {
            stats.errors++;
            log('error', 'Connection error', { connection: index, error: error.message });
        };

        conn.onmessage = function (e) {
            stats.messagesReceived++;
            
            if (shouldSample()) {
                log('debug', 'Message received', {
                    connection: index,
                    message: e.data.toString().substring(0, 100) // Truncate long messages
                });
            }
            
            messageBuffer.push(e.data);
        };
    });

    setInterval(() => {
        log('info', 'Connection statistics', {
            stats,
            messageRate: messageBuffer.length / (BATCH_INTERVAL/1000),
            uniqueMessages: new Set(messageBuffer).size
        });
        messageBuffer = [];
    }, BATCH_INTERVAL);

    process.on('SIGINT', function() {
        console.log("\nGracefully closing all WebSocket connections...");
        wsConns.forEach((conn) => {
            conn.close();
        });
        console.log("All connections closed. Exiting.");
        process.exit(0);
    });
    
    console.log(`Connected to ${uuids.length} WebSocket endpoints. Press Ctrl+C to exit.`);
    console.log(`Total client connections: ${wsConns.length}.`);

    setInterval(() => {}, 5 * 60 * 1000);
}

main();
