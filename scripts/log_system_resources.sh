#!/bin/bash
# This script monitors CPU and memory usage of the system and Docker containers

# Specify the log file
log_file="usage_log.txt"

# Clear the log file before starting
#echo "" > "$log_file"

while :
do 
  # Get the current usage of CPU and memory
  cpuUsage=$(top -bn1 | awk '/Cpu/ { print $2}')
  memUsage=$(free -m | awk '/Mem/{print $3}')

  echo "$(date '+%Y-%m-%d %H:%M:%S')" >> "$log_file"
  echo "$(date '+%Y-%m-%d %H:%M:%S')" >> "$log_file"
  
  # Log the usage
  echo "CPU Usage: $cpuUsage%" >> "$log_file"
  echo "Memory Usage: $memUsage MB" >> "$log_file"

  # Get and log the CPU and memory usage for each process
  echo "Individual Process Usage:" >> "$log_file"
  top -bn1 -o %CPU | head -n 17 | tail -n +8 >> "$log_file"
  echo "---------------------------------------------------" >> "$log_file"

  # Get and log the stats of Docker containers
  echo "Docker Container Stats:" >> "$log_file"
  docker stats --no-stream >> "$log_file"
  echo "---------------------------------------------------" >> "$log_file"

  # Sleep for 1 second
  sleep 1
done

