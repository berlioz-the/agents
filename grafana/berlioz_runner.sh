#!/bin/sh

echo "*** STARTING GRAFANA" 
/run.sh &
status=$?
echo "*** GRAFANA STARTED. STATUS: $status"
pidGrafana=$!
echo "*** PROMETHEUS STARTED. PID: $pidGrafana"
if [ $status -ne 0 ]; then
  echo "Failed to start prometheus: $status"
  exit $status
fi

sleep 5

echo "*** STARTING CONFIGURATOR" 
nohup /bin/berlioz_configurator $pidGrafana &
status=$?
echo "*** CONFIGURATOR STARTED. STATUS: $status"
pidBerlioz=$!
echo "*** CONFIGURATOR STARTED. PID: $pidBerlioz"
if [ $status -ne 0 ]; then
  echo "Failed to start berlioz_configurator: $status"
  exit $status
fi

while sleep 5; do
    # echo "Checking Prometheus. PID=$pidGrafana..."
    if ! kill -0 $pidGrafana > /dev/null 2>&1; then
        echo "Prometheus died. Exiting."
        exit 1
    fi
    
    # echo "Checking Configurator. PID=$pidBerlioz..."
    if ! kill -0 $pidBerlioz > /dev/null 2>&1; then
        echo "Configurator died. Exiting."
        exit 1
    fi
done

echo "*** DONE ***" 

