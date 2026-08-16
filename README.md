### Architecture

#### Key Decisions

* **Monolith:** `REST API + gRPC`. The Telegram bot used for notifications was moved into a separate gRPC microservice.
* **Structure:** Domain-Driven Design and SOLID principles were applied, along with a clear separation of layers. An extended `controllers -> services -> repository` architecture is used. DI (Dependency Injection) was implemented.
* **Separation of Guardrails and Analytics:** Guardrails use `received_at` (when the server received the event), while analytics use `client_time` (when the client sent the event).
* **Autopilot ramp-up:** Implemented in a straightforward way. It gradually increases experiment traffic according to **predefined** rules (`1 -> 5 -> 25 -> 50 -> 100`). Autopilot interacts with guardrail monitoring and stops the ramp-up when a guardrail is triggered.
* **Initial seed data:** Initial metrics, admins, and event types.
* **Decide function and its stickiness:** Implemented according to the documentation at `support.leanplum.com`, except that in our implementation the range is from 0 to 99. A user consistently receives the same decision because it is based on the hash sum and a salt. The formula is `hashsum(userID + experimentID)`, where the experiment ID is used as the `salt`.

---

#### Notifications

* **Delivery channels:**
  The system supports two notification channels: **Telegram** (through a separate gRPC microservice) and **Slack** (sent directly from the main application via webhook). Adding a new notification channel is straightforward — it only requires implementing it in the `notifications` package.

* **Deduplication:**
  To prevent spam in chats, an atomic Redis `SetNX` operation with a TTL (key expiration time) is used. This guarantees that even when identical notification requests are processed concurrently, the same notification will be sent only once within the specified time window.

* **Recipient configuration:**
  Each experiment can have its own notification settings stored in the `notification_settings` table. Two recipient lists are supported: `chat_ids` for Telegram and `slack_webhooks` for Slack. Before sending a notification, the notification service retrieves these lists and sends the message to all specified recipients.

* **Asynchronous processing:**
  All notifications are sent in background goroutines using their own context (`context.Background()`), so notification delivery does not delay the processing of primary requests. Delivery errors are logged but do not affect the execution of the main operation.

---

#### Linters and Formatting

* **Formatting:** The standard `go fmt` tool is used to automatically format the code according to the consistent style adopted by the Go ecosystem.
* **Linting:** Qodana is used for deeper static analysis and detection of potential issues. It checks the code for errors, unused code, and security issues. During local development, it can also generate statistics about detected problems.

---

#### Running and Usage

---

#### Where Code Was Generated

**LLMs/neural networks were used ONLY for routine tasks. Below are a few cases where their use went beyond that:**

* **Validators** — There was not enough time to manually write and then verify every validator, so LLMs were used to assist with this part.
* **DSL** — I reused the DSL from a previous task, where I had used an LLM during its development. However, it was not written entirely by an LLM — I already have experience writing parsers. The main areas where I relied on an LLM while working on the DSL were the `IN` operator and combining the `AND` and `OR` logical operators. They worked correctly separately, but not when used together.

**I wrote all the remaining logic myself. For critical parts of the implementation, I studied the relevant documentation and watched technical videos (`api.slack.com`, `support.leanplum.com`, `youtube.com`).**

--- 

### Additionally
**You can also find the task of this project [here](https://github.com/Central-University-IT-prod/2025-2026-tasks/blob/main/Indiv/Backend/task-en.md)**
