package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func validateParameters(scope string, permissions []string, c *gin.Context) error {

    supportedParams := []string{"key", "scope", "permissions"}

    for p := range c.Request.URL.Query() {
        found := false
        for _, param := range supportedParams {
            if p == param {
                found = true
                break
            }
        }
        if !found {
            return fmt.Errorf("invalid parameter: %s", p)
        }
    }

    if len(permissions) == 1 && permissions[0] == "*" {
        return nil
    }

    var validCommands = []string{"*"}

    switch scope {
    case "bastille":
        for _, cmd := range bastilleSpec.Commands {
            validCommands = append(validCommands, cmd.Command)
        }
    case "rocinante":
        for _, cmd := range rocinanteSpec.Commands {
            validCommands = append(validCommands, cmd.Command)
        }
    case "admin":
        validCommands = []string{"add", "delete", "edit"}
    default:
            return fmt.Errorf("invalid scope")
    }

    for _, p := range permissions {
        found := false
        for _, cmd := range validCommands {
            if p == cmd {
                found = true
                break
            }
        }
        if !found {
            return fmt.Errorf("invalid permission: %s", p)
        }
    }
    return nil
}

// Admin add POST
// @Description Add an API key.
// @Param Authorization header string true "Authentication token (e.g., Bearer <token>)"
// @Param key query string false "key"
// @Param scope query string false "scope"
// @Param permissions query string false "permissions"
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Router /api/v1/admin/add [post]
func AddKeyHandler(c *gin.Context) {

	key := c.Query("key")
	scope := c.Query("scope")
	permissionsQuery := c.Query("permissions")
	var permissions []string

	if key == "" {
		c.JSON(400, gin.H{"error": "Missing key parameter"})
		logRequest("error", "missing key parameter", c, nil, nil)
		return
	}
	if scope == "" {
		c.JSON(400, gin.H{"error": "Missing scope parameter"})
		logRequest("error", "missing scope parameter", c, nil, nil)
		return
	}
	if permissionsQuery == "" {
		c.JSON(400, gin.H{"error": "Missing permissions parameter"})
		logRequest("error", "missing permissions parameter", c, nil, nil)
		return
	} else {
		permissions = append(permissions, strings.Fields(permissionsQuery)...)
	}

	if _, exists := cfg.APIKeys[key]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Key already exists"})
		logRequest("error", "key already exists", c, nil, nil)
		return
	}

	if err := validateParameters(scope, permissions, c); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		logRequest("error", err.Error(), c, permissions, err.Error())
		return
	}

	newKey := APIKeyStruct{
		Permissions: PermissionsStruct{
			Bastille:  []string{},
			Rocinante: []string{},
			Admin:     []string{},
		},
	}

	switch scope {
	case "bastille":
		newKey.Permissions.Bastille = permissions
	case "rocinante":
		newKey.Permissions.Rocinante = permissions
	case "admin":
		newKey.Permissions.Admin = permissions
	default:
		c.JSON(400, gin.H{"error": "Invalid scope"})
		logRequest("error", "invalid scope", c, nil, nil)
		return
	}

	cfg.APIKeys[key] = newKey

	if err := saveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save key"})
		logRequest("error", "failed to save key", c, nil, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Key created"})
	logRequest("info", "key created", c, nil, nil)
}

// Admin edit POST
// @Description Edit an API key.
// @Param Authorization header string true "Authentication token (e.g., Bearer <token>)"
// @Param key query string false "key"
// @Param scope query string false "scope"
// @Param permissions query string false "permissions"
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Router /api/v1/admin/edit [post]
func EditKeyHandler(c *gin.Context) {

	key := c.Query("key")
	scope := c.Query("scope")
	permissionsQuery := c.Query("permissions")
	var permissions []string

	if key == "" {
		c.JSON(400, gin.H{"error": "Missing key parameter"})
		logRequest("error", "missing key parameter", c, nil, nil)
		return
	}
	if scope == "" {
		c.JSON(400, gin.H{"error": "Missing scope parameter"})
		logRequest("error", "missing scope parameter", c, nil, nil)
		return
	}
	if permissionsQuery == "" {
		c.JSON(400, gin.H{"error": "Missing permissions parameter"})
		logRequest("error", "missing permissions parameter", c, nil, nil)
		return
	} else {
		permissions = append(permissions, strings.Fields(permissionsQuery)...)
	}

	keyData, exists := cfg.APIKeys[key]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		logRequest("error", "key not found", c, nil, nil)
		return
	}

	if err := validateParameters(scope, permissions, c); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		logRequest("error", err.Error(), c, permissions, err.Error())
		return
	}

	switch scope {
	case "bastille":
		keyData.Permissions.Bastille = permissions
	case "rocinante":
		keyData.Permissions.Rocinante = permissions
	case "admin":
		keyData.Permissions.Admin = permissions
	default:
		c.JSON(400, gin.H{"error": "Invalid scope"})
		logRequest("error", "invalid scope", c, nil, nil)
		return
	}

	cfg.APIKeys[key] = keyData

	if err := saveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		logRequest("error", "failed to save config", c, nil, nil)
		return
	}

	logRequest("info", "Key updated", c, nil, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Key updated"})
}

// Admin add POST
// @Description Delete an API key.
// @Param Authorization header string true "Authentication token (e.g., Bearer <token>)"
// @Param key query string false "key"
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Router /api/v1/admin/delete [post]
func DeleteKeyHandler(c *gin.Context) {

	key := c.Query("key")

	if key == "" {
		c.JSON(400, gin.H{"error": "Missing key parameter"})
		logRequest("error", "missing key parameter", c, nil, nil)
		return
	} else if len(c.Request.URL.Query()) != 1 {
	        c.JSON(400, gin.H{"error": "Invalid parameters"})
	        logRequest("error", "invalid parameters", c, nil, nil)
	        return
	}

	if _, exists := cfg.APIKeys[key]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		logRequest("error", "key not found", c, nil, nil)
		return
	}

	delete(cfg.APIKeys, key)

	if err := saveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		logRequest("error", "failed to save config", c, nil, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key deleted"})
	logRequest("info", "key deleted", c, nil, nil)
}