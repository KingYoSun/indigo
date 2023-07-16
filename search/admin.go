package search

import (
	"github.com/labstack/echo/v4"
	meilisearch "github.com/meilisearch/meilisearch-go"
	"go.opentelemetry.io/otel"
)

func (s *Server) handleUpdateMeiliIndexSetting(e echo.Context) error {
	_, span := otel.Tracer("server").Start(e.Request().Context(), "HandleMeiliUpdateIndexSettings")
	defer span.End()

	index := e.Param("index")
	var settings meilisearch.Settings

	if err := e.Bind(&settings); err != nil {
		return err
	}

	resp, err := s.meilicli.Index(index).UpdateSettings(&settings)
	if err != nil {
		return err
	}
	return e.JSON(200, resp.Status)
}
