package handlers

import (
	"net/http"

	ras "app/tools"

	"github.com/USACE/filestore"
	"github.com/go-errors/errors" // warning: replaces standard errors
	"github.com/labstack/echo/v4"
)

// IsGeospatial godoc
// @Summary Check if the RAS model has geospatial information
// @Description Check if the RAS model has geospatial information
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/ras/CHURCH HOUSE GULLY/CHURCH HOUSE GULLY.prj"
// @Success 200 {object} bool
// @Router /isgeospatial [get]
func IsGeospatial(fs *filestore.FileStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")
		if definitionFile == "" {
			return c.JSON(http.StatusBadRequest, "Missing query parameter: `definition_file`")
		}

		rm, err := ras.NewRasModel(definitionFile, *fs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, SimpleResponse{http.StatusInternalServerError, err.Error(), err.(*errors.Error).ErrorStack()})
		}
		isIt := rm.IsGeospatial()

		return c.JSON(http.StatusOK, isIt)
	}
}
