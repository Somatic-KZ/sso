# Регистрация обычного пользователя

Регистрация пользователя начинается здесь.

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
  Все коды ошибок, опциональных полей и деталей в спецификации "swagger.yaml".
  
## Верификация телефона: отправка токена верификации

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
  Все коды ошибок и деталей в спецификации "swagger.yaml".

## Верификация телефона: валидация токена верификации

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
    А "refresh_token" можно использовать для получения новых "access_token".  
    Все коды ошибок и детали определены в спецификации "swagger.yaml".
