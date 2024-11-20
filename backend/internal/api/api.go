package api

import (
	"log"
	"net/http"
	"time"

	"efaturas-xtreme/internal/api/requests"
	"efaturas-xtreme/internal/api/responses"
	"efaturas-xtreme/internal/auth"
	"efaturas-xtreme/internal/service"
	"efaturas-xtreme/internal/service/domain"
	"efaturas-xtreme/internal/session"
	"efaturas-xtreme/pkg/errors"
	"efaturas-xtreme/pkg/sse"

	"github.com/gin-gonic/gin"
)

type API struct {
	sse     sse.Server
	auth    *auth.Service
	service service.Service
}

func (api *API) Init(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	auth := r.Group("")
	{
		auth.POST("login", api.login)
	}

	invoices := r.Group("invoices", authMiddleware)
	{
		invoices.GET("", api.getInvoices)
		invoices.POST("", api.fetchNewInvoices)
		invoices.PUT("", api.processInvoices)
		invoices.PUT("categories", api.updateCategoriesInvoices)
		invoices.GET(":invoiceID", api.getInvoiceSSE)
	}

	categories := r.Group("categories", authMiddleware)
	{
		categories.GET("", api.getCategories)
	}
}

func (api *API) login(ctx *gin.Context) {
	var request requests.Create
	if err := ctx.ShouldBind(&request); err != nil {
		_ = ctx.Error(errors.NewInput(err))
		return
	}

	session, err := api.auth.Login(ctx, request.Username, request.Password)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": session.Value})
}

func (api *API) getInvoices(ctx *gin.Context) {
	userID, err := session.GetUserID(ctx)
	if err != nil {
		_ = ctx.Error(errors.New(err))
		return
	}

	invoices, err := api.service.GetInvoices(ctx, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, responses.NewListInvoices(invoices))
}

func (api *API) getInvoiceSSE(ctx *gin.Context) {
	invoiceID := ctx.Param("invoiceID")
	api.sse.Subscribe(ctx.Writer, ctx.Request, invoiceID)
}

func (api *API) fetchNewInvoices(ctx *gin.Context) {
	userID, uname, pword, err := session.GetUser(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	invoices, err := api.service.CreateOrUpdate(ctx, userID, uname, pword)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, responses.NewListInvoices(invoices))
}

func (api *API) processInvoices(ctx *gin.Context) {
	userID, uname, pword, err := session.GetUser(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if err = api.service.ScanInvoices(ctx, userID, uname, pword); err != nil {
		_ = ctx.Error(err)
		return
	}

	invoices, err := api.service.GetInvoices(ctx, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, responses.NewListInvoices(invoices))
}

func (api *API) updateCategoriesInvoices(ctx *gin.Context) {
	var request requests.UpdateInvoices
	if err := ctx.ShouldBind(&request); err != nil {
		_ = ctx.Error(errors.NewInput(err))
		return
	}

	userID, uname, pword, err := session.GetUser(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	invoices, err := api.service.UpdateInvoiceCategories(ctx, userID, uname, pword, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, responses.NewListInvoices(invoices))
}

func (api *API) getCategories(ctx *gin.Context) {
	ctx.Header("Cache-Control", "public, max-age=86400")
	ctx.Header("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))
	ctx.Header("ETag", `W/"12345"`)

	ctx.JSON(http.StatusOK, responses.NewGetCategories())
}

func (api *API) onUpdateInvoice(inv *domain.Invoice, done bool) {
	if err := api.sse.Publish(inv.GetID(), responses.NewInvoiceUpdate(inv, done)); err != nil {
		log.Println("failed to publish:", inv.ID, err)
	}
}

func New(service service.Service, auth *auth.Service, sse sse.Server) *API {
	api := &API{sse: sse, auth: auth, service: service}
	service.SetOnUpdateInvoice(api.onUpdateInvoice)

	return api
}
