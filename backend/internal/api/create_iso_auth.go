package api

import (
	"crypto/subtle"
	"fmt"
	"strings"

	"linux-iso-manager/internal/config"

	"github.com/gin-gonic/gin"
	ldap "github.com/go-ldap/ldap/v3"
)

func createISOAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Server.CreateISOAuthEnabled {
			c.Next()
			return
		}

		if c.Request.Method != "POST" || c.Request.URL.Path != "/api/isos" {
			c.Next()
			return
		}

		username, password, ok := c.Request.BasicAuth()
		if !ok || username == "" || password == "" {
			c.Header("WWW-Authenticate", "Basic realm=\"ISO Create\"")
			ErrorResponse(c, 401, "UNAUTHORIZED", "Authentication required")
			c.Abort()
			return
		}

		authenticated := false
		var err error

		if cfg.Server.LDAPAuthEnabled && cfg.Server.LDAPURL != "" && cfg.Server.LDAPUsersBaseDN != "" {
			authenticated, err = authenticateWithLDAP(cfg, username, password)
		} else if cfg.Server.BasicAuthUsername != "" && cfg.Server.BasicAuthPassword != "" {
			authenticated = subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Server.BasicAuthUsername)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(cfg.Server.BasicAuthPassword)) == 1
		} else {
			ErrorResponse(c, 500, "INTERNAL_ERROR", "Create ISO auth is enabled but no auth backend is configured")
			c.Abort()
			return
		}

		if err != nil || !authenticated {
			c.Header("WWW-Authenticate", "Basic realm=\"ISO Create\"")
			ErrorResponse(c, 401, "UNAUTHORIZED", "Invalid credentials")
			c.Abort()
			return
		}

		c.Next()
	}
}

func authenticateWithLDAP(cfg *config.Config, username, password string) (bool, error) {
	conn, err := ldap.DialURL(cfg.Server.LDAPURL)
	if err != nil {
		return false, fmt.Errorf("ldap dial failed: %w", err)
	}
	defer conn.Close()

	if cfg.Server.LDAPBindDN != "" {
		if err := conn.Bind(cfg.Server.LDAPBindDN, cfg.Server.LDAPBindPassword); err != nil {
			return false, fmt.Errorf("ldap service bind failed: %w", err)
		}
	}

	escapedUser := ldap.EscapeFilter(username)
	filter := strings.ReplaceAll(cfg.Server.LDAPUserFilter, "{user}", escapedUser)
	searchReq := ldap.NewSearchRequest(
		cfg.Server.LDAPUsersBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		filter,
		[]string{"dn"},
		nil,
	)

	searchRes, err := conn.Search(searchReq)
	if err != nil {
		return false, fmt.Errorf("ldap search failed: %w", err)
	}
	if len(searchRes.Entries) != 1 {
		return false, nil
	}

	userDN := searchRes.Entries[0].DN
	if err := conn.Bind(userDN, password); err != nil {
		return false, nil
	}

	return true, nil
}
