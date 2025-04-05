# JobQueueGo

JobQueueGo är ett distribuerat job queue-system skrivet i Go, som hanterar köade jobb, fördelar dem till lediga workers via WebSocket och returnerar resultatet till servern för lagring och senare visning.

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

