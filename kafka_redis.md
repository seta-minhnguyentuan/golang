# 📌 Focused Kafka Integration: Team & Asset Changes

These exercises help integrate **Kafka** into the Team Management and Asset Sharing parts of your microservices system.

---

## 👥 Team Event Notification via Kafka

**Goal**: Emit events whenever team data changes—creation, member updates, or manager changes.

- **Topic**: `team.activity`
- **Events**:
  - `TEAM_CREATED`
  - `MEMBER_ADDED`
  - `MEMBER_REMOVED`
  - `MANAGER_ADDED`
  - `MANAGER_REMOVED`
- **Payload Example**:
  ```json
  {
    "eventType": "MEMBER_ADDED",
    "teamId": "uuid",
    "performedBy": "userId",
    "targetUserId": "userId",
    "timestamp": "ISO-8601"
  }
  ```
- **Consumer Suggestions**:
  - Log to database or ElasticSearch.

---

## 🗂 Asset Change Stream via Kafka

**Goal**: Emit events when assets (folders/notes) are created, updated, deleted, or shared.

- **Topic**: `asset.changes`
- **Events**:
  - `FOLDER_CREATED`, `FOLDER_UPDATED`, `FOLDER_DELETED`
  - `NOTE_CREATED`, `NOTE_UPDATED`, `NOTE_DELETED`
  - `FOLDER_SHARED`, `FOLDER_UNSHARED`
  - `NOTE_SHARED`, `NOTE_UNSHARED`
- **Payload Example**:
  ```json
  {
    "eventType": "NOTE_UPDATED",
    "assetType": "note",
    "assetId": "uuid",
    "ownerId": "userId",
    "actionBy": "userId",
    "timestamp": "ISO-8601"
  }
  ```
- **Consumer Suggestions**:
  - Maintain audit logs.
  - Trigger cache invalidation or real-time UI updates.
  - Index assets into search system.

---

## 🔧 Technical Notes

- Use `confluent-kafka-go` or `segmentio/kafka-go` libraries for Go.
- Ensure all events include `eventType`, `timestamp`, and identifiers.
- Consider retry & error handling in Kafka consumers.

## 🔴 Redis-Focused Exercises

### 1. Real-Time Team Member Cache
- **Key format**: `team:{teamId}:members` → list of `userId`.
- Update Redis cache when `MEMBER_ADDED` / `MEMBER_REMOVED` events occur.
- API reads from Redis, fallback to DB if cache miss.

### 2. Asset Metadata Cache
- **Key format**:
  - `folder:{folderId}` → folder metadata JSON.
  - `note:{noteId}` → note metadata JSON.
- Update or invalidate cache on `asset.changes` events.
- Use write-through strategy.

### 3. Real-Time Access Control Lookup
- **Key format**: `asset:{assetId}:acl` → `{userId: accessType}`.
- Update on `FOLDER_SHARED` / `NOTE_SHARED` events.
- Validate permissions from Redis before DB query.

---