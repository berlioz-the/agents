const _ = require('the-lodash');
const stats = require('../external/docker-stats')
const through = require('through2')
const PromClient = require('prom-client');
const PromRegister = PromClient.register;

class StatsProcessor
{
    constructor(taskRepository)
    {
        this._taskRepository = taskRepository;
        this._setupCollector();

        this._cpu_usage = new PromClient.Gauge({
            name: 'berlioz_cpu_usage',
            help: 'CPU Usage',
            labelNames: ['service', 'kind', 'identity']
        });
        this._memory_usage = new PromClient.Gauge({
            name: 'berlioz_memory_usage',
            help: 'Memory Usage',
            labelNames: ['service', 'kind', 'identity']
        });
    }

    extractInfo()
    {
        return {
            contentType: PromRegister.contentType,
            data: PromRegister.metrics()
        }
    }

    _setupCollector()
    {
        var opts = {
          docker: null, // here goes options for Dockerode
          events: null, // an instance of docker-allcontainers
        
          statsinterval: 1, // downsample stats. Collect a number of statsinterval logs
                             // and output their mean value
        
          // the following options limit the containers being matched
          // so we can avoid catching logs for unwanted containers
        //   matchByName: '*' // optional
        //   matchByImage: /hello-main-web/, // optional
        //   skipByName: /.*pasteur.*/, // optional
        //   skipByImage: /.*dockerfile.*/ // optional
        }
        stats(opts).pipe(through.obj((chunk, enc, cb) => {
            this._handleStatsChunk(chunk)
            
        //   this.push(JSON.stringify(chunk))
        //   this.push('\n')
        //   this.push('\n')
        //   this.push('\n')
        //   this.push('\n')
          cb()
        }))
        //.pipe(process.stdout)
    }

    _handleStatsChunk(chunk)
    {
        var containerInfo = this._taskRepository.getContainer(chunk.fullId);
        if(!containerInfo) {
            return;
        }

        this._cpu_usage.set(containerInfo.StatsLabels, chunk.stats.cpu_stats.cpu_usage.cpu_percent);
        this._memory_usage.set(containerInfo.StatsLabels, chunk.stats.memory_stats.usage);

        // console.log(containerInfo.TaskName + ': ' + chunk.stats.cpu_stats.cpu_usage.cpu_percent)
        // console.log()
        // console.log(chunk.stats.memory_stats.usage)
        // var objStats = this._getObjectStats(containerInfo.Id);
        // objStats['cpu_percent'] = chunk.stats.cpu_stats.cpu_usage.cpu_percent;
        // objStats['memory_usage'] = chunk.stats.memory_stats.usage;
    }


}

module.exports = StatsProcessor