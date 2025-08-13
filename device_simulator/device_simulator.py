import paho.mqtt.client as mqtt
import time
import json
import random
import ssl
import sys
import os

# --- Configuration ---
# MQTT Broker address (read from environment variable, default to localhost)
MQTT_BROKER_HOST = os.getenv("MQTT_BROKER_HOST", "localhost")
MQTT_BROKER_PORT = int(os.getenv("MQTT_BROKER_PORT", 8883))  # Secure TLS port
MQTT_USERNAME = os.getenv("MQTT_USERNAME", "iot_user")
MQTT_PASSWORD = os.getenv("MQTT_PASSWORD", "iot_password")
# Path to the CA certificate (absolute path inside the Docker container due to volume mount)
CA_CERTS_PATH = os.getenv("CA_CERTS_PATH", "/app/emqx_broker/certs/rootCA.crt")

# MQTT topics for different data types
MQTT_TOPIC_TELEMETRY = "iiot/telemetry/machine"
MQTT_TOPIC_PRODUCTION = "iiot/event/production"
MQTT_TOPIC_ALERT = "iiot/event/alert"
MQTT_TOPIC_STATUS_UPDATE = "iiot/event/status_update"  # Topic for status updates

DEVICE_PREFIX = "machine_"
NUM_DEVICES = int(os.getenv("NUM_DEVICES", 3))  # Number of simulated devices
PUBLISH_INTERVAL_TELEMETRY = int(
    os.getenv("PUBLISH_INTERVAL_TELEMETRY", 3)
)  # Telemetry publishing frequency (per device)


# --- MQTT Callback Functions ---
def on_connect(client, userdata, flags, rc, properties):
    """Callback function when the client connects to the MQTT broker."""
    if rc == 0:
        print(f"[{userdata['device_id']}] Connected to MQTT Broker!")
    else:
        print(f"[{userdata['device_id']}] Connection failed, error code {rc}\n")
        sys.exit(1)  # Exit if connection fails


def on_publish(client, userdata, mid, reason_code, properties):
    """Callback function when a message is published."""
    pass  # No need to print for every message in the simulator


def on_disconnect(client, userdata, rc, properties):
    """Callback function when the client disconnects."""
    print(f"[{userdata['device_id']}] Disconnected with result code {rc}")


