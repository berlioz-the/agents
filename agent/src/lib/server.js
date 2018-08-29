const _ = require('the-lodash');
const Boom = require("boom");
const Hapi = require('hapi');
const HAPIWebSocket = require("hapi-plugin-websocket");

class Server
{
    constructor(registry)
    {
        this._registry = registry;

        this._server = new Hapi.Server();
        this._server.connection({ port: process.env.BERLIOZ_LISTEN_PORT_WS, host: '0.0.0.0' });
        this._server.register(HAPIWebSocket);

        this._server.route({
            method: "GET", path: "/{target}",
            handler: (request, reply) => {
                var target = request.params.target;
                var data = registry.getAll(target);
                return reply(data);
            }
        });

        this._server.route({
            method: "POST", path: "/{target}",
            config: {
                payload: { output: "data", parse: true, allow: "application/json" },
                plugins: {
                    websocket: {
                        initially: true,
                        connect: ({ ctx, ws, req }) => {
                            console.log('_setupRoute :: connect. url=' + req.url);
                            var target = this._getTargetFromUrl(req.url);
                            console.log('_setupRoute :: connect. target=' + target);
                            if (target) {
                                registry.registerSubscriber(target, ws);
                            }
                        },
                        disconnect: ({ ctx, ws, req }) => {
                            console.log('_setupRoute :: disconnect. url=' + req.url);
                            var target = this._getTargetFromUrl(req.url);
                            console.log('_setupRoute :: disconnect. target=' + target);
                            if (target) {
                                registry.unregisterSubscriber(target, ws);
                            }
                        }
                    }
                }
            },
            handler: (request, reply) => {
                var target = request.params.target;
                var data = registry.getAll(target);
                let { initially, ws } = request.websocket();
                if (initially) {
                    ws.send(JSON.stringify(data));
                    return reply().code(204);
                } else {
                    return reply(data);
                }
            }
        });
    }

    get server() {
        return this._server;
    }

    run()
    {
        this._server.start((err) => {
            if (err) {
                throw err;
            }
            console.log(`Server running at: ${this._server.info.uri}`);
        });
    }

    _getTargetFromUrl(url)
    {
        var re = new RegExp('\\/(\\S+)' ,"g");
        var match = re.exec(url);
        if (match) {
            return match[1];
        } else {
            console.log('ERROR: counld not parse url: ' + url);
            return null;
        }
    }
}

module.exports = Server;
