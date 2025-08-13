!#/bin/bash

create_rule \
    "mqtt_to_kafka_telemetry" \
    "SELECT topic, payload, clientid, timestamp FROM \"iiot/telemetry/machine/#\"" \
    "kafka_producer_connector" \
    "oee.telemetry" \
    "${COMMON_MESSAGE_VALUE_TEMPLATE}"

create_rule \
    "mqtt_to_kafka_production" \
    "SELECT topic, payload, clientid, timestamp FROM \"iiot/event/production/#\"" \
    "kafka_producer_connector" \
    "oee.production" \
    "${COMMON_MESSAGE_VALUE_TEMPLATE}"

create_rule \
    "mqtt_to_kafka_alert" \
    "SELECT topic, payload, clientid, timestamp FROM \"iiot/event/alert/#\"" \
    "kafka_producer_connector" \
    "oee.alerts" \
    "${COMMON_MESSAGE_VALUE_TEMPLATE}"

create_rule \
    "mqtt_to_kafka_status_update" \
    "SELECT topic, payload, clientid, timestamp FROM \"iiot/event/status_update/#\"" \
    "kafka_producer_connector" \
    "oee.status_updates" \
    "${COMMON_MESSAGE_VALUE_TEMPLATE}"
