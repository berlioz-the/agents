process.env.BERLIOZ_INFRA = 'local-aws'
process.env.BERLIOZ_REGION = 'us-east-1'
process.env.BERLIOZ_INSTANCE_ID = 'abcd1234'
process.env.BERLIOZ_MOCK_AWS_PROFILE = 'croundme'

require('../main/index')
