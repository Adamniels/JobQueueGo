# JobQueueGo

JobQueueGo är ett distribuerat job queue-system skrivet i Go, som hanterar köade jobb, fördelar dem till lediga workers via WebSocket och returnerar resultatet till servern för lagring och senare visning.

## TO RUN

1. pack your program into a zip file and base64 encode it. You can use the following command to do this:
   zip -r ../testprogram.zip .
   base64 -w 0 ../testprogram.zip > base64.txt

2. send the job to server with:
   jq -n --arg input "$(cat base64.txt)" '{"type": "program_zip", "input": $input}' | \
   curl -X POST http://localhost:8080/job \
    -H "Content-Type: application/json" \
    -d @-

Prerequisites:
need to have a makefile that runs with make run

## Arkitekturöversikt

```text
              +---------------------+
              |    HTTP-klienter    |
              |  (t.ex. web UI)     |
              +----------+----------+
                         |
                         | POST /job
                         v
              +----------+----------+
              |      Server / API   |
              |     (main.go)       |
              +----------+----------+
                         |
                         | Lägg till jobb i kön
                         v
                  +-------------+
                  |  JobQueue   | <----+
                  +-------------+      |
                         |             |
                         |             |
                         v             |
              +----------+----------+  |
              |  matchWorkersWithJobs  |
              +----------+----------+  |
                         |             |
        +----------------+-------------+
        | Skickar jobb via WebSocket   |
        v
   +------------+     +------------+     +------------+
   |  Worker 1  |     |  Worker 2  | ... |  Worker N  |
   +------------+     +------------+     +------------+
        |                   |                    |
        | ← Utför jobbet →  |                    |
        | Skickar resultat tillbaka via WebSocket|
        v
   +------------------+
   |  resultstore.go  |
   +------------------+
```
