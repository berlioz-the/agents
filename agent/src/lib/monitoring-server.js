const _ = require('the-lodash');
const Boom = require("boom");
const Hapi = require('hapi');

class MonitoringServer
{
    constructor(statsProcessor, taskRepository)
    {
        this._taskRepository = taskRepository;
        this._statsProcessor = statsProcessor;
        this._server = new Hapi.Server();
        this._server.connection({ port: process.env.BERLIOZ_LISTEN_PORT_MON, host: '0.0.0.0' });

        this._server.route({
            method: 'GET',
            path: '/metrics',
            handler: (request, reply) => {
                var info = this._statsProcessor.extractInfo()
                return reply(info.data).type(info.contentType)
            }
        });

        this._server.route({
            method: 'GET',
            path: '/tasks',
            handler: (request, reply) => {
                var info = this._taskRepository.extractInfo()
                return reply(info).type('application/json')
            }
        });
    }

    run()
    {
        this._server.start((err) => {
            if (err) {
                throw err;
            }
            console.log(`Monitoring Server running at: ${this._server.info.uri}`);
        });
    }
}

module.exports = MonitoringServer;
