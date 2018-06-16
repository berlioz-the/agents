const _ = require('the-lodash');
const Promise = require('the-promise');

class MessageProcessor
{
    constructor(registry)
    {
        this._registry = registry;
    }

    processTargetMessage(item)
    {
        for(var sectionName of _.keys(item.metadata)) {
            var sectionData = item.metadata[sectionName];
            this._processTargetSection(item.id, sectionName, sectionData)
        }
    }

    _processTargetSection(targetId, sectionName, sectionData)
    {
        console.log('[_processTargetSection] ' + targetId + ' :: ' + sectionName + ' : ' + sectionData);
        this._registry.set(targetId, sectionName, sectionData);
    }

}

module.exports = MessageProcessor;
