const _ = require('the-lodash');
const AWS = require('aws-sdk');
const Promise = require('the-promise');

class AwsFetcher
{
    constructor(messageProcessor)
    {
        this._messageProcessor = messageProcessor;
        this._queueUrl = process.env.BERLIOZ_MESSAGE_QUEUE_BERLIOZ_AGENT;

        var credentials = null;
        if (process.env.BERLIOZ_INFRA == 'aws') {
            console.log('RUNNING IN CONTAINER:');
            credentials = new AWS.ECSCredentials();
        } else if (process.env.BERLIOZ_INFRA == 'local-aws') {
            if (process.env.BERLIOZ_MOCK_AWS_PROFILE) {
                console.log('RUNNING LOCALLY:');
                credentials = new AWS.SharedIniFileCredentials({profile: process.env.BERLIOZ_MOCK_AWS_PROFILE});
            }
        } else {
            return;
        }

        AWS.config.update({
            credentials: credentials,
            region: process.env.BERLIOZ_REGION
        });
        console.log('AWS CONFIG:');
        console.log(JSON.stringify(AWS.config));

        this._sqs = new AWS.SQS({});

        this._run();
    }

    _run()
    {
        this._process()
            .catch(error => {
                console.log('[Processor] there was error: ');
                console.log(error);
                return Promise.timeout(5000)
                    .then(() => this._run());
            });
    }

    _fetchMessages()
    {
        console.log('[_fetchMessages]');
        var params = {
            QueueUrl: this._queueUrl,
            WaitTimeSeconds: 20
        };
        return this._sqs.receiveMessage(params).promise()
            .then(data => {
                return data.Messages;
            });
    }

    _process()
    {
        console.log('[_process]');
        return this._fetchMessages()
            .then(messages => Promise.serial(messages, x => this._processMessage(x)))
            .then(() => console.log('[_process] done. waiting till next round...'))
            .then(() => this._process());
    }

    _processMessage(message)
    {
        console.log('[_processMessage] MessageId:' + message.MessageId);
        // console.log(message);
        var item = JSON.parse(message.Body);
        console.log('[_processMessage] Body: ' + JSON.stringify(item));

        this._messageProcessor.processTargetMessage(item);

        return this._deleteMessage(message);
    }

    _deleteMessage(message)
    {
        console.log('[_deleteMessage] ' + message.MessageId);
        var params = {
            QueueUrl: this._queueUrl,
            ReceiptHandle: message.ReceiptHandle
        };
        return this._sqs.deleteMessage(params).promise()
            .then(data => {
            });
    }

}

module.exports = AwsFetcher;
