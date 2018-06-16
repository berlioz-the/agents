const _ = require('the-lodash');
const Boom = require("boom");

class Registry
{
    constructor(server)
    {
        this._server = server;
        this._data = {};
        this._subscribers = {};
    }

    set(target, section, value)
    {
        var currValue = this.get(target, section);
        if (_.isEqual(currValue, value)) {
            return;
        }
        console.log('[REGISTRY] Set: ' + target + ' Section: ' + section);
        if (!(target in this._data)) {
            this._data[target] = {};
        }
        this._data[target][section] = value;
        this._notifyToSubscribers(target, section);
    }

    unset(target, section)
    {
        if (target in this._data) {
            delete this._data[target][section];
            this._notifyToSubscribers(target, section);
        }
    }

    get(target, section)
    {
        if (target in this._data) {
            var val = this._data[target][section];
            if (val) {
                return val;
            }
        }
        return {};
    }

    getAll(target)
    {
        if (target in this._data) {
            return this._data[target];
        } else {
            return {};
        }
    }

    _notifyToSubscribers(target, section)
    {
        var subscribers = this._subscribers[target];
        if (!subscribers) {
            return;
        }
        var data = {};
        data[section] = this.get(target, section);
        var dataStr = JSON.stringify(data);
        console.log('Notifying to subscribers of: ' + target + ', section: ' + section + ' Data: ' + dataStr);
        for(var subscriber of subscribers) {
            console.log('Sending to subscriber... ');
            subscriber.send(dataStr);
        }
    }

    registerSubscriber(target, ws)
    {
        console.log('Register subscriber: ' + target);
        if (!(target in this._subscribers)) {
            this._subscribers[target] = [];
        }
        this._subscribers[target].push(ws);
    }

    unregisterSubscriber(target, ws)
    {
        console.log('Unregister subscriber: ' + target);
        let idx = this._subscribers[target].indexOf(ws);
        this._subscribers[target].splice(idx, 1);
        if (this._subscribers[target].length == 0) {
            delete this._subscribers[target];
        }
    }
}

module.exports = Registry;
