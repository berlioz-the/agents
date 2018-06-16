'use strict';

const Boom = require("boom")
const Hapi = require('hapi');
const HAPIWebSocket = require("hapi-plugin-websocket")

const Registry = require('./lib/registry');
const MessageProcessor = require('./lib/message-processor');

function requireVariable(varName) {
    if (!process.env[varName]) {
        console.log('Error: ' + varName + ' variable is not set.');
        return 1;
    }
    console.log(varName + ': ' + process.env[varName]);
}

requireVariable('BERLIOZ_INFRA');
requireVariable('BERLIOZ_REGION');
requireVariable('BERLIOZ_MESSAGE_QUEUE_BERLIOZ_AGENT');

console.log(JSON.stringify(process.env));

const server = new Hapi.Server();
server.connection({ port: 55555, host: '0.0.0.0' });

server.register(HAPIWebSocket);

var registry = new Registry(server, ['peers', 'endpoints']);
var messageProcessor = new MessageProcessor(registry);

var fetcher = null;
if (process.env.BERLIOZ_INFRA == 'aws' || process.env.BERLIOZ_INFRA == 'local-aws') {
    const Fetcher = require('./lib/fetchers/aws');
    fetcher = new Fetcher(messageProcessor);
} else if (process.env.BERLIOZ_INFRA == 'local') {
    const Fetcher = require('./lib/fetchers/local');
    fetcher = new Fetcher(messageProcessor, server);
} if (process.env.BERLIOZ_INFRA == 'mock') {
    fetcher = null;
}

server.start((err) => {
    if (err) {
        throw err;
    }
    console.log(`Server running at: ${server.info.uri}`);
});

module.exports = messageProcessor;
