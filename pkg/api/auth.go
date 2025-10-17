package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// в беклог - вынести в ENV?
const secretKey = "my_secret_word"

var secretMethod = *jwt.SigningMethodHS256

func generatePasswordHash(password string) string {
	// Создаем хеш sha256
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedBytes := hasher.Sum(nil)

	// Преобразуем хеш в шестнадцатиричную строку
	return hex.EncodeToString(hashedBytes)
}

// comparePasswordHash проверяет, совпадает ли текущий пароль с хешом из токена
func comparePasswordHash(hashFromToken string, currentPassword string) bool {
	currentHash := generatePasswordHash(currentPassword)

	return hashFromToken != currentHash
}

func generateToken(password string) (string, error) {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
		"hash": generatePasswordHash(password),
	}

	token := jwt.NewWithClaims(&secretMethod, claims)
	return token.SignedString([]byte(secretKey))
}

func isValidPassword(token string, pass string) bool {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		// Проверяем метод подписи
		if token.Method.Alg() != secretMethod.Name {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return false
	}

	// Проверяем, валиден ли токен
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Проверяем срок действия токена
		exp, exists := claims["exp"].(float64)
		if !exists || time.Now().After(time.Unix(int64(exp), 0)) {
			return false
		}

		// Проверяем соответствие пароля хешу
		hash, exists := claims["hash"].(string)
		if !exists || comparePasswordHash(hash, pass) {
			return false
		}

		return true
	}

	return false
}

func authHandler(res http.ResponseWriter, req *http.Request, logger *log.Logger) {
	var creds struct {
		Password string `json:"password"`
	}
	globalPass, ok := os.LookupEnv("TODO_PASSWORD")
	if !ok || len(globalPass) == 0 {
		globalPass = ""
	}

	body := req.Body
	defer body.Close()

	err := json.NewDecoder(body).Decode(&creds)
	if err != nil {
		errText := "ошибка десериализации JSON"
		logger.Printf("%s: %v", errText, err)
		err = jsonError(res, errText, http.StatusBadRequest)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	if creds.Password != globalPass {
		errText := "Неверный пароль"
		err := jsonError(res, errText, http.StatusUnauthorized)
		if err != nil {
			logger.Println(err)
		}
		return
	}
	//Далее считаем, что юзер прошел аутентификацию
	token, err := generateToken(creds.Password)
	if err != nil {
		errText := "Ошибка генерации токена"
		logger.Printf("%s: %v", errText, err)
		err = jsonError(res, errText, http.StatusInternalServerError)
		if err != nil {
			logger.Println(err)
		}
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",                     // Будет доступен по всему домену
		HttpOnly: true,                    // Доступна только для браузера, JavaScript не видит
		SameSite: http.SameSiteStrictMode, // Ограничивает доступ третьих сторон
		Secure:   true,                    // Обязательная мера, если сайт работает по HTTPS
		MaxAge:   24 * 60 * 60,            // Жизненный цикл куки — 8 часов
	}
	http.SetCookie(res, cookie)
	//вроде поставил reddirect
	// http.Redirect(res, req, "/", http.StatusSeeOther)
	err = writeJson(res, token, http.StatusOK)
	if err != nil {
		logger.Println(err)
	}
}

// Костыль, понимаю
func Wrap(next func(w http.ResponseWriter, r *http.Request, logger *log.Logger), middle func(w http.ResponseWriter, r *http.Request, logger *log.Logger) bool, logger *log.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Сначала вызываем middleware
		continueProcessing := middle(w, r, logger)

		// Если middleware разрешает продолжение обработки
		if continueProcessing {
			// Вызываем основной обработчик
			next(w, r, logger)
		}
	})
}

// func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc

//	func AuthMiddleware(next func(w http.ResponseWriter, r *http.Request, logger *log.Logger), logger *log.Logger) http.HandlerFunc {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
func AuthMiddleware(w http.ResponseWriter, r *http.Request, logger *log.Logger) bool {
	var token *http.Cookie

	pass := os.Getenv("TODO_PASSWORD")
	// Если пароль не задан - пропускаем проверку
	if len(pass) == 0 {
		return true
	}

	token, err := r.Cookie("token")
	if err == nil && isValidPassword(token.Value, pass) {
		logger.Printf("Прошли прослойку аутентификации %s", r.URL)
		return true
	}
	// дальнейшая логика работает если кука невалидна:
	var errText string
	var status int
	switch {
	case err != http.ErrNoCookie:
		errText = "Cookie processing error"
		status = http.StatusInternalServerError
		logger.Printf("%s: %v", errText, err)
	default:
		errText = "Authentication required"
		status = http.StatusUnauthorized
	}

	err = jsonError(w, errText, status)
	if err != nil {
		logger.Println(err)
	}
	return false
}
