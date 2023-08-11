# json文件

```json
{
  "web_listen_address": ":30007",
  "web_telemetry_path": "/metrics",
  "web_listen_host": "XXX",
  "component_code": "kafka",
  "clusters": [
    {
      "instance_name": "XXX-a",
      "component_env": "qa",
      "component_instance_id": "XXXX",
      "instance_host": "XXX:9092",
      "instance_id": "XXXX",
      "kafka_version": "0.10.2.1",
      "sasl_enabled": "fasle",
      "enable": true,
      "importance": "3",
      "owner": "other"
    },
    {
      "instance_name": "XXX",
      "component_env": "qa",
      "component_instance_id": "XXX",
      "instance_host": "XXX:9092",
      "instance_id": "XXX",
      "kafka_version": "0.10.2.1",
      "sasl_enabled": "false",
      "sasl_mechanism": "plain",
      "sasl_username": "XXX",
      "sasl_password": "XXX",
      "enable": true,
      "importance": "3",
      "owner": "other"
    }
  ]
}
```