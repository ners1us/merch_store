package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enums"
	"github.com/ners1us/merch_store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbTest, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbTest.AutoMigrate(&models.User{}, &models.Purchase{}, &models.CoinTransfer{}, &models.Merch{})
	require.NoError(t, err)

	db = dbTest
	return dbTest
}

func TestHandleSendCoinUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", nil)

	HandleSendCoin(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, enums.ErrUserNotAuthorized.Error(), response["error"])
}

func TestHandleSendCoinInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	sender := models.User{Username: "sender", Password: "password", Coins: 1000}
	require.NoError(t, db.Create(&sender).Error)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", sender)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBufferString("invalid json"))

	HandleSendCoin(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.ErrWrongReqFormat.Error(), response["error"])
}

func TestHandleSendCoinInvalidAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	sender := models.User{Username: "sender", Password: "password", Coins: 1000}
	require.NoError(t, db.Create(&sender).Error)

	requestBody := models.SendCoinRequest{
		ToUser: "receiver",
		Amount: 0,
	}
	bodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", sender)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(bodyBytes))
	ctx.Request.Header.Set("Content-Type", "application/json")

	HandleSendCoin(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var request map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &request))
	assert.Equal(t, enums.ErrCoinsInappropriateAmount.Error(), request["error"])
}

func TestHandleSendCoinInsufficientFunds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	sender := models.User{Username: "sender", Password: "password", Coins: 50}
	receiver := models.User{Username: "receiver", Password: "password", Coins: 1000}
	require.NoError(t, db.Create(&sender).Error)
	require.NoError(t, db.Create(&receiver).Error)

	requestBody := models.SendCoinRequest{
		ToUser: "receiver",
		Amount: 100,
	}
	bodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", sender)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(bodyBytes))
	ctx.Request.Header.Set("Content-Type", "application/json")

	HandleSendCoin(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.ErrInsufficientMoney.Error(), response["error"])
}

func TestHandleInfoUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/info", nil)

	HandleInfo(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.ErrUserNotAuthorized.Error(), response["error"])
}

func TestHandleInfoDBUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	user := models.User{Username: "ghost_user"}
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/info", nil)

	HandleInfo(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.ErrReceivingCoinsInfo.Error(), response["error"])
}

func TestHandleInfoSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	user := models.User{Username: "user", Password: "password", Coins: 1500}
	require.NoError(t, db.Create(&user).Error)

	purchases := []models.Purchase{
		{UserID: user.ID, MerchItem: "t-shirt", CreatedAt: time.Now()},
		{UserID: user.ID, MerchItem: "t-shirt", CreatedAt: time.Now()},
		{UserID: user.ID, MerchItem: "cup", CreatedAt: time.Now()},
	}
	for _, p := range purchases {
		require.NoError(t, db.Create(&p).Error)
	}

	sender := models.User{Username: "sender", Password: "password", Coins: 1002}
	require.NoError(t, db.Create(&sender).Error)
	transferReceived := models.CoinTransfer{
		FromUserID: sender.ID,
		ToUserID:   user.ID,
		Amount:     300,
		CreatedAt:  time.Now(),
	}
	require.NoError(t, db.Create(&transferReceived).Error)

	receiver := models.User{Username: "receiver", Password: "password", Coins: 1052}
	require.NoError(t, db.Create(&receiver).Error)
	transferSent := models.CoinTransfer{
		FromUserID: user.ID,
		ToUserID:   receiver.ID,
		Amount:     200,
		CreatedAt:  time.Now(),
	}
	require.NoError(t, db.Create(&transferSent).Error)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", user)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/info", nil)

	HandleInfo(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.InfoResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	assert.Equal(t, 1500, response.Coins)

	assert.Len(t, response.Inventory, 2)
	for _, item := range response.Inventory {
		if item.Type == "t-shirt" {
			assert.Equal(t, 2, item.Quantity)
		} else if item.Type == "cup" {
			assert.Equal(t, 1, item.Quantity)
		} else {
			t.Errorf("Unexpected inventory item: %s", item.Type)
		}
	}

	assert.Len(t, response.CoinHistory.Received, 1)
	assert.Equal(t, "sender", response.CoinHistory.Received[0].FromUser)
	assert.Equal(t, 300, response.CoinHistory.Received[0].Amount)

	assert.Len(t, response.CoinHistory.Sent, 1)
	assert.Equal(t, "receiver", response.CoinHistory.Sent[0].ToUser)
	assert.Equal(t, 200, response.CoinHistory.Sent[0].Amount)
}

func TestHandleBuyNoItemProvided(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	user := models.User{Username: "user", Password: "password", Coins: 1000}
	require.NoError(t, db.Create(&user).Error)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", user)
	ctx.Params = gin.Params{gin.Param{Key: "item", Value: ""}}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/buy/", nil)

	HandleBuy(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.ErrNotProvidedItem.Error(), response["error"])
}

func TestHandleBuyItemNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	user := models.User{Username: "user", Password: "password", Coins: 1001}
	require.NoError(t, db.Create(&user).Error)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", user)
	ctx.Params = gin.Params{gin.Param{Key: "item", Value: "nonexistent"}}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/buy/nonexistent", nil)

	HandleBuy(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.ErrItemNotFound.Error(), response["error"])
}

func TestHandleBuyInsufficientFunds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	user := models.User{Username: "user", Password: "password", Coins: 100}
	require.NoError(t, db.Create(&user).Error)
	merch := models.Merch{Name: "hoodie", Price: 300}
	require.NoError(t, db.Create(&merch).Error)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", user)
	ctx.Params = gin.Params{gin.Param{Key: "item", Value: "hoodie"}}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/buy/hoodie", nil)

	HandleBuy(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.ErrBuyWithInsufficientMoney.Error(), response["error"])
}

func TestHandleBuySuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	user := models.User{Username: "user", Password: "password", Coins: 1000}
	require.NoError(t, db.Create(&user).Error)
	merch := models.Merch{Name: "cap", Price: 250}
	require.NoError(t, db.Create(&merch).Error)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("user", user)
	ctx.Params = gin.Params{gin.Param{Key: "item", Value: "cap"}}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/buy/cap", nil)

	HandleBuy(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, enums.SuccessfulPurchase.String(), response["message"])

	var updatedUser models.User
	require.NoError(t, db.First(&updatedUser, user.ID).Error)
	assert.Equal(t, 750, updatedUser.Coins)

	var purchase models.Purchase
	err := db.First(&purchase, "user_id = ? AND merch_item = ?", user.ID, "cap").Error
	require.NoError(t, err)
}