# --- Main Device Logic ---
def simulate_device(device_id):
    """Simulates the operation of a single IoT device."""
    client = mqtt.Client(
        callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        client_id=device_id,
        userdata={"device_id": device_id},
        protocol=mqtt.MQTTv5,
    )

    # Assign callback functions
    client.on_connect = on_connect
    client.on_publish = on_publish
    client.on_disconnect = on_disconnect

    # Configure TLS
    # try:
    #     client.tls_set(ca_certs=CA_CERTS_PATH, tls_version=ssl.PROTOCOL_TLSv1_2)
    #     # client.tls_set()
    #     # client.tls_insecure_set(True)
    #     client.username_pw_set(MQTT_USERNAME, MQTT_PASSWORD)
    # except FileNotFoundError:
    #     print(
    #         f"Error: CA certificate not found at {CA_CERTS_PATH}. Ensure certificates are created and path is correct."
    #     )
    #     sys.exit(1)
    # except Exception as e:
    #     print(f"TLS configuration error: {e}")
    #     sys.exit(1)

    try:
        # Connect to the MQTT broker
        client.connect(MQTT_BROKER_HOST, MQTT_BROKER_PORT, 60)
        # Start MQTT loop in a non-blocking mode
        client.loop_start()

        print(f"[{device_id}] Starting sensor data simulation...")

        # Store previous machine status to detect changes
        previous_machine_status = "unknown"

        while True:
            timestamp = int(time.time())

            # --- 1. Simulate Machine Telemetry (continuous) ---
            temperature = round(random.uniform(20.0, 35.0), 2)
            humidity = round(random.uniform(40.0, 80.0), 2)
            vibration_hz = round(random.uniform(10.0, 50.0), 2)
            motor_rpm = random.randint(1000, 3000)
            current_amps = round(random.uniform(5.0, 15.0), 2)

            # Current machine status (can change)
            statuses = ["running", "idle", "fault", "maintenance"]
            current_machine_status = random.choices(
                statuses, weights=[0.7, 0.15, 0.1, 0.05], k=1
            )[0]

            error_code = None
            if current_machine_status == "fault":
                error_codes = [
                    "E001_Overheat",
                    "E002_Jammed",
                    "E003_SensorFailure",
                    "E004_PowerLoss",
                ]
                error_code = random.choice(error_codes)

            telemetry_data = {
                "timestamp": timestamp,
                "device_id": device_id,
                "payload_type": "telemetry",
                "data": {
                    "temperature_celsius": temperature,
                    "humidity_percent": humidity,
                    "vibration_hz": vibration_hz,
                    "motor_rpm": motor_rpm,
                    "current_amps": current_amps,
                    "machine_status": current_machine_status,  # Still send in telemetry
                    "error_code": error_code,
                },
            }
            client.publish(
                f"{MQTT_TOPIC_TELEMETRY}/{device_id}",
                json.dumps(telemetry_data),
                qos=2,
                retain=True,
            )
            # print(f"[{device_id}] Telemetry: {json.dumps(telemetry_data)}") # Uncomment to see detailed logs

            # --- 2. Simulate Production Event (occasionally) ---
            if random.random() < 0.2:  # 20% chance to send a production event
                # Device only sends production event, quality status is pending analysis by server
                production_event_data = {
                    "timestamp": timestamp,
                    "device_id": device_id,
                    "payload_type": "event_production",
                    "data": {
                        "event_type": "unit_produced",
                        "product_sku": f"PROD_XYZ_{random.randint(100, 999):03d}",
                        "unit_count": 1,
                        "batch_id": f"BATCH_{random.randint(1000, 9999):04d}",
                        "quality_status": "pending_analysis",  # Set as pending_analysis
                    },
                }
                client.publish(
                    f"{MQTT_TOPIC_PRODUCTION}/{device_id}",
                    json.dumps(production_event_data),
                    qos=2,
                    retain=True,
                )
                print(
                    f"[{device_id}] Production Event: {json.dumps(production_event_data)}"
                )

            # --- 3. Simulate Error/Alert Event (occasionally) ---
            if random.random() < 0.05:  # 5% chance to send an alert
                alert_type = random.choice(
                    [
                        "VIBRATION_EXCEEDED_THRESHOLD",
                        "TEMPERATURE_CRITICAL",
                        "SENSOR_OFFLINE",
                    ]
                )
                alert_data = {
                    "timestamp": timestamp,
                    "device_id": device_id,
                    "payload_type": "event_alert",
                    "data": {
                        "alert_type": alert_type,
                        "severity": random.choice(
                            ["LOW", "MEDIUM", "HIGH", "CRITICAL"]
                        ),
                        "message": f"Alert: {alert_type} detected on {device_id}",
                        "current_value": round(random.uniform(50.0, 100.0), 2),
                    },
                }
                client.publish(
                    f"{MQTT_TOPIC_ALERT}/{device_id}",
                    json.dumps(alert_data),
                    qos=2,
                    retain=True,
                )
                print(f"[{device_id}] Alert Event: {json.dumps(alert_data)}")

            # --- 4. Simulate Machine Status Update (only when status changes) ---
            if current_machine_status != previous_machine_status:
                status_update_data = {
                    "timestamp": timestamp,
                    "device_id": device_id,
                    "payload_type": "event_status_update",
                    "data": {
                        "old_status": previous_machine_status,
                        "new_status": current_machine_status,
                        "reason": "simulated_change",
                    },
                }
                client.publish(
                    f"{MQTT_TOPIC_STATUS_UPDATE}/{device_id}",
                    json.dumps(status_update_data),
                    qos=2,
                    retain=True,
                )
                print(f"[{device_id}] Status Update: {json.dumps(status_update_data)}")
                previous_machine_status = current_machine_status

            time.sleep(PUBLISH_INTERVAL_TELEMETRY)

    except KeyboardInterrupt:
        print(f"\n[{device_id}] Stopping simulation...")
    except Exception as e:
        print(f"[{device_id}] An error occurred: {e}")
    finally:
        client.loop_stop()
        client.disconnect()
        print(f"[{device_id}] Simulator stopped.")


if __name__ == "__main__":
    import threading

    print("Starting device simulator initialization...")
    threads = []
    for i in range(NUM_DEVICES):
        device_id = f"{DEVICE_PREFIX}{i+1:03d}"
        thread = threading.Thread(target=simulate_device, args=(device_id,))
        threads.append(thread)
        thread.start()
        time.sleep(0.5)

    try:
        for thread in threads:
            thread.join()
    except KeyboardInterrupt:
        print("\nStopping all simulators...")
    print("All simulators stopped.")
