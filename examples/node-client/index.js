const grpc = require('grpc')
const winston = require('winston');
const HttpStreamTransport = require('winston-transport-http-stream')


const logger = winston.createLogger({
    level: 'info',
    format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.json()
    ),
    defaultMeta: { service: 'node-client' },
    transports: [
        new HttpStreamTransport({
            url: 'http://localhost:3000/'
        })
    ],
});

console.time("time")

for (let index = 0; index < 1000; index += 1) {
    logger.info(`Log number ${index}`)

}

console.timeEnd("time")

process.on('SIGTERM', () => {
    console.log('SIGTERM signal received.');
    grpc.closeClient()
    process.exit(0);
});

process.on('SIGINT', () => {
    console.log('SIGINT signal received.');
    grpc.closeClient()
    process.exit(0);
});




