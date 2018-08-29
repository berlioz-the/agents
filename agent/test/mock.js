process.env.BERLIOZ_INFRA = 'mock'
process.env.BERLIOZ_CLUSTER = 'hello'
process.env.BERLIOZ_SECTOR = 'main'
// process.env.BERLIOZ_REGION = 'us-east-1'
// process.env.BERLIOZ_INSTANCE_ID = 'abcd1234'
// process.env.BERLIOZ_MOCK_AWS_PROFILE = 'croundme'
process.env.BERLIOZ_LISTEN_PORT_WS = 55555
process.env.BERLIOZ_LISTEN_PORT_MON = 55556

const fs = require('fs');
const _ = require('the-lodash');

var messageProcessor = require('../src/index');
var str = fs.readFileSync('./data/agent-message.json');
var data = JSON.parse(str);

console.log('---------------------------------------------');

messageProcessor.processTargetMessage(data);
