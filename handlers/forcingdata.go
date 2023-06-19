package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ar-siddiqui/mcat-ras/config"
	"github.com/ar-siddiqui/mcat-ras/tools"

	"github.com/USACE/filestore"
	"github.com/go-errors/errors" // warning: replaces standard errors
	"github.com/labstack/echo/v4"
)

// ForcingData godoc
// @Summary Extract forcing data from flow files
// @Description forcing data from a RAS model given an s3 key
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/ras/CHURCH HOUSE GULLY/CHURCH HOUSE GULLY.prj"
// @Success 200 {object} interface{}
// @Failure 500 {object} SimpleResponse
// @Router /forcingdata [get]
func ForcingData(ac *config.APIConfig) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")
		if definitionFile == "" {
			return c.JSON(http.StatusBadRequest, "Missing query parameter: `definition_file`")
		}

		if !isAModel(ac.FileStore, definitionFile) {
			return c.JSON(http.StatusBadRequest, definitionFile+" is not a valid RAS prj file.")
		}

		data, err := forcingData(definitionFile, ac.FileStore)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, SimpleResponse{http.StatusInternalServerError, fmt.Sprintf("Go error encountered: %v", err.Error()), err.(*errors.Error).ErrorStack()})
		}

		return c.JSON(http.StatusOK, data)
	}
}

func forcingData(definitionFile string, fs *filestore.FileStore) (tools.ForcingData, error) {
	fd := tools.ForcingData{
		Steady:   make(map[string]tools.SteadyData),
		Unsteady: make(map[string]tools.UnsteadyData),
	}

	mfiles, err := modFiles(definitionFile, *fs)
	if err != nil {
		return fd, errors.Wrap(err, 0)
	}

	for _, fp := range mfiles {

		ext := filepath.Ext(fp)

		switch {

		case tools.RasRE.AllFlow.MatchString(ext):

			if err := tools.GetForcingData(&fd, *fs, fp); err != nil {
				return fd, errors.Wrap(err, 0)
			}

		}
	}

	return fd, nil
}
