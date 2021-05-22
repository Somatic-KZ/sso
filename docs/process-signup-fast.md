# Быстрая регистрация пользователя

Эта регистрация используется для того, чтобы создать аккаунт пользователя на этапе оплаты корзины.

  * Метод: PUT 
  * URL: /api/v1/auth/signup/fast
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "email": "user1@example.com",
            "phone": "87078275611"
        }
        
  * Код успешного ответа: 201
  * Тело успешного ответа:
  
          {
              "tdid": "5e8aaeed59e42b3b9be101e0",
              "status": "SignUP successful. Verification Required"
          }
  
  После успешного вызова метода создается новый глобальный идентификатор пользователя "tdid". После прохождения регистрации пользователь не может залогиниться в систему, так как он обязан пройти верификацию телефона.
  Верификацию вручную начинать на этом шаге не нужно. SMS отправится пользователю сразу.
  Все коды ошибок, опциональных полей и деталей в спецификации "swagger.yaml".
  
## Верификация телефона: валидация токена верификации

Является частью процесса быстрой регистрации, вызывается после "/api/v1/auth/signup/fast".

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