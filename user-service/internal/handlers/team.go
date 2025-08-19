package handlers

import (
	"net/http"
	"user-service/internal/models"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	TeamService services.TeamService
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		TeamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req models.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT claims
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	creatorID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	teamResponse, err := h.TeamService.CreateTeam(c.Request.Context(), req, creatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teamResponse)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamIDParam := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	teamResponse, err := h.TeamService.GetTeamByID(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	c.JSON(http.StatusOK, teamResponse)
}

func (h *TeamHandler) GetAllTeams(c *gin.Context) {
	requestorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestorUUID, ok := requestorID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	teams, err := h.TeamService.GetAllTeams(c.Request.Context(), requestorUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

func (h *TeamHandler) AddMember(c *gin.Context) {
	teamIDParam := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req models.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get requestor ID from JWT claims
	requestorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestorUUID, ok := requestorID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.TeamService.AddMember(c.Request.Context(), teamID, userID, requestorUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamIDParam := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	memberIDParam := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	// Get requestor ID from JWT claims
	requestorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestorUUID, ok := requestorID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.TeamService.RemoveMember(c.Request.Context(), teamID, memberID, requestorUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

func (h *TeamHandler) AddManager(c *gin.Context) {
	teamIDParam := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req models.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get requestor ID from JWT claims
	requestorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestorUUID, ok := requestorID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.TeamService.AddManager(c.Request.Context(), teamID, userID, requestorUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager added successfully"})
}

func (h *TeamHandler) RemoveManager(c *gin.Context) {
	teamIDParam := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	managerIDParam := c.Param("managerId")
	managerID, err := uuid.Parse(managerIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID"})
		return
	}

	// Get requestor ID from JWT claims
	requestorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestorUUID, ok := requestorID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.TeamService.RemoveManager(c.Request.Context(), teamID, managerID, requestorUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed successfully"})
}
