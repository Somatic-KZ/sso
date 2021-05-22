# Techno SSO

Сервер единого входа/регистрации, а также, хранения данных о пользовательских аккаунтах.

## Конфигурация

Конфигурирование осуществляется путем передачи переменных окружения. Основные настройки:

  * DATASTORE_URL - URL для соединения с MongoDB. Например: `mongodb://user:secret@host1:27018,host2:27018,host3:27018/sso?authMechanism=SCRAM-SHA-1`
  * JWT_KEY - секретный ключ генерации и валидации JWT токена.
  
Дополнительные настройки:

  * TOKEN_TTL - время жизни `access_token` в минутах.
  * REFRESH_TOKEN_TTL - время жизни `refresh_token` в днях.
  * NOTIFY_PREF - префикс системы, вы вызывающий SMS шлюз.
  * NOTIFY_SEND - идентификатор отправителя SMS. Должен быть "TECHNODOM.".
  * LISTEN - адрес интерфейса сервиса, обслуживающего API.
  * CERT_FILE - путь к файлу сертификата.
  * KEY_FILE - путь к ключу сертификата. 
  * DEBUG - режим отладки.
  * APP_TESTING - режим тестирования (все SMS имеют код "1111").
  
Настройки штрафов:

  * RECOVERY_SPAM_PENALTY - время повторной отсылки SMS при восстановлении.
  * VERIFY_SPAM_PENALTY - время повторной отсылки SMS при верификации с указанным в запросе TDID.
  * VERIFY_LONG_SPAM_PENALTY - время повторной отсылки SMS при верификации без TDID в запросе.

## Примеры операций API

### Вход пользователя

![signin](docs/img/sso-signin.png)

Структура запроса:

  * Метод: POST 
  * URL: /api/v1/auth/signin/phone
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "phone": "77078275615",
            "password": "secret"
        }
        
  * Код успешного ответа: 200
  * Тело успешного ответа содержит JWT токен авторизации, который надо использовать для обращения к точкам API, требующих авторизации:
  
        {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4NzIzYzg2MWRjMmQ0MGQzMjEwY2E3IiwiZXhwIjoxNTg2MTQ5NzQyfQ.5OC4yTE3LSPCLhzWklEafVyl2x0NkKXZ6zsHd9YIPAE",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU5NWE0ZDE2NGZiNDBhMTU3N2Y4ZjE2IiwiZXhwIjoxNTg2ODY3Mjg0fQ.-r5SuAl4ukoxHMsh4kcGAIAzYFneK-m0mtQTqTSi9cw",
            "status": "SignIn success"
        }
          
  После этого вызова пользователь считается залогиненым. Токен приходит в заголовок "Authorization":
  
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4YWFlZWQ1OWU0MmIzYjliZTEwMWUwIiwiZXhwIjoxNTg2MTQ4ODgyfQ.nC92BAooHu-WjyO13zSzMJf7LGVHewrRGxLjDeVEXQg
  
  `refresh_token` можно использовать для получения новых `access_token`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           >
  Все коды ошибок и детали определены в спецификации "swagger.yaml".


### Регистрация пользователя

![signup](docs/img/sso-signup.png)

Структура запроса:

  * Метод: PUT 
  * URL: /api/v1/auth/signup
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "email": "user1@example.com",
            "phone": "87078275611",
            "password": "secret"
        }
        
  * Код успешного ответа: 201
  * Тело успешного ответа:
  
        {
            "tdid": "5e8aaeed59e42b3b9be101e0",
            "status": "SignUP successful. Verification Required"
        }
  
  После успешного вызова метода создается новый глобальный идентификатор пользователя "tdid". После прохождения регистрации пользователь не может залогиниться в систему, так как он обязан пройти верификацию телефона.

### Верификация номера телефона

![verify](docs/img/sso-verify.png)

