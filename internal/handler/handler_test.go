package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enum"
	"github.com/ners1us/merch_store/internal/model"
	"github.com/ners1us/merch_store/internal/repository"
	"github.com/ners1us/merch_store/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var container testcontainers.Container
var db *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	request := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	var err error
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Не удалось запустить контейнер: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Не удалось получить host контейнера: %s", err)
	}
	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Не удалось получить mapped port: %s", err)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=test_db", host, mappedPort.Port())
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %s", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Merch{}, &model.Purchase{}, &model.CoinTransfer{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %s", err)
	}

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("Ошибка остановки контейнера: %s", err)
	}
	os.Exit(code)
}

func setupRouter() *gin.Engine {
	jwtSecret := []byte("elaborate_secret")
	userRepo := repository.NewUserRepository(db)
	merchRepo := repository.NewMerchRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)
	transferRepo := repository.NewCoinTransferRepository(db)

	authService := service.NewAuthService(userRepo, jwtSecret)
	userService := service.NewUserService(userRepo, purchaseRepo, transferRepo)
	merchService := service.NewMerchService(userRepo, merchRepo, purchaseRepo)
	transferService := service.NewTransferService(userRepo, transferRepo)

	authHandler := NewAuthHandler(authService)
	infoHandler := NewInfoHandler(userService)
	buyHandler := NewBuyHandler(merchService)
	sendCoinHandler := NewSendCoinHandler(transferService)

	router := gin.Default()
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/auth", authHandler.HandleAuth)
		apiRoutes.GET("/info", AuthMiddleware(jwtSecret), infoHandler.HandleInfo)
		apiRoutes.POST("/sendCoin", AuthMiddleware(jwtSecret), sendCoinHandler.HandleSendCoin)
		apiRoutes.GET("/buy/:item", AuthMiddleware(jwtSecret), buyHandler.HandleBuy)
	}
	return router
}

func clearDB() {
	db.Exec("TRUNCATE TABLE coin_transfers, purchases, merches, users RESTART IDENTITY CASCADE")
}

func performAuth(t *testing.T, serverURL, username, password string) string {
	authRequest := model.AuthRequest{
		Username: username,
		Password: password,
	}
	payload, err := json.Marshal(authRequest)
	if err != nil {
		t.Fatalf("Ошибка маршалинга запроса: %v", err)
	}
	response, err := http.Post(serverURL+"/api/auth", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Ошибка запроса /api/auth: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200 OK, но получен %d", response.StatusCode)
	}
	var authResponse model.AuthResponse
	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа аутентификации: %v", err)
	}
	return authResponse.Token
}

func TestBuyMerch(t *testing.T) {
	// Arrange
	clearDB()
	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	merchItem := model.Merch{
		Name:  "t-shirt",
		Price: 500,
	}
	db.Create(&merchItem)
	token := performAuth(t, ts.URL, "ners1us", "thelongestpasswordever")

	// Act
	request, err := http.NewRequest("GET", ts.URL+"/api/buy/t-shirt", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса покупки: %v", err)
	}
	defer response.Body.Close()

	// Assert
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200 OK, но получен %d", response.StatusCode)
	}
	var buyResponse map[string]string
	err = json.NewDecoder(response.Body).Decode(&buyResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа покупки: %v", err)
	}
	assert.Equal(t, enum.SuccessfulPurchase.String(), buyResponse["message"])

	// Act
	request, err = http.NewRequest("GET", ts.URL+"/api/info", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer response.Body.Close()

	// Assert
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, но получен %d", response.StatusCode)
	}
	var infoResponse model.InfoResponse
	err = json.NewDecoder(response.Body).Decode(&infoResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
	expectedCoins := 1000 - merchItem.Price
	assert.Equal(t, expectedCoins, infoResponse.Coins)
}

func TestSendCoin(t *testing.T) {
	// Arrange
	clearDB()
	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	tokenSender := performAuth(t, ts.URL, "johnnyBravo", "ilikepie")
	tokenReceiver := performAuth(t, ts.URL, "darthVader", "iamlivingcorpse")

	sendCoinRequest := model.SendCoinRequest{
		ToUser: "darthVader",
		Amount: 200,
	}
	payload, err := json.Marshal(sendCoinRequest)
	if err != nil {
		t.Fatalf("Ошибка маршалинга запроса отправки монет: %v", err)
	}

	// Act
	request, err := http.NewRequest("POST", ts.URL+"/api/sendCoin", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Ошибка создания запроса отправки монет: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+tokenSender)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса отправки монет: %v", err)
	}
	defer response.Body.Close()

	// Assert
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200 OK, но получен %d", response.StatusCode)
	}
	var transferResponse map[string]string
	err = json.NewDecoder(response.Body).Decode(&transferResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа отправки монет: %v", err)
	}
	assert.Equal(t, enum.SuccessfulTransfer.String(), transferResponse["message"])

	// Act
	request, err = http.NewRequest("GET", ts.URL+"/api/info", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса для отправителя: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+tokenSender)
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса для отправителя: %v", err)
	}
	defer response.Body.Close()

	// Assert
	var infoSender model.InfoResponse
	err = json.NewDecoder(response.Body).Decode(&infoSender)
	if err != nil {
		t.Fatalf("Ошибка декодирования для отправителя: %v", err)
	}
	expectedSenderCoins := 800
	assert.Equal(t, expectedSenderCoins, infoSender.Coins)

	// Act
	request, err = http.NewRequest("GET", ts.URL+"/api/info", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса для получателя: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+tokenReceiver)
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса для получателя: %v", err)
	}
	defer response.Body.Close()

	// Assert
	var infoReceiver model.InfoResponse
	err = json.NewDecoder(response.Body).Decode(&infoReceiver)
	if err != nil {
		t.Fatalf("Ошибка декодирования для получателя: %v", err)
	}
	expectedReceiverCoins := 1200
	assert.Equal(t, expectedReceiverCoins, infoReceiver.Coins)
}
