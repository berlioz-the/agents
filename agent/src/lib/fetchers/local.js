const _ = require('the-lodash');
const Promise = require('the-promise');

class LocalFetcher
{
    constructor(messageProcessor, server)
    {
        this._messageProcessor = messageProcessor;

        server.route({
            method: "POST", path: "/report",
            config: {
                handler: (request, reply) => {
                    console.log('RECEIVED: ' + JSON.stringify(request.payload, null, 4))
                    this._processData(request.payload);
                    reply();
                },
                payload: {
                    parse: true
                }
            }
        });
    }

    _processData(messages)
    {
        for(var item of messages)
        {
            this._messageProcessor.processTargetMessage(item);
        }
    }

}

module.exports = LocalFetcher;
