process.env.BERLIOZ_INFRA = 'mock'
// process.env.BERLIOZ_REGION = 'us-east-1'
// process.env.BERLIOZ_INSTANCE_ID = 'abcd1234'
// process.env.BERLIOZ_MOCK_AWS_PROFILE = 'croundme'

const fs = require('fs');
const _ = require('the-lodash');

var messageProcessor = require('../main/index');
var str = fs.readFileSync('./data/agent-message.json');
var data = JSON.parse(str);

console.log('---------------------------------------------');

messageProcessor.processTargetMessage(data);
