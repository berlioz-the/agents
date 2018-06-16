const _ = require('the-lodash');
const Boom          = require("boom")

class Registry
{
    constructor(server, sections)
    {
        this._sections = sections;
        this._server = server;
        this._data = {};
        this._subscribers = {};
        for(var section of this._sections) {
            this._subscribers[section] = {};
            this._data[section] = {};
            this._setupRoute(section);
        }
    }

    set(target, section, value)
    {
        var currValue = this.get(target, section);
        if (_.isEqual(currValue, value)) {
            return;
        }
        console.log('[REGISTRY] Set: ' + target + ' Section: ' + section);
        this._data[section][target] = value;
        this._notifyToSubscribers(section, target);
    }

    unset(target, section)
    {
        delete this._data[section][target];
        this._notifyToSubscribers(section, target);
    }

    get(target, section)
    {
        if (target in this._data[section]) {
            return this._data[section][target];
        }
        return {};
    }

    _notifyToSubscribers(section, target)
    {
        var subscribers = this._subscribers[section][target];
        if (!subscribers) {
            return;
        }
        var data = this.get(target, section);
        console.log('Notifying to subscribers of: ' + target + ' :: ' + section + ' Data: ' + JSON.stringify(data));
        for(var subscriber of subscribers) {
            console.log('Sending to subscriber... ');
            subscriber.send(JSON.stringify(data));
        }
    }

    _registerSubscriber(section, target, ws)
    {
        console.log('Register subscriber: ' + target + ' to section: ' + section);
        if (!(target in this._subscribers[section])) {
            this._subscribers[section][target] = [];
        }
        this._subscribers[section][target].push(ws);
    }

    _unregisterSubscriber(section, target, ws)
    {
        console.log('Unregister subscriber: ' + target + ' from section: ' + section);
        let idx = this._subscribers[section][target].indexOf(ws);
        this._subscribers[section][target].splice(idx, 1);
        if (this._subscribers[section][target].length == 0) {
            delete this._subscribers[section][target];
        }
    }

    _getTargetFromUrl(url, section)
    {
        var re = new RegExp('\\/(\\S+)\\/' + section ,"g");
        var match = re.exec(url);
        if (match) {
            return match[1];
        } else {
            console.log('ERROR: counld not parse url: ' + url + ' for section: ' + section);
            return null;
        }
    }

    _setupRoute(section)
    {
        this._server.route({
            method: "POST", path: "/{target}/" + section,
            config: {
                payload: { output: "data", parse: true, allow: "application/json" },
                plugins: {
                    websocket: {
                        initially: true,
                        connect: ({ ctx, ws, req }) => {
                            console.log('_setupRoute :: connect. url=' + req.url);
                            var target = this._getTargetFromUrl(req.url, section);
                            console.log('_setupRoute :: connect. target=' + target);
                            if (target) {
                                this._registerSubscriber(section, target, ws);
                            }
                        },
                        disconnect: ({ ctx, ws, req }) => {
                            console.log('_setupRoute :: disconnect. url=' + req.url);
                            var target = this._getTargetFromUrl(req.url, section);
                            console.log('_setupRoute :: disconnect. target=' + target);
                            if (target) {
                                this._unregisterSubscriber(section, target, ws);
                            }
                        }
                    }
                }
            },
            handler: (request, reply) => {
                var target = request.params.target;
                var data = this.get(target, section);

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
}

module.exports = Registry;
