package routes

import (
	"strings"

	"dav_docker/utils"

	"github.com/gofiber/fiber/v2"
)

type ConfigItem struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Root     string `json:"root"`
	Running  int    `json:"running"`
	Name     string `json:"name"`
}

type ConfigEditRequest struct {
	Path     *string `json:"path"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Root     *string `json:"root"`
	Running  *int    `json:"running"`
	Name     *string `json:"name"`
}

func JSONResponse(c *fiber.Ctx, ok bool, data any) error {
	return c.JSON(fiber.Map{
		"ok":   ok,
		"data": data,
	})
}

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("token")
	if token == "" {
		return JSONResponse(c, false, "missing token")
	}

	username, err, isExpired := utils.ValidateAccessToken(token)
	if err != nil {
		if isExpired {
			return JSONResponse(c, false, "expired")
		}
		return JSONResponse(c, false, err.Error())
	}

	c.Locals("username", username)
	return c.Next()
}

func HandleConfigList(c *fiber.Ctx) error {
	rows, err := utils.DB.Query("SELECT id, path, username, root, running, name FROM config")
	if err != nil {
		return JSONResponse(c, false, err.Error())
	}
	defer rows.Close()

	var configs []ConfigItem
	for rows.Next() {
		var item ConfigItem
		if err := rows.Scan(&item.ID, &item.Path, &item.Username, &item.Root, &item.Running, &item.Name); err != nil {
			return JSONResponse(c, false, err.Error())
		}
		configs = append(configs, item)
	}

	if configs == nil {
		configs = []ConfigItem{}
	}

	return JSONResponse(c, true, configs)
}

func HandleConfigAdd(c *fiber.Ctx) error {
	var item ConfigItem
	if err := c.BodyParser(&item); err != nil {
		return JSONResponse(c, false, "invalid request body")
	}

	if item.Path == "" || item.Username == "" || item.Password == "" || item.Root == "" || item.Name == "" {
		return JSONResponse(c, false, "path, username, password, root and name are required")
	}

	var pathExists bool
	utils.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM config WHERE path = ?)", item.Path).Scan(&pathExists)
	if pathExists {
		return JSONResponse(c, false, "path already exists")
	}

	encryptedPwd, err := utils.Encrypt(item.Password, utils.GetAESKey())
	if err != nil {
		return JSONResponse(c, false, "encryption failed")
	}

	id := utils.GenerateID()
	_, err = utils.DB.Exec("INSERT INTO config (id, path, username, password, root, running, name) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, item.Path, item.Username, encryptedPwd, item.Root, item.Running, item.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique") {
			return JSONResponse(c, false, "path already exists")
		}
		return JSONResponse(c, false, err.Error())
	}

	return JSONResponse(c, true, "")
}

func HandleConfigEdit(c *fiber.Ctx) error {
	id := c.Params("id")
	var req ConfigEditRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONResponse(c, false, "invalid request body")
	}

	if req.Path == nil && req.Username == nil && req.Password == nil && req.Root == nil && req.Running == nil && req.Name == nil {
		return JSONResponse(c, false, "at least one parameter is required")
	}

	var exists bool
	err := utils.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM config WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return JSONResponse(c, false, "config not found")
	}

	var queryParts []string
	var args []interface{}

	if req.Path != nil {
		var pathExists bool
		utils.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM config WHERE path = ? AND id != ?)", *req.Path, id).Scan(&pathExists)
		if pathExists {
			return JSONResponse(c, false, "path already exists")
		}
		queryParts = append(queryParts, "path = ?")
		args = append(args, *req.Path)
	}
	if req.Username != nil {
		queryParts = append(queryParts, "username = ?")
		args = append(args, *req.Username)
	}
	if req.Password != nil {
		encryptedPwd, err := utils.Encrypt(*req.Password, utils.GetAESKey())
		if err != nil {
			return JSONResponse(c, false, "encryption failed")
		}
		queryParts = append(queryParts, "password = ?")
		args = append(args, encryptedPwd)
	}
	if req.Root != nil {
		queryParts = append(queryParts, "root = ?")
		args = append(args, *req.Root)
	}
	if req.Running != nil {
		queryParts = append(queryParts, "running = ?")
		args = append(args, *req.Running)
	}
	if req.Name != nil {
		queryParts = append(queryParts, "name = ?")
		args = append(args, *req.Name)
	}

	args = append(args, id)
	query := "UPDATE config SET " + strings.Join(queryParts, ", ") + " WHERE id = ?"
	_, err = utils.DB.Exec(query, args...)
	if err != nil {
		return JSONResponse(c, false, err.Error())
	}

	return JSONResponse(c, true, "")
}

func HandleConfigDel(c *fiber.Ctx) error {
	id := c.Params("id")
	var exists bool
	err := utils.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM config WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return JSONResponse(c, false, "config not found")
	}

	_ = utils.StopWebDAVServer(id)

	_, err = utils.DB.Exec("DELETE FROM config WHERE id = ?", id)
	if err != nil {
		return JSONResponse(c, false, err.Error())
	}

	return JSONResponse(c, true, "")
}

func HandleConfigRun(c *fiber.Ctx) error {
	id := c.Params("id")
	var path, username, password, root string
	err := utils.DB.QueryRow("SELECT path, username, password, root FROM config WHERE id = ?", id).Scan(&path, &username, &password, &root)
	if err != nil {
		return JSONResponse(c, false, "config not found")
	}

	if err := utils.StartWebDAVServer(id, path, username, password, root); err != nil {
		return JSONResponse(c, false, err.Error())
	}

	_, err = utils.DB.Exec("UPDATE config SET running = 1 WHERE id = ?", id)
	if err != nil {
		return JSONResponse(c, false, err.Error())
	}

	return JSONResponse(c, true, "")
}

func HandleConfigStop(c *fiber.Ctx) error {
	id := c.Params("id")
	var exists bool
	err := utils.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM config WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return JSONResponse(c, false, "config not found")
	}

	if err := utils.StopWebDAVServer(id); err != nil {
		return JSONResponse(c, false, err.Error())
	}

	_, err = utils.DB.Exec("UPDATE config SET running = 0 WHERE id = ?", id)
	if err != nil {
		return JSONResponse(c, false, err.Error())
	}

	return JSONResponse(c, true, "")
}
