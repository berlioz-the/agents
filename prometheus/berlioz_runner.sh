#!/bin/sh

echo "*** STARTING PROMETHEUS" 
/bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_libraries --web.console.templates=/usr/share/prometheus/consoles &
status=$?
echo "*** PROMETHEUS STARTED. STATUS: $status"
pidPrometheus=$!
echo "*** PROMETHEUS STARTED. PID: $pidPrometheus"
if [ $status -ne 0 ]; then
  echo "Failed to start prometheus: $status"
  exit $status
fi

sleep 5

echo "*** STARTING CONFIGURATOR" 
nohup /bin/berlioz_configurator $pidPrometheus &
status=$?
echo "*** CONFIGURATOR STARTED. STATUS: $status"
pidBerlioz=$!
echo "*** CONFIGURATOR STARTED. PID: $pidBerlioz"
if [ $status -ne 0 ]; then
  echo "Failed to start berlioz_configurator: $status"
  exit $status
fi

while sleep 5; do
    # echo "Checking Prometheus. PID=$pidPrometheus..."
    if ! kill -0 $pidPrometheus > /dev/null 2>&1; then
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

