'use strict';


const Registry = require('./lib/registry');
const Server = require('./lib/server');
const TaskRepository = require('./lib/task-repository');
const StatsProcessor = require('./lib/stats-processor');
const MonitoringServer = require('./lib/monitoring-server');
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

var registry = new Registry(server);
var messageProcessor = new MessageProcessor(registry);
var server = new Server(registry);
var taskRepository = new TaskRepository();
var statsProcessor = new StatsProcessor(taskRepository);
var monitoringServer = new MonitoringServer(statsProcessor, taskRepository);

var fetcher = null;
if (process.env.BERLIOZ_INFRA == 'aws' || process.env.BERLIOZ_INFRA == 'local-aws') {
    const Fetcher = require('./lib/fetchers/aws');
    fetcher = new Fetcher(messageProcessor);
} else if (process.env.BERLIOZ_INFRA == 'local') {
    const Fetcher = require('./lib/fetchers/local');
    fetcher = new Fetcher(messageProcessor, server.server);
} if (process.env.BERLIOZ_INFRA == 'mock') {
    fetcher = null;
}

server.run();
monitoringServer.run();

module.exports = messageProcessor;
