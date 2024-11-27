#!/usr/bin/env python3
import requests
import json
import sys

# Configurer l'URL de l'API
API_URL = "https://example.com/your-api-endpoint"

def send_alert_to_api(alert):
    headers = {
        "Content-Type": "application/json",
    }

    try:
        response = requests.post(API_URL, headers=headers, data=json.dumps(alert))
        if response.status_code == 200:
            print("Alert successfully sent to the API")
        else:
            print(f"Failed to send alert. Status code: {response.status_code}")
    except Exception as e:
        print(f"Error sending alert: {e}")

if __name__ == "__main__":
    # Recevoir l'alerte de Wazuh via stdin
    alert = sys.stdin.read()

    # Convertir l'alerte en JSON
    try:
        alert_json = json.loads(alert)
    except json.JSONDecodeError:
        print("Error decoding alert JSON")
        sys.exit(1)

    # Vérifier le niveau de gravité
    if alert_json.get("rule", {}).get("level", 0) >= 10:  # Gravité HIGH
        send_alert_to_api(alert_json)
