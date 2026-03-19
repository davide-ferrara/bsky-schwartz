# Project: Schwartz-Values Algorithmic Agent for Bluesky

## 1. Vision

The goal of this project is to implement **Algorithmic Sovereignty** on the AT
Protocol (Bluesky). Instead of a black-box algorithm designed for engagement,
this system provides a transparent, value-driven feed based on **Schwartz's
Theory of Basic Human Values**.

The agent will specifically (for now) focus on **Politics** for simplicity,
filtering and re-ranking content according to a "Constitutional" Markdown file
that defines the ethical and social weights of the algorithm.

---

## 2. Core Architecture

The system is built as a set of Go-based microservices running in Docker
containers:

- **Ingestor & Filter (Golang):** Connects to the **Jetstream** websocket to
  consume the Bluesky firehose. It filters for specific keywords related to
  Politics to optimize resource usage.

- **The Schwartz Agent (LLM):** An AI agent (e.g., Llama3 via Ollama or API
  Model) that analyzes the filtered posts. It uses a **Markdown Constitution**
  as its system prompt to assign scores (0 to 6) to Schwartz values (e.g.,
  Universalism, Benevolence, Power, Hedonism). A score of 0 indicates that the
  value does not exist or is not supported in the tweet (i.e., the tweet
  contains content that contradicts a value). If a tweet does contain content
  that supports the value, the scores range from 1 (i.e., the value is slightly
  reflected in the content) to 6 (i.e., the value is strongly reflected).

- **Labeler Service:** Based on `bsky-watch/labeler`, it emits cryptographically
  signed labels for posts that pass the value-based criteria.

- **Feed Generator:** Based on `go-bsky-feed-generator`, it serves the final
  timeline to users by querying the indexed database of evaluated posts.

- **Database (SQLite):** Stores post URIs, metadata, and calculated Schwartz
  scores.

---

## 3. The "Markdown Constitution"

Unlike traditional algorithms, the "logic" of this agent is stored in a
human-readable `.md` file. This file acts as the moral compass for the AI,
defining:

- Definitions of each Schwartz value.
- Scoring criteria for political discourse.
- Examples of content that should be promoted or de-prioritized.

---

## 4. Technical Stack

| Component         | Technology                                           |
| :---------------- | :--------------------------------------------------- |
| **Language**      | Golang (utilizing `indigo` and `jetstream` packages) |
| **Data Stream**   | Bluesky Jetstream (JSON via Websocket)               |
| **Orchestration** | Docker & Docker Compose                              |
| **AI/ML**         | Ollama (Local LLM) or OpenAI API                     |
| **Database**      | SQLite                                               |
| **Protocol**      | AT Protocol (Authenticated Transfer)                 |

---

## 5. Implementation Roadmap

1. **Phase 1:** ✅ DONE - Setup Golang Project, connect to Jetstream, filter posts by political keywords, print incoming data.
2. **Phase 2:** Create the Markdown Constitution with Schwartz values definitions.
3. **Phase 3:** Implement the LLM agent to score posts against the constitution.
4. **Phase 4:** Add SQLite database for storing scored posts.
5. **Phase 5:** Add Docker Compose for orchestration.
6. **Phase 6:** Implement Labeler and Feed Generator services.

---
