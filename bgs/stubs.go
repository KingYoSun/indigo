package bgs

import (
	"encoding/json"
	"io"
	"strconv"

	comatprototypes "github.com/KingYoSun/indigo/api/atproto"
	"github.com/labstack/echo/v4"
	"github.com/meilisearch/meilisearch-go"
	"go.opentelemetry.io/otel"
)

func (s *BGS) RegisterHandlersAppBsky(e *echo.Echo) error {
	return nil
}

func (s *BGS) RegisterHandlersComAtproto(e *echo.Echo) error {
	e.GET("/xrpc/com.atproto.sync.getBlob", s.HandleComAtprotoSyncGetBlob)
	e.GET("/xrpc/com.atproto.sync.getBlocks", s.HandleComAtprotoSyncGetBlocks)
	e.GET("/xrpc/com.atproto.sync.getCheckout", s.HandleComAtprotoSyncGetCheckout)
	e.GET("/xrpc/com.atproto.sync.getCommitPath", s.HandleComAtprotoSyncGetCommitPath)
	e.GET("/xrpc/com.atproto.sync.getHead", s.HandleComAtprotoSyncGetHead)
	e.GET("/xrpc/com.atproto.sync.getRecord", s.HandleComAtprotoSyncGetRecord)
	e.GET("/xrpc/com.atproto.sync.getRepo", s.HandleComAtprotoSyncGetRepo)
	e.GET("/xrpc/com.atproto.sync.listBlobs", s.HandleComAtprotoSyncListBlobs)
	e.GET("/xrpc/com.atproto.sync.listRepos", s.HandleComAtprotoSyncListRepos)
	e.GET("/xrpc/com.atproto.sync.notifyOfUpdate", s.HandleComAtprotoSyncNotifyOfUpdate)
	e.GET("/xrpc/com.atproto.sync.requestCrawl", s.HandleComAtprotoSyncRequestCrawl)
	e.GET("/debug/getRepo", s.HandleDebugGetRecord)
	e.GET("/meili/requestCopyRecord", s.HandleMeiliRequestCopyRecord)
	e.POST("/meili/updateIndexSettings/:index", s.HandleMeiliUpdateIndexSettings)
	e.GET("/meili/search", s.HandleMeiliSearch)
	return nil
}

func (s *BGS) HandleComAtprotoSyncGetBlob(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncGetBlob")
	defer span.End()
	cid := c.QueryParam("cid")
	did := c.QueryParam("did")
	var out io.Reader
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncGetBlob(ctx context.Context,cid string,did string) (io.Reader, error)
	out, handleErr = s.handleComAtprotoSyncGetBlob(ctx, cid, did)
	if handleErr != nil {
		return handleErr
	}
	return c.Stream(200, "application/octet-stream", out)
}

func (s *BGS) HandleComAtprotoSyncGetBlocks(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncGetBlocks")
	defer span.End()

	cids := c.QueryParams()["cids"]
	did := c.QueryParam("did")
	var out io.Reader
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncGetBlocks(ctx context.Context,cids []string,did string) (io.Reader, error)
	out, handleErr = s.handleComAtprotoSyncGetBlocks(ctx, cids, did)
	if handleErr != nil {
		return handleErr
	}
	return c.Stream(200, "application/vnd.ipld.car", out)
}

func (s *BGS) HandleComAtprotoSyncGetCheckout(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncGetCheckout")
	defer span.End()
	commit := c.QueryParam("commit")
	did := c.QueryParam("did")
	var out io.Reader
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncGetCheckout(ctx context.Context,commit string,did string) (io.Reader, error)
	out, handleErr = s.handleComAtprotoSyncGetCheckout(ctx, commit, did)
	if handleErr != nil {
		return handleErr
	}
	return c.Stream(200, "application/vnd.ipld.car", out)
}

func (s *BGS) HandleComAtprotoSyncGetCommitPath(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncGetCommitPath")
	defer span.End()
	did := c.QueryParam("did")
	earliest := c.QueryParam("earliest")
	latest := c.QueryParam("latest")
	var out *comatprototypes.SyncGetCommitPath_Output
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncGetCommitPath(ctx context.Context,did string,earliest string,latest string) (*comatprototypes.SyncGetCommitPath_Output, error)
	out, handleErr = s.handleComAtprotoSyncGetCommitPath(ctx, did, earliest, latest)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *BGS) HandleComAtprotoSyncGetHead(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncGetHead")
	defer span.End()
	did := c.QueryParam("did")
	var out *comatprototypes.SyncGetHead_Output
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncGetHead(ctx context.Context,did string) (*comatprototypes.SyncGetHead_Output, error)
	out, handleErr = s.handleComAtprotoSyncGetHead(ctx, did)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *BGS) HandleComAtprotoSyncGetRecord(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncGetRecord")
	defer span.End()
	collection := c.QueryParam("collection")
	commit := c.QueryParam("commit")
	did := c.QueryParam("did")
	rkey := c.QueryParam("rkey")
	var out io.Reader
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncGetRecord(ctx context.Context,collection string,commit string,did string,rkey string) (io.Reader, error)
	out, handleErr = s.handleComAtprotoSyncGetRecord(ctx, collection, commit, did, rkey)
	if handleErr != nil {
		return handleErr
	}
	return c.Stream(200, "application/vnd.ipld.car", out)
}

