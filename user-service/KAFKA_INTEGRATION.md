# User Service - Kafka Integration

This document explains how Kafka is integrated into the User Service for team activity event streaming.

## Overview

The User Service implements Kafka event publishing for team management operations according to the requirements specified in `kafka_redis.md`. It publishes events to the `team.activity` topic whenever team-related changes occur.

## Event Types

The service publishes the following event types:

### 1. TEAM_CREATED
Published when a new team is created.
```json
{
  "eventType": "TEAM_CREATED",
  "teamId": "uuid",
  "performedBy": "userId",
  "teamName": "Team Name",
  "timestamp": "2023-08-25T10:30:00Z"
}
```

### 2. MEMBER_ADDED
Published when a user is added as a team member.
```json
{
  "eventType": "MEMBER_ADDED",
  "teamId": "uuid",
  "performedBy": "userId",
  "targetUserId": "userId",
  "timestamp": "2023-08-25T10:30:00Z"
}
```

### 3. MEMBER_REMOVED
Published when a team member is removed.
```json
{
  "eventType": "MEMBER_REMOVED",
  "teamId": "uuid",
  "performedBy": "userId", 
  "targetUserId": "userId",
  "timestamp": "2023-08-25T10:30:00Z"
}
```

### 4. MANAGER_ADDED
Published when a user is added as a team manager.
```json
{
  "eventType": "MANAGER_ADDED",
  "teamId": "uuid",
  "performedBy": "userId",
  "targetUserId": "userId", 
  "timestamp": "2023-08-25T10:30:00Z"
}
```

### 5. MANAGER_REMOVED
Published when a team manager is removed.
```json
{
  "eventType": "MANAGER_REMOVED",
  "teamId": "uuid",
  "performedBy": "userId",
  "targetUserId": "userId",
  "timestamp": "2023-08-25T10:30:00Z"
}
```

## Architecture

### Components

1. **TeamServiceWithEvents**: Decorator that wraps the base TeamService and publishes events
2. **TeamActivityEventHandler**: Processes incoming team activity events
3. **Producer**: Kafka producer for publishing events
4. **Consumer**: Kafka consumer for processing events

### Data Flow

```
Team Operation → TeamServiceWithEvents → Database + Kafka Event → TeamActivityEventHandler
```

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_TEAM_ACTIVITY=team.activity
KAFKA_GROUP_ID=user-service-consumer

# Optional Kafka Tuning
KAFKA_MIN_BYTES=10240
KAFKA_MAX_BYTES=10485760
KAFKA_MAX_WAIT=250ms
KAFKA_BATCH_BYTES=1048576
KAFKA_BATCH_TIMEOUT=50ms
```

### Topic Configuration

The service expects the following Kafka topic to exist:
- **Topic**: `team.activity`
- **Partitions**: Recommended 3+ for scalability
- **Replication Factor**: Recommended 3 for production

## Usage

### Starting the Service

```bash
# Start with Kafka enabled
go run cmd/api/main.go
```

The service will:
1. Connect to the database
2. Initialize Kafka producer and consumer
3. Start HTTP server on configured port
4. Start Kafka consumer in background

### Event Publishing

Events are automatically published when team operations are performed through the API:

```bash
# Create a team (publishes TEAM_CREATED + MANAGER_ADDED events)
POST /api/teams
{
  "teamName": "Development Team",
  "managers": [{"userId": "manager-id"}],
  "members": [{"userId": "member-id"}]
}

# Add a member (publishes MEMBER_ADDED event)
POST /api/teams/{teamId}/members
{
  "userId": "new-member-id"
}
```

### Event Consumption

The service includes a sample event handler that logs all team activity events. You can extend `TeamActivityEventHandler` to:

- Store audit logs in database
- Update cache layers (Redis)
- Send notifications
- Trigger analytics
- Update search indexes

## Error Handling

### Event Publishing Failures

- Events are published after database operations succeed
- If event publishing fails, the operation continues (logged but not failed)
- This ensures database consistency isn't compromised by Kafka issues

### Event Consumption Failures

- Failed events are logged and not committed
- Consumer will retry the message
- Consider implementing dead letter queues for persistent failures

## Monitoring

### Metrics to Monitor

1. **Event Publishing**:
   - Event publish success/failure rates
   - Event publish latency
   - Producer lag

2. **Event Consumption**:
   - Consumer lag
   - Processing time per event
   - Error rates by event type

### Logs

The service logs all Kafka operations:
- Event publishing attempts
- Event processing results  
- Connection/configuration issues

## Development

### Testing Events

Use the included test helper:

```go
// In tests
components := app.Wire(kafkaConfig)
err := components.PublishTestTeam(ctx, "test-team-id")
```

### Adding New Event Types

1. Add event type constant to `kafka/types.go`
2. Add event handling to `TeamActivityEventHandler`
3. Publish events from appropriate service methods

## Production Considerations

### Performance
- Use appropriate batch sizes for your throughput requirements
- Monitor consumer lag and scale consumers if needed
- Consider async event publishing for high-throughput scenarios

### Reliability
- Use `RequiredAcks: kafka.RequireAll` for durability
- Implement circuit breakers for Kafka failures
- Use retry mechanisms with exponential backoff

### Security
- Configure SASL/SSL for production Kafka clusters
- Use separate service accounts for producers/consumers
- Implement proper access controls on topics

### Scaling
- Use team ID as message key for ordered processing per team
- Scale consumers horizontally by increasing consumer group instances
- Consider partitioning strategy based on team distribution

## Troubleshooting

### Common Issues

1. **Consumer not receiving messages**:
   - Check topic exists and has correct name
   - Verify consumer group ID is unique
   - Check Kafka broker connectivity

2. **High consumer lag**:
   - Increase number of consumer instances
   - Optimize event handler performance
   - Check for slow database operations

3. **Events not being published**:
   - Check producer configuration
   - Verify topic auto-creation is enabled or topic exists
   - Check network connectivity to Kafka brokers

### Debug Commands

```bash
# Check topic info
kafka-topics --bootstrap-server localhost:9092 --describe --topic team.activity

# Monitor consumer group
kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group user-service-consumer

# Consume messages manually
kafka-console-consumer --bootstrap-server localhost:9092 --topic team.activity --from-beginning
```
