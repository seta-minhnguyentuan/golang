package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"user-service/internal/app"
	"user-service/internal/config"
	"user-service/internal/models"

	"github.com/google/uuid"
)

// Example showing how to use the team service with Kafka events
func main() {
	// Load Kafka configuration
	kafkaCfg := config.LoadKafkaConfig()

	// Wire up components with Kafka support
	components := app.Wire(kafkaCfg)

	ctx := context.Background()

	// Example: Create a new team
	fmt.Println("🚀 Creating a new team...")

	creatorID := uuid.New()
	managerID := uuid.New()
	memberID := uuid.New()

	createTeamReq := models.CreateTeamRequest{
		TeamName: "Development Team",
		Managers: []models.TeamMemberRequest{
			{UserID: managerID.String()},
		},
		Members: []models.TeamMemberRequest{
			{UserID: memberID.String()},
		},
	}

	// This will publish the following events:
	// 1. TEAM_CREATED
	// 2. MANAGER_ADDED (for creator)
	// 3. MANAGER_ADDED (for additional manager)
	// 4. MEMBER_ADDED (for member)
	teamResponse, err := components.Teams.CreateTeam(ctx, createTeamReq, creatorID)
	if err != nil {
		log.Printf("❌ Failed to create team: %v", err)
		return
	}

	fmt.Printf("✅ Team created: %s (ID: %s)\n", teamResponse.TeamName, teamResponse.ID)

	// Give some time for events to be processed
	time.Sleep(100 * time.Millisecond)

	// Example: Add a new member
	fmt.Println("🚀 Adding a new member...")

	newMemberID := uuid.New()

	// This will publish a MEMBER_ADDED event
	err = components.Teams.AddMember(ctx, teamResponse.ID, newMemberID, creatorID)
	if err != nil {
		log.Printf("❌ Failed to add member: %v", err)
		return
	}

	fmt.Printf("✅ Member added: %s\n", newMemberID)

	// Give some time for events to be processed
	time.Sleep(100 * time.Millisecond)

	// Example: Remove a member
	fmt.Println("🚀 Removing a member...")

	// This will publish a MEMBER_REMOVED event
	err = components.Teams.RemoveMember(ctx, teamResponse.ID, newMemberID, creatorID)
	if err != nil {
		log.Printf("❌ Failed to remove member: %v", err)
		return
	}

	fmt.Printf("✅ Member removed: %s\n", newMemberID)

	// Clean up
	if err := components.Producer.Close(); err != nil {
		log.Printf("❌ Error closing producer: %v", err)
	}
	if err := components.Consumer.Close(); err != nil {
		log.Printf("❌ Error closing consumer: %v", err)
	}

	fmt.Println("✅ Example completed successfully!")

	fmt.Println("\n📋 Events that were published:")
	fmt.Println("1. TEAM_CREATED - when the team was created")
	fmt.Println("2. MANAGER_ADDED - for each manager added to the team")
	fmt.Println("3. MEMBER_ADDED - for each member added to the team")
	fmt.Println("4. MEMBER_ADDED - when new member was added")
	fmt.Println("5. MEMBER_REMOVED - when member was removed")

	fmt.Println("\n🔍 Check your Kafka topic 'team.activity' to see these events!")
}
