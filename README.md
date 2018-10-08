# fsagent

fsagent is a Golang Application to perform various standard actions triggered by various events. FSAgent is highly customizable.

## Name

The name FSAgent was originally the shorthand for File System Agent.  
With the support of additional triggers and data sources outside of file system events in the monitored folders, a new meaning for the letters F and S must be found. 
Currently, I prefer the definition Free Service Agent. 

## Why

**WARNING:** *sad true reality*  

When planning and developing new (open source) applications / systems, you usually use the latest technologies (containerization (e.g. Docker), message queues (e.g. RabbitMQ, ZeroMQ, ActiveMQ, ...), databases (e.g. PostgreSQL, MariaDB, Redis, MongoDB, ...), ...), but if you're working for a non-startup Company, you often have to deal with old legacy enterprise applications.  
These applications do not have modern interfaces, many are decades old. The most modern communication channels of these applications are mostly FTP uploads and emails.
I really mean FTP, not those new and fancy SFTP Servers.  
However, many of these applications do not even have the functionality to upload, but can only store files in directories and need other applications such as [Bat](https://en.wikipedia.org/wiki/The_Bat!) or [Outlook](https://en.wikipedia.org/wiki/Microsoft_Outlook) to upload. 
Bat is not fundamentally bad, but if your whole business depends on software like Bat, you have a big problem.  
In my free time, I have written this program for replacing such "interfaces". It is not just a replacement for Bat or Outlook, it monitors directories and executes predefined actions for new files. 
The configuration is kept as simple as possible (and will become even easier) to be done by anyone in IT departments, not just programmers.  
The first goal was the elimination of the biggest pain in the ass. 
Gradually, however, it is also planned to extend fsagent for the service composition/orchestration of other protocols such as [HTTP(S)](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol), [AMQP](https://en.wikipedia.org/wiki/Advanced_Message_Queuing_Protocol), [WebSub (PubSubHubbub)](https://en.wikipedia.org/wiki/WebSub).  

## Install

fsagent can easily installed by the ```go get```-command:

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

### Trigger

Currently there are two triggers availible, the most important trigger is filesystem event trigger based on [fsnotify](github.com/fsnotify/fsnotify).
If this is not possible (e.g. if you work on mounted drives and the fs event comes from a different system) you can use a ticker as trigger.

Trigger | Info
--------|------
fsevent | file system event based on [fsnotify](github.com/fsnotify/fsnotify)
ticker  | checks for new files at a customizable frequenzy 

### Actions

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

## Todo / Contribute

Informations about the [license](https://github.com/SimonWaldherr/fsagent/blob/master/LICENSE), [how to contribute](https://github.com/SimonWaldherr/fsagent/blob/master/CONTRIBUTING.md) and a [list of improvements to do](https://github.com/SimonWaldherr/fsagent/blob/master/TODO.md) are in separate files.
