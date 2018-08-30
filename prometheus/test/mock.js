process.env.BERLIOZ_AGENT_PATH = "ws://127.0.0.1:55555/82d1c32d-19bd-4e8b-a53b-7529e386b7c3"
process.env.BERLIOZ_INFRA = 'mock'
process.env.BERLIOZ_CLUSTER = 'infra'
process.env.BERLIOZ_SECTOR = 'prometheus'


console.log('---------------------------------------------');
require('../src/index');
console.log('---------------------------------------------');