func (s *BGS) HandleComAtprotoSyncGetRepo(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncGetRepo")
	defer span.End()
	did := c.QueryParam("did")
	earliest := c.QueryParam("earliest")
	latest := c.QueryParam("latest")
	var out io.Reader
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncGetRepo(ctx context.Context,did string,earliest string,latest string) (io.Reader, error)
	out, handleErr = s.handleComAtprotoSyncGetRepo(ctx, did, earliest, latest)
	if handleErr != nil {
		return c.JSON(500, handleErr.Error())
	}
	return c.Stream(200, "application/vnd.ipld.car", out)
}

func (s *BGS) HandleComAtprotoSyncListBlobs(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncListBlobs")
	defer span.End()
	did := c.QueryParam("did")
	earliest := c.QueryParam("earliest")
	latest := c.QueryParam("latest")
	var out *comatprototypes.SyncListBlobs_Output
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncListBlobs(ctx context.Context,did string,earliest string,latest string) (*comatprototypes.SyncListBlobs_Output, error)
	out, handleErr = s.handleComAtprotoSyncListBlobs(ctx, did, earliest, latest)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *BGS) HandleComAtprotoSyncListRepos(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncListRepos")
	defer span.End()
	cursor := c.QueryParam("cursor")

	var limit int
	if p := c.QueryParam("limit"); p != "" {
		var err error
		limit, err = strconv.Atoi(p)
		if err != nil {
			return err
		}
	} else {
		limit = 500
	}
	var out *comatprototypes.SyncListRepos_Output
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncListRepos(ctx context.Context,cursor string,limit int) (*comatprototypes.SyncListRepos_Output, error)
	out, handleErr = s.handleComAtprotoSyncListRepos(ctx, cursor, limit)
	if handleErr != nil {
		return handleErr
	}
	return c.JSON(200, out)
}

func (s *BGS) HandleComAtprotoSyncNotifyOfUpdate(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncNotifyOfUpdate")
	defer span.End()
	hostname := c.QueryParam("hostname")
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncNotifyOfUpdate(ctx context.Context,hostname string) error
	handleErr = s.handleComAtprotoSyncNotifyOfUpdate(ctx, hostname)
	if handleErr != nil {
		return handleErr
	}
	return nil
}

func (s *BGS) HandleComAtprotoSyncRequestCrawl(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleComAtprotoSyncRequestCrawl")
	defer span.End()
	hostname := c.QueryParam("hostname")
	var handleErr error
	// func (s *BGS) handleComAtprotoSyncRequestCrawl(ctx context.Context,hostname string) error
	handleErr = s.handleComAtprotoSyncRequestCrawl(ctx, hostname)
	if handleErr != nil {
		return handleErr
	}
	return nil
}

func (s *BGS) HandleDebugGetRecord(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleDebugGetRepo")
	defer span.End()
	did := c.QueryParam("did")
	cid := c.QueryParam("cid")
	rkey := c.QueryParam("rkey")

  out, err := s.handleDebugGetRepoJson(ctx, did, cid, rkey)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, string(out))
}

func (s *BGS) HandleMeiliRequestCopyRecord(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleMeiliRequestCopyRecord")
	defer span.End()
	hostname := c.QueryParam("hostname")
	// func (s *BGS) handleComAtprotoSyncRequestCrawl(ctx context.Context,hostname string) error
	handleErr := s.handleMeiliRequestCopyRecord(ctx, hostname)
	if handleErr != nil {
		return c.JSON(500, handleErr.Error())
	}

	return nil
}

func (s *BGS) HandleMeiliUpdateIndexSettings(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleMeiliUpdateIndexSettings")
	defer span.End()

	index := c.Param("index")
	var settings meilisearch.Settings
	if err := c.Bind(&settings); err != nil {
		return c.JSON(500, err.Error())
	}

	resp, err := s.meilislur.UpdateIndexSetting(ctx, index, settings)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, resp.Status)
}

func (s *BGS) HandleMeiliSearch(c echo.Context) error {
	ctx, span := otel.Tracer("server").Start(c.Request().Context(), "HandleMeiliSearch")
	defer span.End()

	keyword := c.QueryParam("q")
	hostname := c.QueryParam("h")
	sort := c.QueryParam("s")

	if sort == "" {
		sort = ""
	}

	var posts []interface{}
	var err error
	if posts, err = s.meilislur.Search(ctx, keyword, hostname, sort); err != nil {
		return c.JSON(500, err.Error())
	}

	var out []byte

	if out, err = json.Marshal(posts); err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, string(out))
}