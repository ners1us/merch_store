package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ners1us/merch_store/internal/middleware"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ners1us/merch_store/internal/enums"
	"github.com/ners1us/merch_store/internal/models"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var container testcontainers.Container

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

	err = db.AutoMigrate(&models.User{}, &models.Merch{}, &models.Purchase{}, &models.CoinTransfer{})
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
	Init(db, []byte("elaborate_secret"))
	middleware.Init([]byte("elaborate_secret"))

	router := gin.Default()
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/auth", HandleAuth)
		apiRoutes.GET("/info", middleware.AuthMiddleware(), HandleInfo)
		apiRoutes.POST("/sendCoin", middleware.AuthMiddleware(), HandleSendCoin)
		apiRoutes.GET("/buy/:item", middleware.AuthMiddleware(), HandleBuy)
	}
	return router
}

func clearDB() {
	db.Exec("TRUNCATE TABLE coin_transfers, purchases, merches, users RESTART IDENTITY CASCADE")
}

func performAuth(t *testing.T, serverURL, username, password string) string {
	authRequest := models.AuthRequest{
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

	var authResponse models.AuthResponse
	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа аутентификации: %v", err)
	}
	return authResponse.Token
}

func TestBuyMerch(t *testing.T) {
	clearDB()
	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	merchItem := models.Merch{
		Name:  "t-shirt",
		Price: 500,
	}

	db.Create(&merchItem)

	token := performAuth(t, ts.URL, "ners1us", "thelongestpasswordever")

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

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200 OK, но получен %d", response.StatusCode)
	}

	var buyResponse map[string]string
	err = json.NewDecoder(response.Body).Decode(&buyResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа покупки: %v", err)
	}
	if buyResponse["message"] != enums.SuccessfulPurchase.String() {
		t.Errorf("Ожидалось сообщение %q, но получено %q", enums.SuccessfulPurchase.String(), buyResponse["message"])
	}

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

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, но получен %d", response.StatusCode)
	}

	var infoResponse models.InfoResponse
	err = json.NewDecoder(response.Body).Decode(&infoResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	expectedCoins := 1000 - merchItem.Price
	if infoResponse.Coins != expectedCoins {
		t.Errorf("Ожидалось %d монет, но получено %d", expectedCoins, infoResponse.Coins)
	}
}

func TestSendCoin(t *testing.T) {
	clearDB()
	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	tokenSender := performAuth(t, ts.URL, "johnnyBravo", "ilikepie")
	tokenReceiver := performAuth(t, ts.URL, "darthVader", "iamlivingcorpse")

	sendCoinRequest := models.SendCoinRequest{
		ToUser: "darthVader",
		Amount: 200,
	}
	payload, err := json.Marshal(sendCoinRequest)
	if err != nil {
		t.Fatalf("Ошибка маршалинге запроса отправки монет: %v", err)
	}
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

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200 OK, но получен %d", response.StatusCode)
	}
	var transferResponse map[string]string
	err = json.NewDecoder(response.Body).Decode(&transferResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа отправки монет: %v", err)
	}
	if transferResponse["message"] != enums.SuccessfulTransfer.String() {
		t.Errorf("Ожидалось сообщение %q, но получено %q", enums.SuccessfulTransfer.String(), transferResponse["message"])
	}

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

	var infoSender models.InfoResponse
	err = json.NewDecoder(response.Body).Decode(&infoSender)
	if err != nil {
		t.Fatalf("Ошибка декодирования для отправителя: %v", err)
	}
	expectedSenderCoins := 800
	if infoSender.Coins != expectedSenderCoins {
		t.Errorf("У отправителя ожидалось %d монет, но получено %d", expectedSenderCoins, infoSender.Coins)
	}

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

	var infoReceiver models.InfoResponse
	err = json.NewDecoder(response.Body).Decode(&infoReceiver)
	if err != nil {
		t.Fatalf("Ошибка декодирования для получателя: %v", err)
	}
	expectedReceiverCoins := 1200
	if infoReceiver.Coins != expectedReceiverCoins {
		t.Errorf("У получателя ожидалось %d монет, но получено %d", expectedReceiverCoins, infoReceiver.Coins)
	}
}
