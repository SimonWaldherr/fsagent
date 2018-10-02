# fsagent

fsagent is a Golang Application to perform various standard actions triggered by various events. FSAgent is highly customizable.

```go get simonwaldherr.de/go/fsagent```

## Config

fsagent can do many things, these can be definend and configurated with json files.

start the fsagent deamon with ```go run fsagent.go config.json``` or compile a binary (```go build```) and run it with ```./fsagent config.json```.

the **config.json** file could look like:
```json
[
  {
    "verbose": true,
    "debounce": true,
    "folder": "/mnt/prod/Server/Transfer/701/%Y/%m/%d/",
    "trigger": "fsevent",
    "match": "^[0-9]+\\.[Tt][Xx][Tt]$",
    "action": [
      {
        "do": "sleep",
        "config": "/application/go/fsagent/sleep01.json",
        "onSuccess": [
          {
            "do": "mail",
            "config": "/application/go/fsagent/mail01.json",
            "onSuccess": [
              {
                "do": "move",
                "config": "/application/go/fsagent/fileSuccess01.json"
              }
            ],
            "onFailure": [
              {
                "do": "move",
                "config": "/application/go/fsagent/fileFailure01.json"
              }
            ]
          }
        ]
      }
    ]
  }
]
```

the corresponding action files could look like:

**sleep01.json**:
```json
{
  "time": 500
}
```

**mail01.json**:
```json
{
  "name": "mail",
  "subject": "Lorem Ipsum",
  "body": "dolor sit amet",
  "from": "notification@company.tld",
  "to": ["example@domain.tld"],
  "cc": ["example2@domain.tld"],
  "bcc": ["example3@domain.tld"],
  "user": "notification",
  "pass": "spring2018",
  "server": "webmail.domain.tld",
  "port": 587
}
```

**fileSuccess01.json**:
```json
{
  "name": "/mnt/prod/Server/Archive/701/$file_%Y%m%d%H%M%S"
}
```

**fileFailure01.json**:
```json
{
  "name": "/mnt/prod/Server/Error/701/$file_%Y%m%d%H%M%S"
}
```

## Trigger

Currently there are two triggers availible, the most important trigger is filesystem event trigger based on [fsnotify](github.com/fsnotify/fsnotify).
If this is not possible (e.g. if you work on mounted drives and the fs event comes from a different system) you can use a ticker as trigger.

Trigger | Info
--------|------
fsevent | file system event based on [fsnotify](github.com/fsnotify/fsnotify)
ticker  | checks for new files at a customizable frequenzy 

## Actions

There are some ready-made actions, but you can easily create others yourself.

Action     | Info
-----------|------
Copy       | creates a copy of a given file at the specified destination
Delete     | removes a file
Move       | moves a file to a new location
Decompress | decompresses a file
Compress   | compresses a file
HttpPostR. | sends the content of a file in a HTTP Post Request Body
SendMail   | sends the file as mail attachment
Sleep      | waits for a specified duration
