```mermaid
flowchart TD
A[EntriesLoadedMsg] -->|Name entered| B[StartRunningMsg]
B -->|Space pressed| C[StopRunningMsg]
C --> D[AddEntryMsg]
D --> E[EntryAddedMsg]
```
