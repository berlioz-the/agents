const _ = require('the-lodash');
const Promise = require('the-promise');
const Docker = require('dockerode')

class TaskRepository
{
    constructor()
    {
        var opts = undefined
        this._docker = new Docker(opts)

        this._containers = {};
        this._activeContainerIds = {};
        this._isScheduled = false;
        this._isRefreshing = false;
        this.refresh();
    }

    refresh()
    {
        return this._doRefresh();        
    }

    extractInfo()
    {
        return this._containers;
    }

    scheduleRefresh()
    {
        // console.log('[scheduleRefresh]')
        if (this._isScheduled) {
            return;
        }
        this._isScheduled = true;
        setTimeout(() => {
            this._doRefresh()
        }, 20000)
    }

    getContainer(id)
    {
        var containerInfo = this._containers[id];
        if (!containerInfo) {
            this._doRefresh();
            return null;
        }
        if (containerInfo.ignore) {
            return null;
        }
        return containerInfo;
    }

    _doRefresh()
    {
        if (this._isRefreshing) {
            return;
        }
        this._isRefreshing = true;
        this._isScheduled = false;
        return this._docker.listContainers()
            .then(containers => {
                this._activeContainerIds = {};
                return Promise.serial(containers, x => this._processContainer(x))
                    .then(() => {
                        for(var id of _.keys(this._containers))
                        {
                            if (!this._activeContainerIds[id]) {
                                delete this._containers[id];
                            }
                        }
                        // console.log(this._containers)
                    });
            })
            .then(() => {
                this._isRefreshing = false;
            })
            .catch(error => {
                this._isRefreshing = false;
                console.log('[TaskRepository] failed to list containers: ' + error);
                this.scheduleRefresh();
            })
            .then(() => this.scheduleRefresh())
            ;
    }

    _processContainer(container)
    {
        this._activeContainerIds[container.Id] = true;

        var containerInfo = this._containers[container.Id];
        if (containerInfo) {
            return;
        }

        console.log('[DOCKER::INSPECT]')

        return this._docker.getContainer(container.Id).inspect()
            .then(data => {
                var envList = this._parseEnvVariables(data); 

                var containerInfo = {
                    Id: container.Id
                };
                if (envList['BERLIOZ_CLUSTER'])
                {
                    containerInfo.ignore = false;
                    containerInfo.Image = container.Image;
                    containerInfo.State = container.State;
    
                    containerInfo.cluster = envList['BERLIOZ_CLUSTER'];
                    containerInfo.sector = envList['BERLIOZ_SECTOR'];
                    containerInfo.service = envList['BERLIOZ_SERVICE'];
                    containerInfo.identity = envList['BERLIOZ_IDENTITY'];
    
                    containerInfo.TaskName = [containerInfo.cluster, containerInfo.sector, containerInfo.service, containerInfo.identity].join('-');
                    containerInfo.StatsLabels = {
                        service: [containerInfo.cluster, containerInfo.sector, containerInfo.service].join('-'),
                        identity: containerInfo.identity,
                        kind: 'task'
                    };
                }
                else 
                {
                    containerInfo.ignore = true;
                }
                this._containers[container.Id] = containerInfo;
            })
            .catch(error => {
                console.log('[TaskRepository] failed to inspect container: ' + error);
                this.scheduleRefresh();
            })
    }

    _parseEnvVariables(containerData)
    {
        var env = {};
        for(var envLine of containerData.Config.Env)
        {
            var index = envLine.indexOf('=');
            env[envLine.substring(0, index)] = envLine.substring(index + 1);
        }
        return env;
    }

}

module.exports = TaskRepository