Является частью процесса регистрации, вызывается после "/api/v1/auth/signup".

  * Метод: PUT 
  * URL: /api/v1/auth/verify/phone
  * Тело: JSON
  * Минимальная структура запроса содержит сгенерированный ранее идентификатор пользователя "tdid" и номер верифицированного телефона:
  
        {
            "phone": "87078275611",
            "tdid": "5e8aaeed59e42b3b9be101e0"
        }
          
  * Код успешного ответа: 201
  * Тело успешного ответа:
  
        {
            "status": "Token was created successfully"
        }
          
  После успешного вызова пользователю отправляется на телефон СМС с **токеном верификации** (4 цифры).

#### Верификация телефона: валидация токена верификации

Является частью процесса регистрации, вызывается после "/api/v1/auth/verify/phone".

  * Метод: PUT 
  * URL: /api/v1/auth/verify/token
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "token": "7266",
            "tdid": "5e8aaeed59e42b3b9be101e0"
        }
                  
  * Код успешного ответа: 200
  * Тело успешного ответа содержит JWT токен авторизации, который надо использовать для обращения к точкам API, требующих аутентификации:
  
        {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4YWFlZWQ1OWU0MmIzYjliZTEwMWUwIiwiZXhwIjoxNTg2MTQ4ODgyfQ.nC92BAooHu-WjyO13zSzMJf7LGVHewrRGxLjDeVEXQg",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU5NWE0ZDE2NGZiNDBhMTU3N2Y4ZjE2IiwiZXhwIjoxNTg2ODY3Mjg0fQ.-r5SuAl4ukoxHMsh4kcGAIAzYFneK-m0mtQTqTSi9cw",
            "status": "SignIn success"
        }
  
    После этого вызова пользователь считается зарегистрированным и может выполять вход и выход из системы. Полученный "access_token" токен можно сохранить клиенту и передавать его в заголовке:
    
        Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4YWFlZWQ1OWU0MmIzYjliZTEwMWUwIiwiZXhwIjoxNTg2MTQ4ODgyfQ.nC92BAooHu-WjyO13zSzMJf7LGVHewrRGxLjDeVEXQg
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            >
    А `refresh_token` можно использовать для получения новых `access_token`. 
    Поле "tdid" не является обязательным в запросе, если оно не указано, то повторно вызвать точку верификации можно будет через бОльший по времени штраф.
     
### Восстановление доступа

![recovery](docs/img/sso-recovery.png)

Процесс восстановления доступа к аккаунту состоить из двух обязательных вызовов API:

  * Создание пользователю токена на восстановление (СМС код, состоящий из четырех цифр)
  * Установка пароля аккаунту
  
Также есть необязательный вызов, проверяющий валидность токена на восстановление.

#### Восстановление доступа по номеру телефона

Это обязательный шаг, генерирующий четырехзначный код и отсылающий его пользователю в виде SMS.

  * Метод: PUT 
  * URL: /api/v1/auth/recovery/phone
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "phone": "87078275625"
        }
        
  * Код успешного ответа: 201
  * Тело успешного ответа содержит статус:
  
        {
          "status": "Recovery token was created successfully"
        }
  
#### Валидация токена на восстановление

Это опциональный шаг проверки валидности SMS кода.

  * Метод: PUT 
  * URL: /api/v1/auth/recovery/phone/token
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "phone": "87078275625",
            "token": "1111"
        }
        
  * Код успешного ответа: 200
  * Тело успешного ответа содержит статус:
  
        {
          "status": "Token is valid"
        }
  
#### Установка пароля аккаунту

Это обязательный шаг, на котором пользователь устанавливает новый пароль своему аккаунту. После первого успешного вызова этого метода с валидным токеном на восстановление - токен перестает существовать.

  * Метод: PUT 
  * URL: /api/v1/auth/recovery/password
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "phone": "87078275625",
            "token": "1111",
            "password": "secret"
        }
        
  * Код успешного ответа: 200
  * Тело успешного ответа содержит статус:
  
        {
          "status": "Your password was changed"
        }
