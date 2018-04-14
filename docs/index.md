---
layout: default
---

# Getting Started

Download command beat from [GitHub Releases page](https://github.com/uphy/commandbeat/releases/tag/{{ page.version }}).
You can see [Metricbeat installation](https://www.elastic.co/guide/en/beats/metricbeat/current/metricbeat-installation.html) guide as a reference.

Also you can install using docker.

At first, create your config file.

```bash
$ mkdir -p config
$ cat << EOF > config/commandbeat.yml
commandbeat:
  tasks:
    example:
      command: echo "hello world"
output:
  elasticsearch:
    hosts:
      - "YOUR_ELASTICSEARCH_HOST:9200"
EOF
```

Start commandbeat with the above config file.

```bash
$ docker run -it --rm -v "$(pwd)/config:/etc/commandbeat" uphy/commandbeat:{{ page.version }}
```

# Config file format

## Define the tasks

Task defines a command execution settings, command itself and scheduling.
Schedule format is based on cron.

```yaml
commandbeat:
  tasks:
    task1:
      command: echo "hello world"
      schedule: "@every 10s"
    task2:
      command: uptime
      schedule: "0 * * * *"
```

All commands are treated as shell script.  You can write shell script directly in command.

```yaml
commandbeat:
  tasks:
    task2:
      command: |
        #!/bin/bash

        if [ $(grep "hello" somefile > /dev/null;$?) == "0" ]; then
          echo hello exist
        fi
      schedule: "@every 10s"
```

## Parsing the command output

Commandbeat can parse the command output(stdout) and split fields.

### CSV Parser

```yaml
commandbeat:
  tasks:
    parser-csv:
      command: echo "hello,world,1,1.5,true"
      parser:
        type: csv
        fields:
          - name: field1
            type: string
          - name: field2
            type: string
          - name: field3
            type: int
          - name: field4
            type: float
          - name: field5
            type: bool
```

Multiple line csv also available.

```yaml
commandbeat:
  tasks:
    parser-csv-multi:
      command: sh -c 'echo "2018/03/10 11:17:10,hello\n2018/03/10 11:17:12,hello"'
      schedule: "@every 10s"
      parser:
        type: csv
        fields:
          - name: "@timestamp"
            type: timestamp
            format: 2006/01/02 15:04:05
          - name: message
            type: string
```

### JSON Parser

```yaml
commandbeat:
  tasks:
    parser-json:
      command: echo '{"message":"hello world", "year":2018}'
      parser:
        type: json
```